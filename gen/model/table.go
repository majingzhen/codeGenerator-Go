package model

import "encoding/xml"

type Table struct {
	XMLName   xml.Name `xml:"Table"`
	Name      string   `gorm:"table_name" xml:"Name,attr"`
	Code      string   `xml:"Code,attr"`
	HtmlPkg   string   `xml:"HtmlPkg,attr"`
	GoPkg     string   `xml:"GoPkg,attr"`
	DBCreator string   `xml:"DBCreator,attr"`
	Comment   string   `gorm:"table_comment" xml:"Comment,attr"`
	Columns   []Column `xml:"Column"`
}
