// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

//nolint:dupl
package yaml_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/tools/schema/model"
	"github.com/iotaledger/wasp/tools/schema/model/yaml"
)

func TestConvert(t *testing.T) {
	type args struct {
		path string
	}
	type wants struct {
		out *model.SchemaDef
	}
	type test struct {
		args  args
		wants wants
	}

	tests := map[string]func(_ *testing.T) test{
		"successfully test1": func(_ *testing.T) test {
			return test{
				args: args{
					path: "testdata/test1.yaml",
				},
				wants: wants{
					out: &model.SchemaDef{
						Name: model.DefElt{
							Val:  "SchemaComment",
							Line: 1,
						},
						Description: model.DefElt{
							Val:  "test description",
							Line: 2,
						},
						Events: model.DefMapMap{
							model.DefElt{Val: "TestEvent1", Line: 7, Comment: "// line comment for TestEvent1"}: &model.DefMap{
								model.DefElt{Val: "eventParam11", Line: 8, Comment: "// line comment for eventParam11"}: &model.DefElt{
									Val:  "String",
									Line: 8,
								},
							},
							model.DefElt{Val: "TestEvent2", Line: 9, Comment: "// line comment for TestEvent2"}: &model.DefMap{
								model.DefElt{Val: "eventParam21", Line: 10, Comment: "// line comment for eventParam21"}: &model.DefElt{
									Val:  "String",
									Line: 10,
								},
								model.DefElt{Val: "eventParam22", Line: 11, Comment: "// line comment for eventParam22"}: &model.DefElt{
									Val:  "String",
									Line: 11,
								},
							},
						},
						Structs: model.DefMapMap{
							model.DefElt{Val: "TestStruct1", Line: 16}: &model.DefMap{
								model.DefElt{Val: "x1", Line: 17}: &model.DefElt{
									Val:  "Int32",
									Line: 17,
								},
								model.DefElt{Val: "y1", Line: 18}: &model.DefElt{
									Val:  "Int32",
									Line: 18,
								},
							},
							model.DefElt{Val: "TestStruct2", Line: 20, Comment: "// comment for TestStruct2"}: &model.DefMap{
								model.DefElt{Val: "x2", Line: 21, Comment: "// comment for x2"}: &model.DefElt{
									Val:  "Int32",
									Line: 21,
								},
								model.DefElt{Val: "y2", Line: 25, Comment: "// comment for y2 1\n// comment for y2 2"}: &model.DefElt{
									Val:  "Int32",
									Line: 26,
								},
							},
						},
						Views: model.FuncDefMap{
							model.DefElt{Val: "TestView1", Line: 29, Comment: "// comment for TestView1"}: &model.FuncDef{
								Line:    29,
								Comment: "// comment for TestView1",
								Access: model.DefElt{
									Val:     "owner",
									Line:    30,
									Comment: "// comment for access",
								},
								Params: model.DefMap{
									model.DefElt{Val: "name", Line: 32, Comment: "// comment for name"}: &model.DefElt{
										Val:  "String",
										Line: 32,
									},
								},
								Results: model.DefMap{
									model.DefElt{Val: "length", Line: 34, Comment: "// comment for length"}: &model.DefElt{
										Val:  "Uint32",
										Line: 34,
									},
								},
							},
						},
					},
				},
			}
		},
		"successfully test2": func(_ *testing.T) test {
			return test{
				args: args{
					path: "testdata/test2.yaml",
				},
				wants: wants{
					out: &model.SchemaDef{
						Name: model.DefElt{
							Val:  "SchemaComment",
							Line: 1,
						},
						Description: model.DefElt{
							Val:  "test description",
							Line: 2,
						},
						Events: model.DefMapMap{
							model.DefElt{Val: "TestEvent1", Line: 15, Comment: "// header comment for TestEvent1 1\n// header comment for TestEvent1 2"}: &model.DefMap{
								model.DefElt{Val: "eventParam1", Line: 22, Comment: "// header comment for eventParam1 1\n// header comment for eventParam1 2"}: &model.DefElt{
									Val:  "String",
									Line: 22,
								},
							},
							model.DefElt{Val: "TestEvent2", Line: 34, Comment: "// line comment for TestEvent2 1"}: &model.DefMap{
								model.DefElt{Val: "eventParam2", Line: 38, Comment: "// header comment for eventParam2 1\n// header comment for eventParam2 2"}: &model.DefElt{
									Val:  "String",
									Line: 41,
								},
							},
						},
						Structs: model.DefMapMap{
							model.DefElt{Val: "TestStruct", Line: 49, Comment: "// comment for TestStruct 1"}: &model.DefMap{
								model.DefElt{Val: "x", Line: 54, Comment: "// comment for x 1\n// comment for x 2"}: &model.DefElt{
									Val:  "Int32",
									Line: 54,
								},
								model.DefElt{Val: "y", Line: 57, Comment: "// comment for y 1"}: &model.DefElt{
									Val:  "Int32",
									Line: 58,
								},
							},
						},
					},
				},
			}
		},
		"successfully test3": func(_ *testing.T) test {
			return test{
				args: args{
					path: "testdata/test3.yaml",
				},
				wants: wants{
					out: &model.SchemaDef{
						Name: model.DefElt{
							Val:  "SchemaComment",
							Line: 1,
						},
						Description: model.DefElt{
							Val:  "test description",
							Line: 2,
						},
						Events: model.DefMapMap{
							model.DefElt{Val: "TestEvent1", Line: 7}: &model.DefMap{
								model.DefElt{Val: "eventParam11", Line: 8}: &model.DefElt{
									Val:  "String",
									Line: 8,
								},
							},
							model.DefElt{Val: "TestEvent2", Line: 9}: &model.DefMap{
								model.DefElt{Val: "eventParam21", Line: 10}: &model.DefElt{
									Val:  "String",
									Line: 10,
								},
								model.DefElt{Val: "eventParam22", Line: 11}: &model.DefElt{
									Val:  "String",
									Line: 11,
								},
							},
						},
						Structs: model.DefMapMap{
							model.DefElt{Val: "TestStruct1", Line: 16}: &model.DefMap{
								model.DefElt{Val: "x", Line: 17}: &model.DefElt{
									Val:  "Int32",
									Line: 17,
								},
								model.DefElt{Val: "y", Line: 18}: &model.DefElt{
									Val:  "Int32",
									Line: 18,
								},
							},
							model.DefElt{Val: "TestStruct2", Line: 20}: &model.DefMap{
								model.DefElt{Val: "x", Line: 21}: &model.DefElt{
									Val:  "Int32",
									Line: 22,
								},
							},
						},
						Funcs: model.FuncDefMap{
							model.DefElt{Val: "TestFunc1", Line: 25}: &model.FuncDef{
								Line: 25,
								Access: model.DefElt{
									Val:  "owner",
									Line: 26,
								},
								Params: model.DefMap{
									model.DefElt{Val: "name", Line: 28}: &model.DefElt{
										Val:  "String",
										Line: 28,
									},
									model.DefElt{Val: "value", Line: 29}: &model.DefElt{
										Val:  "String",
										Line: 29,
									},
								},
								Results: model.DefMap{
									model.DefElt{Val: "length", Line: 31}: &model.DefElt{
										Val:  "Uint32",
										Line: 31,
									},
								},
							},
							model.DefElt{Val: "TestFunc2", Line: 32}: &model.FuncDef{
								Line: 32,
								Access: model.DefElt{
									Val:  "owner",
									Line: 33,
								},
								Params: model.DefMap{
									model.DefElt{Val: "name", Line: 35}: &model.DefElt{
										Val:  "String",
										Line: 35,
									},
									model.DefElt{Val: "value", Line: 36}: &model.DefElt{
										Val:  "String",
										Line: 36,
									},
								},
								Results: model.DefMap{
									model.DefElt{Val: "length", Line: 38}: &model.DefElt{
										Val:  "Uint32",
										Line: 38,
									},
								},
							},
						},
						Views: model.FuncDefMap{
							model.DefElt{Val: "TestView1", Line: 42}: &model.FuncDef{
								Line: 42,
								Access: model.DefElt{
									Val:  "owner",
									Line: 43,
								},
								Params: model.DefMap{
									model.DefElt{Val: "name", Line: 45}: &model.DefElt{
										Val:  "String",
										Line: 45,
									},
									model.DefElt{Val: "id", Line: 46}: &model.DefElt{
										Val:  "Int32",
										Line: 46,
									},
								},
								Results: model.DefMap{
									model.DefElt{Val: "length", Line: 48}: &model.DefElt{
										Val:  "Uint32",
										Line: 48,
									},
								},
							},
							model.DefElt{Val: "TestView2", Line: 49}: &model.FuncDef{
								Line: 49,
								Access: model.DefElt{
									Val:  "owner",
									Line: 50,
								},
								Params: model.DefMap{
									model.DefElt{Val: "name", Line: 52}: &model.DefElt{
										Val:  "String",
										Line: 52,
									},
									model.DefElt{Val: "id", Line: 53}: &model.DefElt{
										Val:  "Int32",
										Line: 53,
									},
								},
								Results: model.DefMap{
									model.DefElt{Val: "length", Line: 55}: &model.DefElt{
										Val:  "Uint32",
										Line: 55,
									},
								},
							},
						},
						Typedefs: model.DefMap{
							model.DefElt{Val: "TestTypedef", Line: 58}: &model.DefElt{
								Val:  "String[]",
								Line: 58,
							},
						},
						State: model.DefMap{
							model.DefElt{Val: "TestState", Line: 64}: &model.DefElt{
								Val:  "Int64[]",
								Line: 64,
							},
						},
					},
				},
			}
		},
		"successfully test4": func(_ *testing.T) test {
			return test{
				args: args{
					path: "testdata/test4.yaml",
				},
				wants: wants{
					out: &model.SchemaDef{
						Copyright: "// This is the testing copyright message\n// copyright message is the comment block\n// ahead of copyright yaml tag.\n// No other yaml item should exist ahead of\n// the copyright message, and the value of\n// copyright tag should leave empty.",
						Name: model.DefElt{
							Val:  "SchemaComment",
							Line: 8,
						},
						Description: model.DefElt{
							Val:  "test description",
							Line: 9,
						},
					},
				},
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			file, err := os.Open(tt.args.path)
			require.NoError(t, err)
			in, err := io.ReadAll(file)
			require.NoError(t, err)

			def := &model.SchemaDef{}
			root := yaml.Parse(in)
			require.NotNil(t, root)
			err = yaml.Convert(root, def)
			require.NoError(t, err)
			require.Equal(t, tt.wants.out.Name, def.Name)
			require.Equal(t, tt.wants.out.Description, def.Description)
			require.Equal(t, tt.wants.out.Events, def.Events)
			require.Equal(t, tt.wants.out.Structs, def.Structs)
			require.Equal(t, tt.wants.out.Typedefs, def.Typedefs)
			require.Equal(t, tt.wants.out.State, def.State)
			require.Equal(t, tt.wants.out.Funcs, def.Funcs)
			require.Equal(t, tt.wants.out.Views, def.Views)
		})
	}
}
