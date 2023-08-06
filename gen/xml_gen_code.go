package gen

import (
	"codeGenerator-Go/gen/model"
	"codeGenerator-Go/global"
	"codeGenerator-Go/utils"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func Xml2Code(xmlFilePath string) {
	// xml 解析为实体
	file, err := os.ReadFile(xmlFilePath)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	v := &model.MDB{}
	err = xml.Unmarshal(file, &v)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	GenLanguage(v.Modules)
}

func GenLanguage(modules []model.Module) {
	switch global.GVA_VP.GetString("gen_code.language") {
	case "go":
		// 生成model
		genCode(modules, "./tmpl/go/model.tmpl", "model", "", "", false)
		genCode(modules, "./tmpl/go/dao.tmpl", "model", "dao", "", false)
		genCode(modules, "./tmpl/go/dao_init.tmpl", "model", "", "", true)

		// 生成service
		genCode(modules, "./tmpl/go/service.tmpl", "service", "service", "", false)
		genCode(modules, "./tmpl/go/service_init.tmpl", "service", "", "", true)
		genCode(modules, "./tmpl/go/view.tmpl", "view", "view", "service", false)
		genCode(modules, "./tmpl/go/view_page.tmpl", "view", "view_page", "service", false)
		genCode(modules, "./tmpl/go/view_utils.tmpl", "view", "view_utils", "service", false)
		genCode(modules, "./tmpl/go/view_init.tmpl", "view", "", "service", true)
		// 生成api
		genCode(modules, "./tmpl/go/api.tmpl", "api", "api", "", false)
		genCode(modules, "./tmpl/go/api_init.tmpl", "api", "", "", true)
		// 生成router
		genCode(modules, "./tmpl/go/router.tmpl", "router", "router", "", false)
	case "vue":
		// 生成vue
		genCode(modules, "./tmpl/vue/index.tmpl", "", "", "", false)
		genCode(modules, "./tmpl/vue/add.tmpl", "", "add", "", false)
		genCode(modules, "./tmpl/vue/update.tmpl", "", "update", "", false)
		genCode(modules, "./tmpl/vue/detail.tmpl", "", "detail", "", false)
		genCode(modules, "./tmpl/vue/js.tmpl", "", "", "js", true)
	}

}

func genCode(modules []model.Module, tmplFile string, packageName, suffix string, appendPage string, isInit bool) {
	model, err := os.ReadFile(tmplFile)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	template, err := template.New(packageName).Parse(string(model))
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	for _, module := range modules {
		handleTable(&module, template, packageName, suffix, appendPage, isInit)
	}
}

func handleTable(module *model.Module, template *template.Template, packageName, suffix string, appendPage string, isInit bool) {
	tables := module.Tables
	for j := 0; j < len(tables); j++ {
		var templateModel model.TemplateModel
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

func WriteFile(filePath string, templateModel model.TemplateModel, template *template.Template) {
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
