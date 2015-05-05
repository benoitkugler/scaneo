package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

const (
	srcFile = "testdata/structs.go"
)

type structToken struct {
	Name   string
	Fields []string
	Types  []string
}

func main() {
	structToks := make([]structToken, 0, 8)

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, srcFile, nil, 0)
	if err != nil {
		log.Println("'syntax error' - parser probably")
		log.Fatal(err)
	}

	// ast.Print(fset, astf)
	for _, dec := range astf.Decls {
		structTok := structToken{
			Fields: make([]string, 0, 8),
			Types:  make([]string, 0, 8),
		}

		genDec, isGenDec := dec.(*ast.GenDecl)
		if !isGenDec {
			continue
		}

		for _, spec := range genDec.Specs {
			typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
			if !isTypeSpec {
				continue
			}

			structTok.Name = typeSpec.Name.Name

			structType, isStructType := typeSpec.Type.(*ast.StructType)
			if !isStructType {
				continue
			}

			for _, field := range structType.Fields.List {
				for _, ident := range field.Names {
					structTok.Fields = append(structTok.Fields, ident.Name)
				}

				switch fieldType := field.Type.(type) {
				case *ast.Ident:
					structTok.Types = append(structTok.Types, fieldType.Name)
				case *ast.SelectorExpr:
					ident, isIdent := fieldType.X.(*ast.Ident)
					if !isIdent {
						continue
					}

					structTok.Types = append(structTok.Types,
						fmt.Sprintf("%s.%s", ident.Name, fieldType.Sel.Name))
				}

			}

			structToks = append(structToks, structTok)
		}
	}

	for _, sd := range structToks {
		fmt.Println(sd.Name)

		if len(sd.Fields) != len(sd.Types) {
			log.Println(sd)
			log.Fatal("Ahhh!!")
		}

		for i, _ := range sd.Fields {
			fmt.Println("   ", sd.Fields[i], sd.Types[i])
		}
	}
}
