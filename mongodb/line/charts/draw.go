package charts

import (
	"github.com/chenjiandongx/go-echarts/charts"
	"net/http"
	"os"
)

func Handle(arguments []LineArguments) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		page := charts.NewPage(charts.RouterOpts{
			URL: "line",
		})
		for _, argument := range arguments {
			page.Add(baseLine(argument.Title, argument.Name, argument.Names, argument.Values))

		}
		f, err := os.Create("html/line.html")
		if err != nil {
			page.Render(w, f)
		}
	}
}

type LineArguments struct {
	Title  string
	Name   string
	Names  []string
	Values []int
}

func baseLine(title, name string, names []string, values []int) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: title})
	line.AddXAxis(names).AddYAxis(name, values)
	return line
}
