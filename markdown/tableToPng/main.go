package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os/exec"
	"strings"
)

var (
	tmpl = `<html>
<head>
    <style>
        pre{font-family:Menlo,Monaco,Consolas,'Courier New',monospace;direction:ltr;text-align:left;white-space:pre;word-spacing:normal;word-break:normal;padding:1em;margin:.5em 0;overflow:auto;line-height:1.5;tab-size:4;hyphens:none;color:#383a42;background-color:#eaeaeb !important;border:#d1d1d2;border-radius:3px}pre[class*="language-"]{padding:1em}code[class*="language-"] .token.comment,pre[class*="language-"] .token.comment,code[class*="language-"] .token.prolog,pre[class*="language-"] .token.prolog,code[class*="language-"] .token.doctype,pre[class*="language-"] .token.doctype,code[class*="language-"] .token.cdata,pre[class*="language-"] .token.cdata{color:#a0a1a7}code[class*="language-"] .namespace,pre[class*="language-"] .namespace{opacity:.7}code[class*="language-"] .token.constant,pre[class*="language-"] .token.constant{color:#986801}code[class*="language-"] .token.boolean,pre[class*="language-"] .token.boolean,code[class*="language-"] .token.number,pre[class*="language-"] .token.number,code[class*="language-"] .token.function-name,pre[class*="language-"] .token.function-name{color:#986801}code[class*="language-"] .token.tag,pre[class*="language-"] .token.tag{color:#e45649}code[class*="language-"] .token.symbol,pre[class*="language-"] .token.symbol{color:#0184bc}code[class*="language-"] .token.selector,pre[class*="language-"] .token.selector{color:#a626a4}code[class*="language-"] .token.attr-name,pre[class*="language-"] .token.attr-name{color:#986801}code[class*="language-"] .token.string,pre[class*="language-"] .token.string{color:#50a14f}code[class*="language-"] .token.char,pre[class*="language-"] .token.char{color:#0184bc}code[class*="language-"] .token.url,pre[class*="language-"] .token.url{color:#0184bc}code[class*="language-"] .token.operator,pre[class*="language-"] .token.operator{color:#383a42}code[class*="language-"] .token.atrule,pre[class*="language-"] .token.atrule,code[class*="language-"] .token.attr-value,pre[class*="language-"] .token.attr-value,code[class*="language-"] .token.keyword,pre[class*="language-"] .token.keyword{color:#a626a4}code[class*="language-"] .token.function,pre[class*="language-"] .token.function{color:#4078f2}code[class*="language-"] .token.class-name,pre[class*="language-"] .token.class-name{color:#c18401}code[class*="language-"] .token.variable,pre[class*="language-"] .token.variable{color:#986801}code[class*="language-"] .token.regex,pre[class*="language-"] .token.regex{color:#0184bc}code[class*="language-"] .token.important,pre[class*="language-"] .token.important{color:#e45649}code[class*="language-"] .token.important,pre[class*="language-"] .token.important,code[class*="language-"] .token.bold,pre[class*="language-"] .token.bold{font-weight:bold}code[class*="language-"] .token.italic,pre[class*="language-"] .token.italic{font-style:italic}code[class*="language-"] .token.entity,pre[class*="language-"] .token.entity{cursor:help}pre[data-line]{position:relative;padding:1em 0 1em 3em}pre[data-line] .line-highlight-wrapper{position:absolute;top:0;left:0;background-color:transparent;display:block;width:100%%}pre[data-line] .line-highlight{position:absolute;left:0;right:0;padding:inherit 0;margin-top:1em;background:rgba(153,122,102,0.08);background:linear-gradient(to right, rgba(153,122,102,0.1) 70%%, rgba(153,122,102,0));pointer-events:none;line-height:inherit;white-space:pre}pre[data-line] .line-highlight:before,pre[data-line] .line-highlight[data-end]:after{content:attr(data-start);position:absolute;top:.4em;left:.6em;min-width:1em;padding:0 .5em;background-color:rgba(153,122,102,0.4);color:#f5f2f0;font:bold 65%%/1.5 sans-serif;text-align:center;vertical-align:.3em;border-radius:999px;text-shadow:none;box-shadow:0 1px white}pre[data-line] .line-highlight[data-end]:after{content:attr(data-end);top:auto;bottom:.4em}html body{font-family:"Helvetica Neue",Helvetica,"Segoe UI",Arial,freesans,sans-serif;font-size:16px;line-height:1.6;color:#383a42;background-color:#fafafa;overflow:initial;box-sizing:border-box;word-wrap:break-word}html body>:first-child{margin-top:0}html body h1,html body h2,html body h3,html body h4,html body h5,html body h6{line-height:1.2;margin-top:1em;margin-bottom:16px;color:#000}html body h1{font-size:2.25em;font-weight:300;padding-bottom:.3em}html body h2{font-size:1.75em;font-weight:400;padding-bottom:.3em}html body h3{font-size:1.5em;font-weight:500}html body h4{font-size:1.25em;font-weight:600}html body h5{font-size:1.1em;font-weight:600}html body h6{font-size:1em;font-weight:600}html body h1,html body h2,html body h3,html body h4,html body h5{font-weight:600}html body h5{font-size:1em}html body h6{color:#5e616e}html body strong{color:#000}html body del{color:#5e616e}html body a:not([href]){color:inherit;text-decoration:none}html body a{color:#0184bc;text-decoration:none}html body a:hover{color:#01a0e5;text-decoration:none}html body img{max-width:100%%}html body>p{margin-top:0;margin-bottom:16px;word-wrap:break-word}html body>ul,html body>ol{margin-bottom:16px}html body ul,html body ol{padding-left:2em}html body ul.no-list,html body ol.no-list{padding:0;list-style-type:none}html body ul ul,html body ul ol,html body ol ol,html body ol ul{margin-top:0;margin-bottom:0}html body li{margin-bottom:0}html body li.task-list-item{list-style:none}html body li>p{margin-top:0;margin-bottom:0}html body .task-list-item-checkbox{margin:0 .2em .25em -1.8em;vertical-align:middle}html body .task-list-item-checkbox:hover{cursor:pointer}html body blockquote{margin:16px 0;font-size:inherit;padding:0 15px;color:#5e616e;border-left:4px solid #d1d1d2}html body blockquote>:first-child{margin-top:0}html body blockquote>:last-child{margin-bottom:0}html body hr{height:4px;margin:32px 0;background-color:#d1d1d2;border:0 none}html body table{margin:10px 0 15px 0;border-collapse:collapse;border-spacing:0;display:block;width:100%%;overflow:auto;word-break:normal;word-break:keep-all}html body table th{font-weight:bold;color:#000}html body table td,html body table th{border:1px solid #d1d1d2;padding:6px 13px}html body dl{padding:0}html body dl dt{padding:0;margin-top:16px;font-size:1em;font-style:italic;font-weight:bold}html body dl dd{padding:0 16px;margin-bottom:16px}html body code{font-family:Menlo,Monaco,Consolas,'Courier New',monospace;font-size:.85em !important;color:#000;background-color:#eaeaeb;border-radius:3px;padding:.2em 0}html body code::before,html body code::after{letter-spacing:-0.2em;content:"\00a0"}html body pre>code{padding:0;margin:0;font-size:.85em !important;word-break:normal;white-space:pre;background:transparent;border:0}html body .highlight{margin-bottom:16px}html body .highlight pre,html body pre{padding:1em;overflow:auto;font-size:.85em !important;line-height:1.45;border:#d1d1d2;border-radius:3px}html body .highlight pre{margin-bottom:0;word-break:normal}html body pre code,html body pre tt{display:inline;max-width:initial;padding:0;margin:0;overflow:initial;line-height:inherit;word-wrap:normal;background-color:transparent;border:0}html body pre code:before,html body pre tt:before,html body pre code:after,html body pre tt:after{content:normal}html body p,html body blockquote,html body ul,html body ol,html body dl,html body pre{margin-top:0;margin-bottom:16px}html body kbd{color:#000;border:1px solid #d1d1d2;border-bottom:2px solid #c1c1c2;padding:2px 4px;background-color:#eaeaeb;border-radius:3px}@media print{html body{background-color:#fafafa}html body h1,html body h2,html body h3,html body h4,html body h5,html body h6{color:#000;page-break-after:avoid}html body blockquote{color:#5e616e}html body pre{page-break-inside:avoid}html body table{display:table}html body img{display:block;max-width:100%%;max-height:100%%}html body pre,html body code{word-wrap:break-word;white-space:pre}}.markdown-preview{width:100%%;height:100%%;box-sizing:border-box}.markdown-preview .pagebreak,.markdown-preview .newpage{page-break-before:always}.markdown-preview pre.line-numbers{position:relative;padding-left:3.8em;counter-reset:linenumber}.markdown-preview pre.line-numbers>code{position:relative}.markdown-preview pre.line-numbers .line-numbers-rows{position:absolute;pointer-events:none;top:1em;font-size:100%%;left:0;width:3em;letter-spacing:-1px;border-right:1px solid #999;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none}.markdown-preview pre.line-numbers .line-numbers-rows>span{pointer-events:none;display:block;counter-increment:linenumber}.markdown-preview pre.line-numbers .line-numbers-rows>span:before{content:counter(linenumber);color:#999;display:block;padding-right:.8em;text-align:right}.markdown-preview .mathjax-exps .MathJax_Display{text-align:center !important}.markdown-preview:not([for="preview"]) .code-chunk .btn-group{display:none}.markdown-preview:not([for="preview"]) .code-chunk .status{display:none}.markdown-preview:not([for="preview"]) .code-chunk .output-div{margin-bottom:16px}.scrollbar-style::-webkit-scrollbar{width:8px}.scrollbar-style::-webkit-scrollbar-track{border-radius:10px;background-color:transparent}.scrollbar-style::-webkit-scrollbar-thumb{border-radius:5px;background-color:rgba(150,150,150,0.66);border:4px solid rgba(150,150,150,0.66);background-clip:content-box}html body[for="html-export"]:not([data-presentation-mode]){position:relative;width:100%%;height:100%%;top:0;left:0;margin:0;padding:0;overflow:auto}html body[for="html-export"]:not([data-presentation-mode]) .markdown-preview{position:relative;top:0}@media screen and (min-width:914px){html body[for="html-export"]:not([data-presentation-mode]) .markdown-preview{padding:2em calc(50%% - 457px + 2em)}}@media screen and (max-width:914px){html body[for="html-export"]:not([data-presentation-mode]) .markdown-preview{padding:2em}}@media screen and (max-width:450px){html body[for="html-export"]:not([data-presentation-mode]) .markdown-preview{font-size:14px !important;padding:1em}}@media print{html body[for="html-export"]:not([data-presentation-mode]) #sidebar-toc-btn{display:none}}html body[for="html-export"]:not([data-presentation-mode]) #sidebar-toc-btn{position:fixed;bottom:8px;left:8px;font-size:28px;cursor:pointer;color:inherit;z-index:99;width:32px;text-align:center;opacity:.4}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] #sidebar-toc-btn{opacity:1}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc{position:fixed;top:0;left:0;width:300px;height:100%%;padding:32px 0 48px 0;font-size:14px;box-shadow:0 0 4px rgba(150,150,150,0.33);box-sizing:border-box;overflow:auto;background-color:inherit}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc::-webkit-scrollbar{width:8px}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc::-webkit-scrollbar-track{border-radius:10px;background-color:transparent}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc::-webkit-scrollbar-thumb{border-radius:5px;background-color:rgba(150,150,150,0.66);border:4px solid rgba(150,150,150,0.66);background-clip:content-box}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc a{text-decoration:none}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc ul{padding:0 1.6em;margin-top:.8em}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc li{margin-bottom:.8em}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .md-sidebar-toc ul{list-style-type:none}html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .markdown-preview{left:300px;width:calc(100%% -  300px);padding:2em calc(50%% - 457px -  150px);margin:0;box-sizing:border-box}@media screen and (max-width:1274px){html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .markdown-preview{padding:2em}}@media screen and (max-width:450px){html body[for="html-export"]:not([data-presentation-mode])[html-show-sidebar-toc] .markdown-preview{width:100%%}}html body[for="html-export"]:not([data-presentation-mode]):not([html-show-sidebar-toc]) .markdown-preview{left:50%%;transform:translateX(-50%%)}html body[for="html-export"]:not([data-presentation-mode]):not([html-show-sidebar-toc]) .md-sidebar-toc{display:none}
        /* Please visit the URL below for more information: */
        /*   https://shd101wyy.github.io/markdown-preview-enhanced/#/customize-css */

    </style>
</head>
<body>
%s
</body>
</html>`
)

