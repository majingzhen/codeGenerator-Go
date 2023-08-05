package gen

import (
	"codeGenerator-Go/global"
	"codeGenerator-Go/utils"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"text/template"
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
	XMLName   xml.Name `xml:"Table"`
	Name      string   `xml:"Name,attr"`
	Code      string   `xml:"Code,attr"`
	HtmlPkg   string   `xml:"HtmlPkg,attr"`
	GoPkg     string   `xml:"GoPkg,attr"`
	DBCreator string   `xml:"DBCreator,attr"`
	Comment   string   `xml:"Comment,attr"`
	Columns   []Column `xml:"Column"`
}

type TemplateModel struct {
	Table       *Table
	Author      string
	FileName    string
	PackageName string
	StructName  string
	ObjectName  string
	DateTime    string
	ProjectName string
	ModuleName  string
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
	IsTime    bool
}

func Xml2Code(xmlFilePath string) {
	// xml 解析为实体
	file, err := os.ReadFile(xmlFilePath)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	v := &MDB{}
	err = xml.Unmarshal(file, &v)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	genLanguage(v)
}

func genLanguage(v *MDB) {
	switch global.GVA_VP.GetString("gen_code.language") {
	case "go":
		// 生成model
		genCode(v, "./tmpl/go/model.tmpl", "model", "", "", false)
		genCode(v, "./tmpl/go/dao.tmpl", "model", "dao", "", false)
		genCode(v, "./tmpl/go/dao_init.tmpl", "model", "", "", true)

		// 生成service
		genCode(v, "./tmpl/go/service.tmpl", "service", "service", "", false)
		genCode(v, "./tmpl/go/service_init.tmpl", "service", "", "", true)
		genCode(v, "./tmpl/go/view.tmpl", "view", "view", "service", false)
		genCode(v, "./tmpl/go/view_page.tmpl", "view", "view_page", "service", false)
		genCode(v, "./tmpl/go/view_utils.tmpl", "view", "view_utils", "service", false)
		genCode(v, "./tmpl/go/view_init.tmpl", "view", "", "service", true)
		// 生成api
		genCode(v, "./tmpl/go/api.tmpl", "api", "api", "", false)
		genCode(v, "./tmpl/go/api_init.tmpl", "api", "", "", true)
		// 生成router
		genCode(v, "./tmpl/go/router.tmpl", "router", "router", "", false)
	case "vue":
		// 生成vue
		genCode(v, "./tmpl/vue/index.tmpl", "", "", "", false)
		genCode(v, "./tmpl/vue/add.tmpl", "", "add", "", false)
		genCode(v, "./tmpl/vue/update.tmpl", "", "update", "", false)
		genCode(v, "./tmpl/vue/detail.tmpl", "", "detail", "", false)
		genCode(v, "./tmpl/vue/js.tmpl", "", "", "js", true)
	}

}

func genCode(v *MDB, tmplFile string, packageName, suffix string, appendPage string, isInit bool) {
	model, err := os.ReadFile(tmplFile)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	template, err := template.New(packageName).Parse(string(model))
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	for _, module := range v.Modules {
		handleTable(&module, template, packageName, suffix, appendPage, isInit)
	}
}

func handleTable(module *Module, template *template.Template, packageName, suffix string, appendPage string, isInit bool) {
	tables := module.Tables
	for j := 0; j < len(tables); j++ {
		var templateModel TemplateModel
		// 字段处理
		table := &tables[j]
		templateModel.StructName = utils.ToTitle(table.Code)
		templateModel.FileName = strings.ToLower(table.Code)
		templateModel.ObjectName = utils.ToCamelCase(table.Code)
		templateModel.DateTime = utils.GetCurTimeStr()
		templateModel.Author = table.DBCreator
		templateModel.PackageName = packageName
		templateModel.ModuleName = module.Code
		templateModel.ProjectName = global.GVA_VP.GetString("project_name")
		for k := range table.Columns {
			column := &table.Columns[k]
			column.FieldName = utils.ToTitle(column.Code)
			column.JsonField = utils.ToCamelCase(column.Code)
			column.FieldType = utils.ConvertDbTypeToGoType(column.Type)
			if column.FieldType == "time.Time" {
				column.IsTime = true
			}
		}
		var filePath string

		filePath = global.GVA_VP.GetString("gen_code.out_path") + global.GVA_VP.GetString("gen_code.language") + "/"

		if appendPage != "" {
			filePath += module.Code + "/" + templateModel.FileName + "/" + appendPage + "/" + templateModel.PackageName + "/"
		} else {
			filePath += module.Code + "/" + templateModel.FileName + "/" + templateModel.PackageName + "/"
		}
		// 判断文件夹是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// 文件夹不存在，创建文件夹
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				panic(err)
			}
			fmt.Printf("文件夹创建成功: %s\n", filePath)
		}
		templateModel.Table = table
		var path string
		if global.GVA_VP.GetString("gen_code.language") == "go" && isInit {
			path = filePath + "init.go"
		} else if global.GVA_VP.GetString("gen_code.language") == "vue" && isInit {
			path = filePath + templateModel.FileName + ".js"
		} else {
			if suffix != "" {
				path = filePath + templateModel.FileName + "_" + suffix + "." + global.GVA_VP.GetString("gen_code.language")
			} else {
				path = filePath + templateModel.FileName + "." + global.GVA_VP.GetString("gen_code.language")
			}
		}
		// 生成写入文件
		WriteFile(path, templateModel, template)
	}
}

func WriteFile(filePath string, templateModel TemplateModel, template *template.Template) {
	create, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	defer func(create *os.File) {
		err := create.Close()
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
	}(create)

	err = template.Execute(create, templateModel)

	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	fmt.Printf("Create File:%v\n", create.Name())
}
