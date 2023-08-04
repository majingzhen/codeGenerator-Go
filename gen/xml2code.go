package gen

import (
	"encoding/xml"
	"fmt"
	"go-gen2pdm/utils"
	"html/template"
	"os"
	"strings"
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
	DateTime    string
}

type Column struct {
	XMLName   xml.Name `xml:"Column"`
	Name      string   `xml:"Name,attr"`
	Code      string   `xml:"Code,attr"`
	Type      string   `xml:"Type,attr"`
	Size      string   `xml:"Size,attr"`
	Required  string   `xml:"Required,attr"`
	FieldName string
	FieldType string
	JsonField string
}

func Xml2Code(xmlFilePath string) {
	// xml 解析为实体
	file, err := os.ReadFile(xmlFilePath)
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
	template, err := template.New("model").Parse(string(model))
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}

	for i := 0; i < len(v.Modules); i++ {
		module := &v.Modules[i]
		handleTable(module.Tables, module, template)
	}

}

func handleTable(tables []Table, module *Module, template *template.Template) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	for j := 0; j < len(tables); j++ {
		// 字段处理
		table := &tables[j]
		table.StructName = utils.ToTitle(table.Code)
		table.FileName = strings.ToLower(table.Code)
		table.DateTime = formattedTime
		table.Author = table.DBCreator
		table.PackageName = "model"
		for k := range table.Columns {
			column := &table.Columns[k]
			column.FieldName = utils.ToTitle(column.Code)
			column.JsonField = utils.ToCamelCase(column.Code)
			column.FieldType = utils.ConvertDbTypeToGoType(column.Type)
		}
		filePath := utils.GetProjectPath() + "/gen/code/" + module.Code + "/" + table.PackageName + "/"
		// 判断文件夹是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// 文件夹不存在，创建文件夹
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s文件夹创建成功/n", filePath)
		}
		// 生成写入文件
		WriteFile(filePath, table, template)
	}
}

func WriteFile(filePath string, table *Table, template *template.Template) {
	create, err := os.Create(filePath + table.FileName + ".go")
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	defer func(create *os.File) {
		err := create.Close()
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
	}(create)
	err = template.Execute(create, table)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	fmt.Printf("Create File:%v\n", create.Name())
}
