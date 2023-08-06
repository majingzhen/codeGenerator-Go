package model

import "encoding/xml"

type Column struct {
	XMLName   xml.Name `xml:"Column"`
	Name      string   `xml:"Name,attr"`
	Code      string   `xml:"Code,attr"`
	Type      string   `xml:"Type,attr"`
	Size      string   `xml:"Size,attr"`
	Required  string   `xml:"Required,attr"`
	Comment   string   `xml:"Comment,attr"`
	FieldName string
	FieldType string
	JsonField string
	IsTime    bool
}
