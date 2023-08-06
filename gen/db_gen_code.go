package gen

import (
	"codeGenerator-Go/gen/model"
	"codeGenerator-Go/global"
	"codeGenerator-Go/utils"
	"fmt"
	"strings"
)

type TableInfo struct {
	TableName    string `gorm:"column:TABLE_NAME"`
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

type ColumnInfo struct {
	ColumnName             string `gorm:"column:COLUMN_NAME"`
	ColumnType             string `gorm:"column:COLUMN_TYPE"`
	CharacterMaximumLength string `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
	ColumnComment          string `gorm:"column:COLUMN_COMMENT"`
}

func DbGenCode() {
	set := make(map[string]bool)
	// xml 解析为实体
	var tables []TableInfo
	err := global.GOrmDao.Raw("SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?", global.GVA_VP.GetString("database.db_name")).Scan(&tables).Error
	if err != nil {
		// 处理错误
		fmt.Printf("error:%v\n", err)
	}
	for i := 0; i < len(tables); i++ {
		temp := tables[i].TableName
		if strings.Index(temp, "_") > 0 {
			set[temp[0:strings.Index(temp, "_")]] = true
		}
	}
	var models []model.Module
	for key := range set {
		modelObj := model.Module{}
		modelObj.Code = key
		var modelTables []model.Table
		for i := 0; i < len(tables); i++ {
			temp := tables[i].TableName
			if strings.Index(temp, key) == 0 {
				var madelTable model.Table
				madelTable.Code = tables[i].TableName
				madelTable.Comment = tables[i].TableComment
				// 处理字段
				getColumnByTableName(madelTable.Name)
				modelTables = append(modelTables, madelTable)
			}
		}
		modelObj.Tables = modelTables
		models = append(models, modelObj)
	}

	//file, err := os.ReadFile(xmlFilePath)
	//if err != nil {
	//	fmt.Printf("error:%v\n", err)
	//}
	//v := &model.MDB{}
	//err = xml.Unmarshal(file, &v)
	//if err != nil {
	//	fmt.Printf("error:%v\n", err)
	//}
	//genLanguage(v)
	GenLanguage(models)
}

func getColumnByTableName(tableName string) (res []model.Column) {
	var columns []ColumnInfo
	err := global.GOrmDao.Raw("SELECT column_name,column_type,CHARACTER_MAXIMUM_LENGTH,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;", global.GVA_VP.GetString("database.db_name"), tableName).Scan(&columns).Error
	if err != nil {
		// 处理错误
		fmt.Printf("error:%v\n", err)
	}
	for _, dbColumn := range columns {
		var column model.Column
		column.FieldName = utils.ToTitle(dbColumn.ColumnName)
		column.JsonField = utils.ToCamelCase(dbColumn.ColumnName)
		column.FieldType = utils.ConvertDbTypeToGoType(dbColumn.ColumnType)
		if column.FieldType == "time.Time" {
			column.IsTime = true
		}
	}
	return
}
