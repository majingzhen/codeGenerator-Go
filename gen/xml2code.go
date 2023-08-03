package gen

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"time"
)

type MDB struct {
	XMLName xml.Name `xml:"MDB"`
	Modules []Module `xml:"Module"`
	Tables  []Table  `xml:"Table"`
}

type Module struct {
	XMLName xml.Name `xml:"Module"`
	Name    string   `xml:"Name,attr"`
	Code    string   `xml:"Code,attr"`
	Tables  []Table  `xml:"Table"`
}

type Table struct {
	XMLName     xml.Name `xml:"Table"`
	Name        string   `xml:"Name,attr"`
	Code        string   `xml:"Code,attr"`
	HtmlPkg     string   `xml:"HtmlPkg,attr"`
	GoPkg       string   `xml:"GoPkg,attr"`
	DBCreator   string   `xml:"DBCreator,attr"`
	Comment     string   `xml:"Comment,attr"`
	Columns     []Column `xml:"Column"`
	Author      string
	FileName    string
	PackageName string
	StructName  string
	DateTime    time.Time
}

type Column struct {
	XMLName   xml.Name `xml:"Column"`
	Name      string   `xml:"Name,attr"`
	Code      string   `xml:"Code,attr"`
	Type      string   `xml:"Type,attr"`
	Size      string   `xml:"Size,attr"`
	Required  string   `xml:"Required,attr"`
	JsonField string
}

func Xml2Code() {
	// xml 解析为实体
	file, err := os.ReadFile("./out/out.xml")
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	v := MDB{}
	err = xml.Unmarshal(file, &v)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}

	model, err := os.ReadFile("./tmpl/model.tmpl")
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	files, err := template.New("model").Parse(string(model))
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	for i := 0; i < len(v.Modules); i++ {
		module := v.Modules[i]
		for j := 0; j < len(module.Tables); j++ {
			err = files.Execute(os.Stdout, module.Tables[j])
			if err != nil {
				fmt.Printf("error:%v\n", err)
			}
		}
	}

}
