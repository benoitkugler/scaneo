package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

var (
	testFiles = []string{
		"testdata/access.go",
		"testdata/declarations.go",
		"testdata/methods.go",
		"testdata/types.go",
	}

	testFilesLen = len(testFiles)

	fileStructsMap = map[string][]structToken{
		"testdata/access.go": []structToken{
			{
				Name: "Exported",
				Fields: []fieldToken{
					{Name: "A", Type: "int"},
					{Name: "B", Type: "int"},
				},
			},
			{
				Name: "unexported",
				Fields: []fieldToken{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "int"},
				},
			},
			{
				Name: "ExAndUn",
				Fields: []fieldToken{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "int"},
				},
			},
			{
				Name: "unAndEx",
				Fields: []fieldToken{
					{Name: "A", Type: "int"},
					{Name: "B", Type: "int"},
				},
			},
		},
		"testdata/declarations.go": []structToken{
			{
				Name: "t0",
				Fields: []fieldToken{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "bool"},
				},
			},
			{
				Name: "t1",
				Fields: []fieldToken{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "bool"},
				},
			},
			{
				Name: "t2",
				Fields: []fieldToken{
					{Name: "a", Type: "string"},
					{Name: "b", Type: "byte"},
				},
			},
			{
				Name: "t3",
				Fields: []fieldToken{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "int"},
					{Name: "c", Type: "int"},
					{Name: "d", Type: "bool"},
					{Name: "e", Type: "bool"},
					{Name: "f", Type: "bool"},
				},
			},
			{
				Name: "t4",
				Fields: []fieldToken{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "bool"},
				},
			},
		},
		"testdata/methods.go": []structToken{
			{
				Name: "Post",
				Fields: []fieldToken{
					{Name: "ID", Type: "int"},
					{Name: "SemURL", Type: "string"},
					{Name: "Created", Type: "time.Time"},
					{Name: "Modified", Type: "time.Time"},
					{Name: "Published", Type: "pq.NullTime"},
					{Name: "Draft", Type: "bool"},
					{Name: "Title", Type: "string"},
					{Name: "Body", Type: "string"},
				},
			},
		},
		"testdata/types.go": []structToken{
			{
				Name: "boolean",
				Fields: []fieldToken{
					{Name: "a", Type: "bool"},
				},
			},
			{
				Name: "numerics",
				Fields: []fieldToken{
					{Name: "a", Type: "uint8"},
					{Name: "b", Type: "uint16"},
					{Name: "c", Type: "uint32"},
					{Name: "d", Type: "uint64"},
					{Name: "e", Type: "int8"},
					{Name: "f", Type: "int16"},
					{Name: "g", Type: "int32"},
					{Name: "h", Type: "int64"},
					{Name: "i", Type: "float32"},
					{Name: "j", Type: "float64"},
					{Name: "k", Type: "complex64"},
					{Name: "l", Type: "complex128"},
					{Name: "m", Type: "byte"},
					{Name: "n", Type: "rune"},
					{Name: "o", Type: "uint"},
					{Name: "p", Type: "int"},
					{Name: "q", Type: "uintptr"},
				},
			},
			{
				Name: "str",
				Fields: []fieldToken{
					{Name: "a", Type: "string"},
				},
			},
			{
				Name: "structs",
				Fields: []fieldToken{
					{Name: "a", Type: "sql.NullString"},
				},
			},
			{
				Name: "slices",
				Fields: []fieldToken{
					{Name: "a", Type: "[]bool"},
					{Name: "b", Type: "[]time.Time"},
					{Name: "c", Type: "[]*byte"},
					{Name: "d", Type: "[]*sql.NullString"},
				},
			},
			{
				Name: "pointers",
				Fields: []fieldToken{
					{Name: "a", Type: "*bool"},
					{Name: "b", Type: "*time.Time"},
					{Name: "c", Type: "*[]byte"},
					{Name: "d", Type: "*[]sql.NullString"},
				},
			},
		},
	}
)