func main() {
	data := []byte(`|  表头   | 表头  |
|  ----  | ----  |
| 单元格  | 单元格 |
| 单元格  | 单元格 |`)
	output := blackfriday.Run(data)

	file, _ := ioutil.TempFile("", "md*.html")
	fmt.Fprintf(file, tmpl, string(output))
	file.Close()
	c := ImageOptions{Input: file.Name(), Format: "png", Output: "/tmp/table.png", BinaryPath: "wkhtmltoimage"}

	if _, err := GenerateImage(&c); err != nil {
		panic(err.Error())
	}
}

// ImageOptions represent the options to generate the image.
type ImageOptions struct {
	// BinaryPath the path to your wkhtmltoimage binary. REQUIRED
	//
	// Must be absolute path e.g /usr/local/bin/wkhtmltoimage
	BinaryPath string
	// Input is the content to turn into an image. REQUIRED
	//
	// Can be a url (http://example.com), a local file (/tmp/example.html), or html as a string (send "-" and set the Html value)
	Input string
	// Format is the type of image to generate
	//
	// jpg, png, svg, bmp supported. Defaults to local wkhtmltoimage default
	Format string
	// Height is the height of the screen used to render in pixels.
	//
	// Default is calculated from page content. Default 0 (renders entire page top to bottom)
	Height int
	// Width is the width of the screen used to render in pixels.
	//
	// Note that this is used only as a guide line. Default 1024
	Width int
	// Quality determines the final image quality.
	//
	// Values supported between 1 and 100. Default is 94
	Quality int
	// Html is a string of html to render into and image.
	//
	// Only needed to be set if Input is set to "-"
	Html string
	// Output controls how to save or return the image.
	//
	// Leave nil to return a []byte of the image. Set to a path (/tmp/example.png) to save as a file.
	Output string
}

