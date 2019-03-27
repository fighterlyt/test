package main

import (
	"bufio"
	"flag"
	"github.com/fighterlyt/test/logView/model"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	ssh  = ""
	path = ""
)

func main() {
	declFlag()
	flag.Parse()

	app := tview.NewApplication()
	flex := tview.NewFlex()
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	text.SetTitle("日志")
	text.SetBorder(true)
	flex.AddItem(text, 0, 1, true)

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	flex.AddItem(info, 30, 1, false)
	_, _, _, height := text.GetRect()

	console := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	console.SetBorder(true)
	flex.AddItem(console, 30, 20, false)
	ch := make(chan string,  1)
	go func() {
		for data := range ch {
			console.SetText(data)
		}
	}()
	info.SetTitle("信息")
	info.SetBorder(true)
	if file, err := os.Open(path); err != nil {
		text.SetText("打开文件失败" + err.Error())
	} else {
		reader := bufio.NewReader(file)
		data := model.NewData(height,ch)

		read := func() {
			section := model.NewSection(height)
			for i := 0; i < height; i++ {
				if data, err := reader.ReadBytes('\n'); err != nil {
					if err == io.EOF {
						break
					}
					text.SetText("读取失败" + err.Error())
				} else {
					element := model.NewElement(string(data))
					section.AddElement(element)
				}
			}
			data.Add(section)
			text.SetText(section.String())
		}
		read()
		text.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
			info.SetTitle(key.Name())
			switch key.Key() {
			case tcell.KeyPgDn:
				contents := data.Get(data.GetSectionIndex())
				if contents != nil {
					data.SetSectionIndex(data.GetSectionIndex() + 1)
					text.SetText(contents.String())
					info.SetText(data.String() + "buffer")

				} else {
					info.SetText(data.String() + "read")
					read()

				}

			case tcell.KeyPgUp:

				index:=data.GetSectionIndex() - 2
				contents := data.Get(index)
				if contents != nil {
					data.SetSectionIndex(index+1)

					info.SetText(data.String() + key.Name())
					text.SetText(contents.String())
				} else {
					text.SetText("未找到" + strconv.Itoa(index))
				}
			default:

				return key
			}

			return nil
		})
	}

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

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
