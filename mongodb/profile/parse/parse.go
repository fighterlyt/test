package parse

import (
	"bytes"
	"fmt"
)

/*本模块代码是将使用 bson 描述的mongo 数据转换为 struct 定义,相比于 json,bson多了一些类型
*	Timestamp
*	BinData
*	NumberLong
*	ISODate
基本思路,将所有的{理解为一个新的结构,而[]理解为数组
*/

type FieldKind int

const (
	SimpleKind FieldKind = iota
	SliceKind
	StructKind
)

type Struct struct {
}

type Element struct {
	Name  string
	tag   string
	Index int
	Data  Field
}
type Field interface {
	Kind() FieldKind
}

type Slice struct {
	elemTypeStr string
}

func (s Slice) Kind() FieldKind {
	return SliceKind
}

type SimpleField struct {
	typeStr string
}

func (s SimpleField) Kind() FieldKind {
	return SimpleKind
}

type StructField struct {
	typeStr *Struct
}

func (s StructField) Kind() FieldKind {
	return StructKind
}

func Parse(data []byte) (*Struct, error) {
	if bytes.Count(data, []byte("{")) != bytes.Count(data, []byte(`}`)) {
		return nil, fmt.Errorf("大括号不匹配")
	}
	if bytes.Count(data, []byte("[")) != bytes.Count(data, []byte(`]`)) {
		return nil, fmt.Errorf("方括号不匹配")
	}

}

//解析单个对象描述.不包括开始结束的{}
func pasrseStruct(data []byte) (*Struct, error) {
	start := 0
	semi := bytes.Index(data[start:], []byte(":")) //没有冒号
	if semi == -1 {
		return nil, nil
	}

}
