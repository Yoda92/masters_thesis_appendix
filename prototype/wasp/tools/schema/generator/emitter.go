// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/tools/schema/model"
	"github.com/iotaledger/wasp/tools/schema/model/yaml"
)

const enableLog = false

const (
	KeyArray     = "array"
	KeyBaseType  = "basetype"
	KeyCore      = "core"
	KeyEvent     = "event"
	KeyEvents    = "events"
	KeyExist     = "exist"
	KeyFunc      = "func"
	KeyFuncs     = "funcs"
	KeyInit      = "init"
	KeyMandatory = "mandatory"
	KeyMap       = "map"
	KeyMut       = "mut"
	KeyParam     = "param"
	KeyParams    = "params"
	KeyProxy     = "proxy"
	KeyPtrs      = "ptrs"
	KeyResult    = "result"
	KeyResults   = "results"
	KeyState     = "state"
	KeyStruct    = "struct"
	KeyStructs   = "structs"
	KeyThis      = "this"
	KeyTrue      = "true"
	KeyTypeDef   = "typedef"
	KeyTypeDefs  = "typedefs"
	KeyView      = "view"
)

var emitKeyRegExp = regexp.MustCompile(`\$[a-zA-Z_][a-zA-Z_0-9]*`)

func (g *Generator) indent() {
	g.tab++
}

func (g *Generator) undent() {
	g.tab--
}

func (g *Generator) log(text string) {
	if !enableLog {
		return
	}

	for i := 0; i < g.tab; i++ {
		fmt.Print("  ")
	}
	fmt.Println(text)
}

// emit processes "$#emit template"
// It processes all lines in the named template
// If the template is non-existent nothing will happen
// Any line starting with a special "$#" directive will recursively be processed
// An unknown directive will result in an error
func (g *Generator) emit(template string) {
	g.log("$#emit " + template)
	g.indent()
	defer g.undent()

	lines := strings.Split(g.templates[template], "\n")
	for i := 1; i < len(lines)-1; i++ {
		// replace any placeholder keys
		line := emitKeyRegExp.ReplaceAllStringFunc(lines[i], func(key string) string {
			text, ok := g.keys[key[1:]]
			if ok {
				return text
			}
			return "key???:" + key
		})
		line = strings.ReplaceAll(line, "\r", "\n")
		// remove concatenation markers
		line = strings.ReplaceAll(line, "$+", "")

		// line contains special directive?
		space := strings.Index(line, " ")
		if space <= 2 || line[:2] != "$#" {
			// no special directive, just emit line
			g.println(line)
			continue
		}

		// now process special directive
		switch line[2:space] {
		case "emit":
			g.emit(strings.TrimSpace(line[7:]))
		case "each":
			g.emitEach(line)
		case "func":
			g.emitFunc(line)
		case "if":
			g.emitIf(line)
		case "set":
			g.emitSet(line)
		default:
			g.error(line)
		}
	}
}

// emitEach processes "$#each array template"
// It processes the template for each item in the array
// Produces an error if the array key is unknown
func (g *Generator) emitEach(line string) {
	g.log(line)
	g.indent()
	defer g.undent()

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		g.error(line)
		return
	}

	template := parts[2]
	switch parts[1] {
	case KeyEvent:
		g.emitEachField(g.currentEvent.Fields, template)
	case KeyEvents:
		g.emitEachEvent(g.s.Events, template)
	case KeyFunc:
		g.emitEachFunc(g.s.Funcs, template)
	case KeyMandatory:
		g.emitEachMandatoryField(template)
	case KeyParam:
		g.emitEachField(g.currentFunc.Params, template)
	case KeyParams:
		g.emitEachField(g.s.Params, template)
	case KeyResult:
		g.emitEachField(g.currentFunc.Results, template)
	case KeyResults:
		g.emitEachField(g.s.Results, template)
	case KeyState:
		g.emitEachField(g.s.StateVars, template)
	case KeyStruct:
		g.emitEachField(g.currentStruct.Fields, template)
	case KeyStructs:
		g.emitEachStruct(g.s.Structs, template)
	case KeyTypeDef:
		g.emitEachField(g.s.Typedefs, template)
	default:
		g.emitEachLine(line, parts[1], template)
	}
}

func (g *Generator) emitEachEvent(events []*model.Struct, template string) {
	for _, g.currentEvent = range events {
		g.log("currentEvent: " + g.currentEvent.Name.Val)
		g.setMultiKeyValues("evtName", g.currentEvent.Name.Val)
		g.keys["eventComment"] = g.currentEvent.Name.Comment
		g.emit(template)
	}
}

func (g *Generator) emitEachField(fields []*model.Field, template string) {
	maxCamelLength := 0
	maxSnakeLength := 0
	for _, g.currentField = range fields {
		camelLen := len(g.currentField.Name)
		if maxCamelLength < camelLen {
			maxCamelLength = camelLen
		}
		snakeLen := len(snake(g.currentField.Name))
		if maxSnakeLength < snakeLen {
			maxSnakeLength = snakeLen
		}
	}

	for _, g.currentField = range fields {
		g.log("currentField: " + g.currentField.Name)
		g.setFieldKeys(true, maxCamelLength, maxSnakeLength)
		g.emit(template)
	}
}

