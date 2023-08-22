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
		genCode(modules, "./tmpl/go/model.tmpl", "model", "", "", false, false)
		genCode(modules, "./tmpl/go/dao.tmpl", "model", "_dao", "", false, false)
		genCode(modules, "./tmpl/go/dao_init.tmpl", "model", "", "", true, false)

		// 生成service
		genCode(modules, "./tmpl/go/service.tmpl", "service", "_service", "", false, false)
		genCode(modules, "./tmpl/go/service_init.tmpl", "service", "", "", true, false)
		genCode(modules, "./tmpl/go/view.tmpl", "view", "_view", "service", false, false)
		genCode(modules, "./tmpl/go/view_utils.tmpl", "view", "_view_utils", "service", false, false)
		genCode(modules, "./tmpl/go/view_init.tmpl", "view", "", "service", true, false)
		// 生成api
		genCode(modules, "./tmpl/go/api.tmpl", "api", "_api", "", false, false)
		genCode(modules, "./tmpl/go/api_init.tmpl", "api", "", "", true, false)
		// 生成router
		genCode(modules, "./tmpl/go/router.tmpl", "router", "_router", "", false, false)
	case "vue":
		// 生成vue
		genCode(modules, "./tmpl/vue/index.tmpl", "", "", "", false, true)
		genCode(modules, "./tmpl/vue/add.tmpl", "", "_add", "", false, false)
		genCode(modules, "./tmpl/vue/update.tmpl", "", "_update", "", false, false)
		genCode(modules, "./tmpl/vue/detail.tmpl", "", "_detail", "", false, false)
		genCode(modules, "./tmpl/vue/js.tmpl", "", "", "js", true, false)
	}

}

func genCode(modules []model.Module, tmplFile string, packageName, suffix string, appendPage string, isInit bool, isIndex bool) {
	model, err := os.ReadFile(tmplFile)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	template, err := template.New(packageName).Parse(string(model))
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	for _, module := range modules {
		handleTable(&module, template, packageName, suffix, appendPage, isInit, isIndex)
	}
}

func handleTable(module *model.Module, template *template.Template, packageName, suffix string, appendPage string, isInit bool, isIndex bool) {
	tables := module.Tables
	for j := 0; j < len(tables); j++ {
		var templateModel model.TemplateModel
		// 字段处理
		table := &tables[j]
		templateModel.StructName = utils.ToTitle(table.Code)
		if isIndex {
			templateModel.FileName = "index"
		} else {
			templateModel.FileName = strings.ToLower(table.Code)
		}
		templateModel.LastPathName = strings.ToLower(table.Code)
		templateModel.ObjectName = utils.ToCamelCase(table.Code)
		templateModel.DateTime = utils.GetCurTimeStr()
		templateModel.Author = table.DBCreator
		templateModel.PackageName = packageName
		templateModel.ModuleName = module.Code
		templateModel.ImportPath = global.GVA_VP.GetString("project.import_path")
		templateModel.CodePath = global.GVA_VP.GetString("project.code_path")
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
			filePath += module.Code + "/" + templateModel.LastPathName + "/" + appendPage + "/" + templateModel.PackageName + "/"
		} else {
			filePath += module.Code + "/" + templateModel.LastPathName + "/" + templateModel.PackageName + "/"
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
				path = filePath + templateModel.FileName + suffix + "." + global.GVA_VP.GetString("gen_code.language")
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
