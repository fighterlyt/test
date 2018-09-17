package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	Parse("/Users/mac/go/src/PositionMonitor/models")
}
func Parse(dir string) {
	fset := token.NewFileSet()

	f, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.AllErrors|parser.ParseComments)

	strucs := make(map[string]*ast.TypeSpec, 100)
	commnents:=make(map[string][]*ast.CommentGroup)
	//methods:=make(map[string]ast.Decl,100)
	if err != nil {
		panic(err.Error())
	} else {

		for key, p := range f {
			fmt.Println(key, p.Name, p.Scope.String())
			if p.Scope != nil {
				for f, object := range p.Scope.Objects {
					fmt.Println("scope", f, object)
				}
			}

			for f, object := range p.Imports {
				fmt.Println("导入", f, object.Kind)
			}
			for _, file := range p.Files {
				cmap := ast.NewCommentMap(fset, file, file.Comments)
				fmt.Println(cmap)
				for _, decl := range file.Decls {
					switch decl.(type) {
					case *ast.GenDecl:
						t := decl.(*ast.GenDecl)
						if t.Tok == token.TYPE {
							for _, spec := range t.Specs {
								s := spec.(*ast.TypeSpec)

								if t,ok:=s.Type.(*ast.StructType);ok{
									strucs[s.Name.String()] = s
									commnents[s.Name.String()]=cmap[t]
								}
							}

						}
					}
				}
			}
		}
		fmt.Println(fset)
		for k,spec := range strucs {
			comment:=""
			for _,c:=range commnents[k]{
				if len(comment)==0{
					comment=c.Text()
				}else{
					comment+="\n"+c.Text()

				}
			}
			fmt.Println(k,comment)

			for _,field:=range spec.Type.(*ast.StructType).Fields.List{
				spew.Dump(field)
				if field.Comment!=nil{
					fmt.Println(field.Comment.Text())
				}
				if field.Tag!=nil{
					fmt.Println(field.Tag.Value)
				}
			}
		}
	}
}
