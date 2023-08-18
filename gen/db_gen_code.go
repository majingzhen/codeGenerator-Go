package gen

import (
	"codeGenerator-Go/gen/model"
	"codeGenerator-Go/global"
	"fmt"
	"strings"
)

type TableInfo struct {
	TableName    string `gorm:"column:TABLE_NAME"`
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

type ColumnInfo struct {
	ColumnName             string `gorm:"column:COLUMN_NAME"`
	DataType               string `gorm:"column:DATA_TYPE"`
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
				columns := getColumnByTableName(madelTable.Code)
				madelTable.Columns = columns
				modelTables = append(modelTables, madelTable)
			}
		}

		modelObj.Tables = modelTables
		models = append(models, modelObj)
	}

	GenLanguage(models)
}

func getColumnByTableName(tableName string) (res []model.Column) {
	var columns []ColumnInfo
	err := global.GOrmDao.Raw("SELECT COLUMN_NAME,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;", global.GVA_VP.GetString("database.db_name"), tableName).Scan(&columns).Error
	if err != nil {
		// 处理错误
		fmt.Printf("error:%v\n", err)
	}
	for _, dbColumn := range columns {
		var column model.Column
		column.Code = dbColumn.ColumnName
		column.Type = dbColumn.DataType
		column.Comment = dbColumn.ColumnComment
		res = append(res, column)
	}
	return
}
