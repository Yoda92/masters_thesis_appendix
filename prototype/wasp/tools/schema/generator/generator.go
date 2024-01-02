// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/iotaledger/wasp/tools/schema/model"
)

// TODO nested structs
// TODO handle case where owner is type AgentID[]

type IGenerator interface {
	Build() error
	Cleanup()
	GenerateImplementation() error
	GenerateInterface() error
	GenerateTests() error
	GenerateWasmStub() error
	IsLatest() bool
}

type Generator struct {
	currentEvent  *model.Struct
	currentField  *model.Field
	currentFunc   *model.Func
	currentStruct *model.Struct
	emitters      map[string]func(g *Generator)
	extension     string
	file          *os.File
	folder        string
	funcRegexp    *regexp.Regexp
	keys          model.StringMap
	language      string
	newTypes      map[string]bool
	path          string
	rootFolder    string
	s             *model.Schema
	subFolder     string
	tab           int
	templates     model.StringMap
	tmp           bool
	typeDependent model.StringMapMap
}

const spaces = "                                             "

func (g *Generator) init(s *model.Schema, typeDependent model.StringMapMap, templates []map[string]string) {
	g.s = s
	g.typeDependent = typeDependent

	g.emitters = map[string]func(g *Generator){}
	g.keys = model.StringMap{}
	g.newTypes = map[string]bool{}
	g.templates = model.StringMap{}

	config := templates[0]
	g.language = config["language"]
	g.extension = config["extension"]
	g.rootFolder = config["rootFolder"]
	g.subFolder = config["subFolder"]
	g.funcRegexp = regexp.MustCompile(config["funcRegexp"])

	g.addTemplates(commonTemplates)
	for _, template := range templates {
		g.addTemplates(template)
	}

	g.setCommonKeys()
}

func (g *Generator) addTemplates(t model.StringMap) {
	for k, v := range t {
		g.templates[k] = v
	}
}

func (g *Generator) build(compiler string, args string) error {
	command := compiler
	cmd := exec.Command(command, strings.Split(args, " ")...)
	var stdout strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(stdout.String())
		return errors.New("build failed")
	}
	return err
}

func (g *Generator) cleanCommonFiles() {
	g.generateCommonFolder("", false)
	g.cleanFolder(g.folder)

	g.generateCommonFolder("impl", true)
	g.cleanSourceFile("thunks")
	g.cleanSourceFile("../main")
	g.cleanFolder(g.folder + "../pkg")
}

func (g *Generator) cleanFolder(folder string) {
	_ = os.RemoveAll(folder)
	_ = os.Remove(folder)
}

func (g *Generator) cleanSourceFile(name string) {
	path := g.folder + name + g.extension
	_ = os.Remove(path)
}

func (g *Generator) cleanSourceFileIfSame(generator func() error) {
	g.tmp = true
	err := generator()
	g.tmp = false
	if err != nil {
		panic(err)
	}
	newFile, err := os.ReadFile(g.path + ".tmp")
	_ = os.Remove(g.path + ".tmp")
	if err != nil {
		panic(err)
	}
	oldFile, _ := os.ReadFile(g.path)
	if bytes.Equal(oldFile, newFile) {
		_ = os.Remove(g.path)
	}
}

func (g *Generator) createFile(path string, overwrite bool, generator func()) (err error) {
	g.path = path
	if g.tmp {
		path += ".tmp"
	}
	if !overwrite && g.exists(path) == nil {
		return nil
	}
	g.file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = g.file.Close() }()
	generator()
	return nil
}

func (g *Generator) createSourceFile(name string, mustExist bool, macro ...string) error {
	path := g.folder + name + g.extension
	if !mustExist {
		_ = os.Remove(path)
		return nil
	}
	if len(macro) == 1 {
		name = macro[0]
	}
	return g.createFile(path, true, func() {
		g.emit("warning")
		g.emit("copyrightMessage")
		g.emit(name + g.extension)
	})
}

func (g *Generator) error(what string) {
	g.println("???:" + what)
}

func (g *Generator) exists(path string) (err error) {
	_, err = os.Stat(path)
	return err
}

func (g *Generator) funcName(f *model.Func) string {
	name := f.Kind + capitalize(f.Name)
	if g.language == "Rust" {
		name = snake(name)
	}
	return name
}

