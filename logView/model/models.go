package model

import (
	"fmt"
	"strings"
	"time"
)

type Data struct {
	sections     []*Section
	index        int
	upper        int
	lower        int
	sectionIndex int
	wrapped      bool
	ch           chan string
}

func NewData(count int, ch chan string) *Data {
	return &Data{
		sections: make([]*Section, count, count),
		ch:       ch,
	}
}
func (d Data) realIndex(index int) int {
	return index % len(d.sections)
}
func (d Data) GetSectionIndex() int {
	return d.sectionIndex
}
func (d *Data) SetSectionIndex(index int) {
	d.sectionIndex = index
}
func (d *Data) Add(section *Section) {
	d.sections[d.realIndex(d.upper)] = section
	if d.upper == d.lower && d.wrapped {
		d.lower = d.realIndex(d.lower + 1)

	}
	if d.upper > len(d.sections)-1 {
		d.wrapped = true
	}
	d.addUpper()

	d.index = d.realIndex(d.index + 1)

	d.sectionIndex++
	d.ch <- "add"+time.Now().String()

}
func (d *Data) addUpper() {
	d.upper = d.realIndex(d.upper + 1)

}

func (d *Data) sectionIndexLower() int {
	if d.index > d.lower {
		return d.sectionIndex - d.index + d.lower
	} else if d.index == d.lower {
		return d.sectionIndex - len(d.sections)
	} else {
		return d.sectionIndex - (d.index + len(d.sections) - d.lower)

	}
}

func (d *Data) sectionIndexUpper() int {
	if d.index <= d.upper {
		return d.sectionIndex + d.upper - d.index
	} else {
		return d.sectionIndex + d.upper - d.index + len(d.sections)
	}
}
func (d *Data) Get(index int) *Section {
	//println("beforeGet",index,d.String())
	//defer func(){
	//	println("afterGet",index,d.String())
	//}()
	lower := d.sectionIndexLower()
	upper := d.sectionIndexUpper()

	if index < upper && index >= lower {
		index = index - d.sectionIndex + d.index

		if index < 0 {
			index += len(d.sections)
		}
		return d.sections[d.realIndex(index)]
	}
	return nil
}
func (d Data) String() string {
	return fmt.Sprintf("index=%d,sectionIndex=%d,upper=%d,lower=%d\n", d.index, d.sectionIndex, d.upper, d.lower)
}

//Section 一个章节
type Section struct {
	tops  []*Element
	count int
}

func NewSection(count int) *Section {
	return &Section{
		tops: make([]*Element, 0, count),
	}
}

func (s *Section) AddElement(element *Element) {
	s.tops = append(s.tops, element)
	s.count += element.count
}

func (s Section) String() string {
	strs := make([]string, 0, len(s.tops))
	for _, top := range s.tops {
		strs = append(strs, top.data)
	}
	return strings.Join(strs, "\n")
}

//Element 一组语句
type Element struct {
	data     string     //自身的记录
	children []*Element //子节点
	count    int        //本组总长度
}

func NewElement(data string) *Element {
	return &Element{
		data:  data,
		count: len(data),
	}
}

func (e *Element) AddChild(child *Element) {
	e.children = append(e.children, child)
	e.count += child.count
}
