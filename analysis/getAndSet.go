package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

/* 目的是为 struct 自动生成Get/Set方法*/
var (
	fileName   = ""
	structName = ""
)

func init() {
	flag.StringVar(&fileName, "fileName", "", "fileName")
	flag.StringVar(&structName, "structName", "", "结构名")
}

var s = `package test// 保存了需要使用mongo聚合管道操作取出的数据的缓存
type MarketQuotesMongoPipeData struct {
	marketCode          string
	nextMongoPipeTimeMS int64            // 下一次使用mongo聚合管道操作取出数据的时间
	sortClosingPrice    *bson.Decimal128 // 最后一次收盘价
	yesterDayTopPrice   *bson.Decimal128 // 昨天最高价
	todayOpenningPrice  *bson.Decimal128 // 昨天开盘价
	_24HLowPrice        *bson.Decimal128 // 24小时最低价
	_24HDealCount       *bson.Decimal128 // 24小时成交量
	_24HTurnover        *bson.Decimal128 // 24小时成交额
}`

func main() {
	flag.Parse()
	if fileName==""{
		fileName="/Users/mac/go/src/radar-trade/dynamic_quote_util.go"
	}
	if structName==""{
		structName="MarketQuotesMongoPipeData"
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("文件错误" + err.Error())
	}
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fileName, string(data), parser.ParseComments)
	if err != nil {
		panic("解析错误" + err.Error())
	}
	FilterStruct(f, fset, string(data), structName)
}

func FilterStruct(file *ast.File, fset *token.FileSet, source string, structName string) {

	ast.Inspect(file, func(x ast.Node) bool {
		switch node := x.(type) {
		case *ast.GenDecl:
			switch node.Tok {
			case token.TYPE:

			case token.IMPORT:

			}
		case *ast.TypeSpec:
			if _, ok := node.Type.(*ast.StructType); ok {
				if node.Name.Name == structName {
					println(generate(node, source))
				}

			}

		}
		return true
	})

}
func generate(s *ast.TypeSpec, source string) string {
	typeName := s.Name.Name
	receiver := strings.ToLower(typeName[:1])
	builder := &strings.Builder{}
	for _, field := range s.Type.(*ast.StructType).Fields.List {
		for _, name := range field.Names {
			fieldType:=strings.TrimSpace(source[field.Type.Pos()-1:field.Type.End()])
			fmt.Fprintf(builder, "func (%s %s) Get%s () %s{\n\treturn %s.%s\n}\n", receiver, typeName, strings.Title(name.Name),fieldType, receiver, name.Name)
			fmt.Fprintf(builder, "func (%s *%s) Set%s (%s %s ) {\n\t %s.%s=%s\n}\n", receiver, typeName, strings.Title(name.Name), name.Name, fieldType, receiver, name.Name, name.Name)

		}
	}
	return builder.String()
}