func (g *Generator) generateCommonFolder(postfix string, withSubFolder bool) {
	g.folder = g.rootFolder + "/" + g.s.PackageName + postfix + "/"
	if g.s.CoreContracts {
		g.folder = g.rootFolder + "/wasmlib/" + g.s.PackageName + "/"
	}
	if withSubFolder && g.subFolder != "" {
		g.folder += g.subFolder + "/"
		if g.s.CoreContracts {
			g.folder = g.subFolder + "/" + g.s.PackageName + "/"
		}
	}
}

func (g *Generator) generateFuncs() error {
	scFileName := g.folder + "funcs" + g.extension
	if g.exists(scFileName) != nil {
		// generate initial SC function file
		return g.createFile(scFileName, false, func() {
			g.emit("copyrightMessage")
			g.emit("funcs" + g.extension)
		})
	}

	// append missing SC functions to existing code file

	// scan existing file for function names
	existing := make(model.StringMap)
	lines := make([]string, 0)
	err := g.scanExistingCode(scFileName, &existing, &lines)
	if err != nil {
		return err
	}

	// save old one from overwrite
	scOriginal := g.folder + "funcs.bak"
	err = os.Rename(scFileName, scOriginal)
	if err != nil {
		return err
	}

	err = g.createFile(scFileName, false, func() {
		// make copy of original file
		for _, line := range lines {
			g.println(line)
		}

		// append any new funcs
		for _, g.currentFunc = range g.s.Funcs {
			if existing[g.funcName(g.currentFunc)] == "" {
				g.setFuncKeys(false, 0, 0)
				g.emit("funcSignature")
			}
		}
	})
	if err != nil {
		return err
	}
	return os.Remove(scOriginal)
}

func (g *Generator) generateImplementation() error {
	g.generateCommonFolder("impl", true)
	err := os.MkdirAll(g.folder, 0o755)
	if err != nil {
		return err
	}
	err = g.createSourceFile("thunks", true)
	if err != nil {
		return err
	}
	return g.generateFuncs()
}

func (g *Generator) generateInterface() error {
	err := os.MkdirAll(g.folder, 0o755)
	if err != nil {
		return err
	}
	err = g.createSourceFile("consts", true)
	if err != nil {
		return err
	}
	err = g.createSourceFile("events", !g.s.CoreContracts && len(g.s.Events) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("eventhandlers", len(g.s.Events) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("structs", len(g.s.Structs) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("typedefs", len(g.s.Typedefs) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("params", len(g.s.Params) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("results", len(g.s.Results) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("state", !g.s.CoreContracts && len(g.s.StateVars) != 0)
	if err != nil {
		return err
	}
	return g.createSourceFile("contract", true)
}

func (g *Generator) GenerateTests() error {
	err := os.MkdirAll("test", 0o755)
	if err != nil {
		return err
	}

	// do not overwrite existing file
	name := strings.ToLower(g.s.PackageName)
	filename := "test/" + name + "_test.go"
	return g.createFile(filename, false, func() {
		g.emit("test.go")
	})
}

func (g *Generator) IsLatest() bool {
	g.generateCommonFolder("", true)

	info, err := os.Stat(g.folder + "consts" + g.extension)
	if err == nil && info.ModTime().After(g.s.SchemaTime) {
		// fmt.Printf("skipping %s code generation\n", g.language)
		return true
	}

	fmt.Printf("generating %s code\n", g.language)
	return false
}

func (g *Generator) openFile(path string, processor func() error) (err error) {
	g.file, err = os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = g.file.Close() }()
	return processor()
}

func (g *Generator) println(a ...interface{}) {
	_, _ = fmt.Fprintln(g.file, a...)
}

func (g *Generator) scanExistingCode(path string, existing *model.StringMap, lines *[]string) error {
	return g.openFile(path, func() error {
		scanner := bufio.NewScanner(g.file)
		for scanner.Scan() {
			line := scanner.Text()
			matches := g.funcRegexp.FindStringSubmatch(line)
			if matches != nil {
				(*existing)[matches[1]] = line
			}
			*lines = append(*lines, line)
		}
		return scanner.Err()
	})
}