func (g *Generator) emitEachFunc(funcs []*model.Func, template string) {
	maxCamelLength := 0
	maxSnakeLength := 0
	for _, g.currentFunc = range funcs {
		camelLen := len(g.currentFunc.Name)
		if maxCamelLength < camelLen {
			maxCamelLength = camelLen
		}
		snakeLen := len(snake(g.currentFunc.Name))
		if maxSnakeLength < snakeLen {
			maxSnakeLength = snakeLen
		}
	}

	for _, g.currentFunc = range funcs {
		g.log("currentFunc: " + g.currentFunc.Name)
		g.setFuncKeys(true, maxCamelLength, maxSnakeLength)
		g.emit(template)
	}
}

// emitEachLine emits multi-line text, like multi-line comments
func (g *Generator) emitEachLine(line string, key string, template string) {
	text, ok := g.keys[key]
	if !ok {
		g.error(line)
		return
	}
	if text != "" {
		lines := strings.Split(text, "\n")
		for _, nextLine := range lines {
			g.keys["nextLine"] = nextLine
			g.emit(template)
		}
	}
}

func (g *Generator) emitEachMandatoryField(template string) {
	mandatoryFields := make([]*model.Field, 0)
	for _, g.currentField = range g.currentFunc.Params {
		fld := g.currentField
		if !fld.IsOptional && fld.IsBaseType && !fld.IsArray && fld.MapKey == "" {
			mandatoryFields = append(mandatoryFields, g.currentField)
		}
	}
	g.emitEachField(mandatoryFields, template)
}

func (g *Generator) emitEachStruct(structs []*model.Struct, template string) {
	for _, g.currentStruct = range structs {
		g.log("currentStruct: " + g.currentStruct.Name.Val)
		g.setMultiKeyValues("strName", g.currentStruct.Name.Val)
		g.keys["structComment"] = g.currentStruct.Name.Comment
		g.emit(template)
	}
}

// emitFunc processes "$#func emitter"
// It can call back into go code to emit more complex stuff
// Produces an error if emitter is unknown
func (g *Generator) emitFunc(line string) {
	g.log(line)
	g.indent()
	defer g.undent()

	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		g.error(line)
		return
	}

	emitter, ok := g.emitters[parts[1]]
	if ok {
		emitter(g)
		return
	}
	g.error(line)
}

// emitIf processes "$#if condition template [elseTemplate]"
// It processes template when the named condition is true
// It processes the optional elseTemplate when the named condition is false
// Produces an error if named condition is unknown

func (g *Generator) emitIf(line string) {
	g.log(line)
	g.indent()
	defer g.undent()

	parts := strings.Split(line, " ")
	if len(parts) < 3 || len(parts) > 4 {
		g.error(line)
		return
	}

	if g.emitIfCondition(parts[1]) {
		g.emit(parts[2])
		return
	}

	// else branch?
	if len(parts) == 4 {
		g.emit(parts[3])
	}
}

//nolint:gocyclo
func (g *Generator) emitIfCondition(key string) bool {
	switch key {
	case KeyArray:
		return g.currentField.IsArray
	case KeyBaseType:
		return g.currentField.IsBaseType
	case KeyCore:
		return g.s.CoreContracts
	case KeyEvent:
		return len(g.currentEvent.Fields) != 0
	case KeyEvents:
		return len(g.s.Events) != 0
	case KeyExist:
		return g.newTypes[g.keys[KeyProxy]]
	case KeyFunc:
		return g.keys["kind"] == KeyFunc
	case KeyFuncs:
		return len(g.s.Funcs) != 0
	case KeyInit:
		return g.currentFunc.Name == KeyInit
	case KeyMandatory:
		return !g.currentField.IsOptional
	case KeyMap:
		return g.currentField.MapKey != ""
	case KeyMut:
		return g.keys[KeyMut] == "Mutable"
	case KeyParam:
		return len(g.currentFunc.Params) != 0
	case KeyParams:
		return len(g.s.Params) != 0
	case KeyPtrs:
		return len(g.currentFunc.Params) != 0 || len(g.currentFunc.Results) != 0
	case KeyResult:
		return len(g.currentFunc.Results) != 0
	case KeyResults:
		return len(g.s.Results) != 0
	case KeyState:
		return len(g.s.StateVars) != 0
	case KeyStructs:
		return len(g.s.Structs) != 0
	case KeyThis:
		return g.currentField.Alias == KeyThis
	case KeyTypeDef:
		return g.fieldIsTypeDef()
	case KeyTypeDefs:
		return len(g.s.Typedefs) != 0
	case KeyView:
		return g.keys["kind"] == KeyView
	default:
		return g.keys[key] != ""
	}
}