// GenerateImage creates an image from an input.
// It returns the image ([]byte) and any error encountered.
func GenerateImage(options *ImageOptions) ([]byte, error) {
	arr, err := buildParams(options)
	if err != nil {
		return []byte{}, err
	}

	if options.BinaryPath == "" {
		return []byte{}, errors.New("BinaryPath not set")
	}

	cmd := exec.Command(options.BinaryPath, arr...)

	if options.Html != "" {
		cmd.Stdin = strings.NewReader(options.Html)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	println(len(output), strings.Join(arr, " "))

	// trimmed := cleanupOutput(output, options.Format)

	return nil, err
}

// buildParams takes the image options set by the user and turns them into command flags for wkhtmltoimage
// It returns an array of command flags.
func buildParams(options *ImageOptions) ([]string, error) {
	a := []string{}

	if options.Input == "" {
		return []string{}, errors.New("Must provide input")
	}

	// silence extra wkhtmltoimage output
	// might want to add --javascript-delay too?
	// a = append(a, "-q")
	// a = append(a, "--disable-plugins")
	//
	// a = append(a, "--format")
	// if options.Format != "" {
	// 	a = append(a, options.Format)
	// } else {
	// 	a = append(a, "png")
	// }
	//
	// if options.Height != 0 {
	// 	a = append(a, "--height")
	// 	a = append(a, strconv.Itoa(options.Height))
	// }
	//
	// if options.Width != 0 {
	// 	a = append(a, "--width")
	// 	a = append(a, strconv.Itoa(options.Width))
	// }
	//
	// if options.Quality != 0 {
	// 	a = append(a, "--quality")
	// 	a = append(a, strconv.Itoa(options.Quality))
	// }
	a = append(a, "--encoding")
	a = append(a, "utf-8")

	// url and output come last
	if options.Input != "-" {
		// make sure we dont pass stdin if we aren't expecting it
		options.Html = ""
	}

	a = append(a, options.Input)

	if options.Output == "" {
		a = append(a, "-")
	} else {
		a = append(a, options.Output)
	}
	return a, nil
}

func cleanupOutput(img []byte, format string) []byte {
	buf := new(bytes.Buffer)
	switch {
	case format == "png":
		decoded, err := png.Decode(bytes.NewReader(img))
		for err != nil {
			img = img[1:]
			if len(img) == 0 {
				break
			}
			decoded, err = png.Decode(bytes.NewReader(img))
		}
		png.Encode(buf, decoded)
		return buf.Bytes()
	case format == "jpg":
		decoded, err := jpeg.Decode(bytes.NewReader(img))
		for err != nil {
			img = img[1:]
			if len(img) == 0 {
				break
			}
			decoded, err = jpeg.Decode(bytes.NewReader(img))
		}
		jpeg.Encode(buf, decoded, nil)
		return buf.Bytes()
		// case format == "svg":
		// 	return img
	}
	return img
}