func TestFilenames(t *testing.T) {
	inputPaths := []string{"testdata/", testFiles[0]}

	files, err := filenames(inputPaths)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if testFilesLen != len(files) {
		t.Error("unexpected file count")
		t.Errorf("expected: %d; found: %d\n", testFilesLen, len(files))
		t.FailNow()
	}

	sort.Strings(files)

	for i := range files {
		filename := filepath.Base(files[i])
		testFilename := filepath.Base(testFiles[i])

		if testFilename != filename {
			t.Error("unexpected filename")
			t.Errorf("expected: %s; found: %s\n", testFilename, filename)
			t.Error("files:", files)
			t.Error("testFiles:", testFiles)
			t.FailNow()
		}
	}
}

func TestWhitelist(t *testing.T) {
	whitelist := map[string]struct{}{
		"Exported":   struct{}{},
		"unexported": struct{}{},
	}
	expectedToks := len(whitelist)

	toks, err := parseCode(testFiles[0], whitelist)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expectedToks != len(toks) {
		t.Error("unexpected struct tokens length")
		t.Errorf("expected: %d; found: %d\n", expectedToks, len(toks))
	}
}

func TestParseCode(t *testing.T) {
	noFilter := make(map[string]struct{}, 0)
	for fPath, structToks := range fileStructsMap {
		// get all struct tokens for a given file
		toks, err := parseCode(fPath, noFilter)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if len(structToks) != len(toks) {
			t.Error("unexpected struct tokens length")
			t.Errorf("expected: %d; found: %d\n", len(structToks), len(toks))
			t.FailNow()
		}

		for i := range toks {
			if structToks[i].Name != toks[i].Name {
				t.Error("unexpected struct name")
				t.Errorf("expected: %s; found: %s\n", structToks[i].Name, toks[i].Name)
				t.FailNow()
			}

			if len(structToks[i].Fields) != len(toks[i].Fields) {
				t.Error("unexpected struct fields length")
				t.Error("file:", fPath)
				t.Error("struct:", structToks[i].Name)
				t.Errorf("expected: %d; found: %d\n", len(structToks[i].Fields), len(toks[i].Fields))
				t.Error("expected:", structToks[i].Fields)
				t.Error("found:", toks[i].Fields)
				t.FailNow()
			}

			for j := range toks[i].Fields {
				if structToks[i].Fields[j].Name != toks[i].Fields[j].Name {
					t.Error("unexpected struct field name")
					t.Error("file:", fPath)
					t.Error("struct:", structToks[i].Name)
					t.Errorf("expected: %s; found: %s\n", structToks[i].Fields[j].Name, toks[i].Fields[j].Name)
					t.FailNow()
				}

				if structToks[i].Fields[j].Type != toks[i].Fields[j].Type {
					t.Error("unexpected struct field type")
					t.Error("file:", fPath)
					t.Error("struct:", structToks[i].Name)
					t.Error("field:", structToks[i].Fields[j].Name)
					t.Errorf("expected: %s; found: %s\n", structToks[i].Fields[j].Type, toks[i].Fields[j].Type)
					t.FailNow()
				}
			}
		}
	}
}

func TestGenFile(t *testing.T) {
	toks := fileStructsMap["testdata/access.go"][:2]

	expectedFuncNames := []string{
		"scanExported",
		"scanUnexported",
		"scanExporteds",
		"scanUnexporteds",
	}

	fout, err := ioutil.TempFile(os.TempDir(), "scaneo-test-")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.Remove(fout.Name()) // comment this line to examin generated code
	defer fout.Close()

	// file, package, unexport, tokens
	if err := genFile(fout, "testing", true, toks); err != nil {
		t.Error(err)
		t.FailNow()
	}

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, fout.Name(), nil, 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	scanFuncs := make([]string, 0, len(toks))
	for _, dec := range astf.Decls {
		funcDecl, isFuncDecl := dec.(*ast.FuncDecl)
		if !isFuncDecl {
			continue
		}

		scanFuncs = append(scanFuncs, funcDecl.Name.String())
	}

	if len(toks)*2 != len(scanFuncs) {
		t.Error("unexpected number of scan functions found")
		t.Errorf("expected: %d; found: %d\n", len(toks)*2, len(scanFuncs))
		t.FailNow()
	}

	for i := range expectedFuncNames {
		if expectedFuncNames[i] != scanFuncs[i] {
			t.Error("unexpected scan function found")
			t.Errorf("expected: %s; found: %s\n", expectedFuncNames[i], scanFuncs[i])
		}
	}

}