// emitSet processes "$#set key value"
// It sets the specified key to value, which can be anything
// Just make sure there is a space after the key name
// The special key "exist" is used to add a newly generated type
// It can be used to prevent duplicate types from being generated
func (g *Generator) emitSet(line string) {
	g.log(line)

	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		g.error(line)
		return
	}

	key := parts[1]
	value := line[len(parts[0])+len(key)+2:]
	g.keys[key] = value

	if key == KeyExist {
		g.newTypes[value] = true
	}
}

func (g *Generator) fieldIsTypeDef() bool {
	for _, typeDef := range g.s.Typedefs {
		if typeDef.Name == g.currentField.Type {
			g.currentField = typeDef
			g.setFieldKeys(false, 0, 0)
			return true
		}
	}
	return false
}

func (g *Generator) setCommonKeys() {
	g.keys["env_wasmlib"] = ""
	g.keys["env_wasmvmhost"] = ""
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			g.keys["env_"+parts[0]] = strings.ReplaceAll(parts[1], "\\n", "\n")
		}
	}
	g.keys["false"] = ""
	g.keys[KeyTrue] = KeyTrue
	g.keys["nil"] = ""
	g.keys["space"] = " "
	g.keys["package"] = g.s.PackageName
	g.keys["package"+g.s.PackageName] = KeyTrue
	g.keys["Package"] = g.s.ContractName
	g.setMultiKeyValues("pkgName", g.s.ContractName)
	g.keys["module"] = moduleName + strings.Replace(moduleCwd[len(modulePath):], "\\", "/", -1)
	scName := g.s.PackageName
	if g.s.CoreContracts {
		// strip off "core" prefix
		scName = scName[4:]
	}
	g.keys["scName"] = scName
	g.keys["hscName"] = isc.Hn(scName).String()
	g.keys["scDesc"] = g.s.Description
	g.keys[yaml.KeyAuthor] = g.s.Author
	g.keys[yaml.KeyCopyright] = g.s.Copyright
	g.keys[yaml.KeyVersion] = g.s.Version
	g.keys[yaml.KeyLicense] = g.s.License
	g.keys[yaml.KeyRepository] = g.s.Repository
}

func (g *Generator) setFieldKeys(pad bool, maxCamelLength, maxSnakeLength int) {
	g.setMultiKeyValues("fldName", g.currentField.Name)
	g.setMultiKeyValues("fldType", g.currentField.Type)
	g.setMultiKeyValues("fldMapKey", g.currentField.MapKey)

	isArray := ""
	if g.currentField.IsArray {
		isArray = KeyTrue
	}
	g.keys["fldIsArray"] = isArray
	g.keys["fldIsMap"] = g.currentField.MapKey

	g.keys["fldAlias"] = g.currentField.Alias
	g.keys["fldComment"] = g.currentField.Comment
	g.keys["eventFldComment"] = g.currentField.Comment

	if pad {
		g.keys["fldPad"] = spaces[:maxCamelLength-len(g.keys["fldName"])]
		g.keys["fld_pad"] = spaces[:maxSnakeLength-len(g.keys["fld_name"])]
	}

	for fieldName, typeValues := range g.typeDependent {
		fieldValue := typeValues[g.currentField.Type]
		if fieldValue == "" {
			// get default value for this field
			// TODO make this smarter w.r.t. maps and arrays?
			fieldValue = typeValues[""]
		}
		g.keys[fieldName] = fieldValue

		if fieldName[:3] == "fld" {
			// we also want the 'fldKey' variant to facilitate the map key type
			fieldValue = typeValues[g.currentField.MapKey]
			if fieldValue == "" {
				// get default value for this field
				// TODO make this smarter w.r.t. maps and arrays?
				fieldValue = typeValues[""]
			}
			g.keys["fldKey"+fieldName[3:]] = fieldValue
		}
	}
}

func (g *Generator) setFuncKeys(pad bool, maxCamelLength, maxSnakeLength int) {
	g.setMultiKeyValues("funcName", g.currentFunc.Name)
	g.setMultiKeyValues("kind", g.currentFunc.Kind)
	g.keys["funcAlias"] = g.currentFunc.Alias
	g.keys["hFuncName"] = g.currentFunc.Hname.String()
	grant := g.currentFunc.Access.Val
	index := strings.Index(grant, "//")
	if index >= 0 {
		grant = strings.TrimSpace(grant[:index])
	}
	g.setMultiKeyValues("funcAccess", grant)
	g.keys["funcAccessComment"] = g.currentFunc.Access.Comment
	g.keys["funcComment"] = g.currentFunc.Comment
	if pad {
		g.keys["funcPad"] = spaces[:maxCamelLength-len(g.keys["funcName"])]
		g.keys["func_pad"] = spaces[:maxSnakeLength-len(g.keys["func_name"])]
	}
}

func (g *Generator) setMultiKeyValues(key, value string) {
	value = uncapitalize(value)
	g.keys[key] = filterIDorVM(value)
	g.keys[capitalize(key)] = filterIDorVM(capitalize(value))
	g.keys[snake(key)] = snake(value)
	g.keys[upper(snake(key))] = upper(snake(value))
}
