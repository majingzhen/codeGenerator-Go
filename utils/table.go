package utils

import "strings"

func ConvertDbTypeToGoType(dbType string) string {
	dbType = strings.ToLower(dbType)
	switch dbType {
	case "int":
		return "int"
	case "varchar":
		return "string"
	case "text":
		return "string"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	case "datetime":
		return "time.Time"
	case "date":
		return "time.Time"
	case "bigint":
		return "int"
	case "tinyint":
		return "int8"
	default:
		return "interface{}"
	}
}
