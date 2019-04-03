package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fighterlyt/gocui"
	"github.com/fighterlyt/test/logView/model"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	ssh    = ""
	path   = ""
	text   *gocui.View
	err    error
	data   *model.Data
	reader *bufio.Reader
	height int
	event  *gocui.View
)

func read() {
	section := model.NewSection(height)
	for i := 0; i < height; i++ {
		if line, err := reader.ReadBytes('\n'); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(text, "读取失败"+err.Error())
		} else {

			element := model.NewElement(string(line))
			section.AddElement(element)
		}
	}
	data.Add(section)
	text.Clear()
	fmt.Fprintln(text, section.String())
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	ch := make(chan string, 1)
	text, err = g.SetView("log", 0, 0, maxX-10, maxY)
	text.Wrap = true
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

	}
	text.Title = "日志"
	event, err = g.SetView("event", maxX-10, 0, maxX, maxY)
	event.Title = "事件"
	_, height = text.Size()
	if file, err := os.Open(path); err != nil {
		fmt.Fprintf(text, "打开文件失败"+err.Error())
	} else {
		reader = bufio.NewReader(file)
		data = model.NewData(height, ch)

		read()

	}
	return nil
}

func next(g *gocui.Gui, v *gocui.View) error {
	contents := data.Get(data.GetSectionIndex() + 1)
	if contents != nil {
		data.SetSectionIndex(data.GetSectionIndex())
		text.Clear()
		fmt.Fprintln(text, contents.String())
		//info.SetText(data.String() + "buffer")

	} else {
		//info.SetText(data.String() + "read")
		read()
	}
	return nil
}
func previous(g *gocui.Gui, v *gocui.View) error {
	panic("previous")
	if g==nil &&  v==nil{
		panic("nil")
	}
	event.Clear()
	fmt.Fprintln(event,gocui.KeyPgup)
	index := data.GetSectionIndex() - 1
	contents := data.Get(index)
	if contents != nil {
		data.SetSectionIndex(index + 1)
		fmt.Fprintln(v, contents.String())

		//info.SetText(data.String() + key.Name())
	} else {
		//info.SetText(data.String() + key.Name())

		fmt.Fprintln(v, "未找到"+strconv.Itoa(index))
	}
	return nil
}
func main() {
	declFlag()
	flag.Parse()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//g.Highlight = true
	//g.SelFgColor = gocui.ColorRed

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, previous); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, next); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func declFlag() {
	flag.StringVar(&ssh, "ssh", "", "ssh服务器地址")
	flag.StringVar(&path, "path", "", "日志文件地址")
}

func convert(data [][]byte) string {
	strs := make([]string, len(data))
	for _, element := range data {
		strs = append(strs, string(element))
	}
	return strings.Join(strs, "\n")
}
