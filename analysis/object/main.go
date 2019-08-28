package main

import (
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dirName := "/Users/mac/projects/ants/user"
	fset := token.NewFileSet() // positions are relative to fset

	fd, err := os.Open(dirName)
	if err != nil {
		panic(err.Error())
	}
	list, err := fd.Readdir(-1)
	if err != nil {
		panic(err.Error())
	}
	structs := make(map[token.Pos]*structType, 0)
	for _, d := range list {
		if strings.Contains(d.Name(), "model.go") {
			filename := filepath.Join(dirName, d.Name())
			if src, err := parser.ParseFile(fset, filename, nil, parser.ParseComments); err == nil {
				ast.Inspect(src, collectStructs(&structs))
			}
		}
	}

	for _, s := range structs {
		start := fset.Position(s.node.Pos()).Line
		end := fset.Position(s.node.End()).Line

		ast.Inspect(s.node, getTag(fset, start, end, s))
	}
	for _, s := range structs {
		println(s.name)
		for _, tag := range s.fields {
			fmt.Printf("\t%s:%s:%s\n", tag.Name, tag.Tag, tag.T)
		}
	}

}

type structType struct {
	name   string
	node   *ast.StructType
	fields []Field
}

type Field struct {
	Name string
	Tag  string
	T    string
}

func collectStructs(structs *map[token.Pos]*structType) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		var t ast.Expr
		var structName string
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Type == nil {
				return true

			}

			structName = x.Name.Name
			t = x.Type
		case *ast.CompositeLit:
			t = x.Type
		case *ast.ValueSpec:
			structName = x.Names[0].Name
			t = x.Type
		}

		x, ok := t.(*ast.StructType)
		if !ok {
			return true
		}

		(*structs)[x.Pos()] = &structType{
			name: structName,
			node: x,
		}
		return true
	}
}
func getTag(fset *token.FileSet, start, end int, s *structType) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		x, ok := n.(*ast.StructType)
		if !ok {
			return true
		}

		for _, f := range x.Fields.List {
			line := fset.Position(f.Pos()).Line

			if !(start <= line && line <= end) {
				continue
			}

			if f.Tag == nil {
				f.Tag = &ast.BasicLit{}
			}

			fieldName := ""
			if len(f.Names) != 0 {
				fieldName = f.Names[0].Name
			}

			// anonymous field
			if f.Names == nil {
				ident, ok := f.Type.(*ast.Ident)
				if !ok {
					continue
				}

				fieldName = ident.Name
			}

			res, err := process(f.Tag.Value)
			if err == nil {
				typeName := ""
				switch e := f.Type.(type) {
				case *ast.Ident:
					typeName = e.Name
				case *ast.StarExpr:
					typeName = e.X.(*ast.Ident).Name
				}
				s.fields = append(s.fields, Field{
					Name: fieldName,
					Tag:  res,
					T:    typeName,
				})
			}

			f.Tag.Value = res
		}

		return true
	}
}
func process(tagVal string) (string, error) {
	var tag string
	if tagVal != "" {
		var err error
		tag, err = strconv.Unquote(tagVal)
		if err != nil {
			return "", err
		}
	}

	tags, err := structtag.Parse(tag)
	if err != nil {
		return "", err
	}
	tagStruct, err := tags.Get("bson")
	if err != nil {
		return "", err
	}
	return tagStruct.Name, nil

}
