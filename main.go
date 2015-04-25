package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	// "golang.org/x/text/width"
)

const root = "http://mygodoc.oschina.mopaas.com/"

// "http://mygodoc.oschina.mopaas.com/pkg/archive_tar.htm"
// "https://godoc.org/net"

func main() {
	urls, err := index()
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(urls)
	//return
	var url string
	var b []byte
	var f *os.File
	for _, url = range urls {
		b, err = gofmt(root + "pkg/" + url + ".htm")
		if err != nil {
			fmt.Println(url, err)
			err = nil
			continue
		}

		path := strings.Replace(url, "_", "/", -1)
		err = os.MkdirAll(path, 0777)
		if err != nil {
			break
		}

		f, err = os.Create(path + "/doc_zh_CN.go")
		if err == nil {
			_, err = f.Write(b)
			f.Close()
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		fmt.Println(url, err)
	}
}

func gofmt(url string) (b []byte, err error) {
	var buf bytes.Buffer
	m, err := docgo(url)
	if err != nil {
		return
	}

	for _, txt := range m["PDoc"] {
		buf.WriteString(txt)
		buf.WriteString("\n")
	}

	buf.WriteString("package " + m["Name"][0])
	buf.WriteString("\n")
	for _, txt := range m["Consts"] {
		buf.WriteString(txt)
		buf.WriteString("\n")
	}
	for _, txt := range m["Vars"] {
		buf.WriteString(txt)
		buf.WriteString("\n")
	}
	for _, txt := range m["Decl"] {
		buf.WriteString(txt)
		buf.WriteString("\n")
	}

	b, err = format.Source(buf.Bytes())
	return
}

//
// #stdlib ~ .dir a[href^="./"]

func index() (urls []string, err error) {
	doc, err := goquery.NewDocument(root)
	if err != nil {
		return
	}

	doc.Find(`#stdlib ~ .dir a[href^="./"]`).Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			urls = append(urls, href[6:len(href)-4])
		}
	})
	return urls, nil
}

// pre 标签注释, 缩进格式化, 行首前导 "// "
func preText(text string) string {
	// pre 开始没有回车, 结尾常有一个回车
	return "// \t" + strings.Replace(text, "\n", "\n// \t", -1)
}

// 注释, 行首前导 "// "
func toText(n *goquery.Selection) (text string) {
	text = strings.Trim(n.Text(), "\n")
	if n.Is("pre") {
		return preText(text)
	}
	if n.Is("p") {
		return "// " + strings.Replace(text, "\n", "\n// ", -1)
	}
	return ""
}

func docgo(url string) (m map[string][]string, err error) {
	// doc, err := goquery.NewDocument(root + "/pkg/" + url + ".htm")
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}

	// package name
	nodes := doc.Find(`#pkg-overview`)
	text := strings.TrimPrefix(nodes.Text(), "package ")
	if text == "" {
		return nil, fmt.Errorf("losted package name: %s", url)
	}
	m = map[string][]string{}
	m["Name"], text = []string{text}, ""
	n := nodes.Next()

	// skip import "package/path"
	if strings.HasPrefix(n.Text(), `import "`) {
		n = n.Next()
	}

	for ; n.Is("p,pre"); n = n.Next() {

		if n.Is("h3") {
			break
		}
		txt := toText(n)
		if txt == "" {
			continue
		}
		if text == "" {
			text = txt
		} else {
			text += "\n//\n" + txt
		}
	}

	if text != "" {
		m["PDoc"], text = []string{text}, ""
	}

	// 文档中 const, var, type, func 使用 h3[id], methods 使用 h4[id].
	// 顶级定义从 h3 标签开始, 到另一个 h3 标签结束.
	// 查找第一个非 #pkg-index 的 h3 标签.
	// 其他 id: skip pkg-index, pkg-examples, pkg-files
	nodes = doc.Find("h3[id]")

	nodes.Each(func(_ int, n *goquery.Selection) {
		if err != nil {
			return
		}
		id, exists := n.Attr("id")
		if !exists {
			return
		}

		var key string
		switch id {
		case "pkg-constants":
			key = "Consts"
		case "pkg-variables":
			key = "Vars"
		default: // type, func
			if strings.HasPrefix(id, "pkg-") {
				return
			}
			// 简化声明分级算法, 适用于真实 DOC 的结构
			key = "Decl"
		}

		// h3[id] 后跟 pre, h4[id] 后必跟 pre[class] 为代码, p 为文档
		// 先提取 h3 级别的代码.
		n = n.Next()
		if !n.Is("pre") {
			err = fmt.Errorf("expected pre: %s#%s\n", url, id)
			return
		}

		var code, comments string
		code = n.Text()

		// http://skill.events/553a7d3e64656d2c3a000000
		// 提取注释
		for {
			n = n.Next()
			// 后置 code
			if !n.Is("p") && !n.Is("pre") {
				if code != "" {
					if comments != "" {
						m[key] = append(m[key], "")
					}
					m[key] = append(m[key], comments)
					m[key] = append(m[key], code)
				}
				// 结束
				if !n.Is("h4") {
					return
				}
				code, comments = "", ""
				// h4[id] 后跟 method 定义代码
				n = n.Next()
				if !n.Is("pre.funcdecl") {
					err = fmt.Errorf("expected pre.funcdecl: %s#%s\n", url, id)
					return
				}
				code = n.Text()
				continue
			}

			// 新 code 或者连续的 code(无注释)
			if comments == "" && n.Is("pre") {
				if code != "" {
					m[key] = append(m[key], "")
					m[key] = append(m[key], code)
				}
				code = n.Text()
				continue
			}

			// 注释分段
			if comments == "" {
				comments = toText(n)
			} else {
				comments += "\n//\n" + toText(n)
			}
		}
	})

	return
}

// via github.com/golang-china/golangdoc/docgen/main.go
const tmplPackageText = `// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

{{if .PDoc}}{{/* template comments */}}{{/*

-------------------------------------------------------------------------------
-- PACKAGE DOCUMENTATION
-------------------------------------------------------------------------------

*/}}{{.PDoc}}
{{end}}package {{.Name}}
{{/*

-------------------------------------------------------------------------------
-- CONSTANTS
-------------------------------------------------------------------------------

*/}}{{with .Consts}}{{range .}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- VARIABLES
-------------------------------------------------------------------------------

*/}}{{with .Vars}}{{range .}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- FUNCTIONS
-------------------------------------------------------------------------------

*/}}{{with .Funcs}}{{range .}}
{{comment_text .Name .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES
-------------------------------------------------------------------------------

*/}}{{with .Types}}{{range .}}{{$typeName := .Name}}
{{comment_text .Name .Doc "" "\t"}}
{{node .Decl}}
{{/*

-------------------------------------------------------------------------------
-- TYPES.CONSTANTS
-------------------------------------------------------------------------------

*/}}{{if .Consts}}{{range .Consts}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.VARIABLES
-------------------------------------------------------------------------------

*/}}{{if .Vars}}{{range .Vars}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.FUNCTIONS
-------------------------------------------------------------------------------

*/}}{{if .Funcs}}{{range .Funcs}}
{{comment_text .Name .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
TYPES.METHODS
-------------------------------------------------------------------------------

*/}}{{if .Methods}}{{range .Methods}}
{{comment_text (printf "%s.%s" $typeName .Name) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.END
-------------------------------------------------------------------------------

*/}}{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- END
-------------------------------------------------------------------------------

*/}}{{end}}
`
