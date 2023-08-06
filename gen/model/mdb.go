package model

import "encoding/xml"

type MDB struct {
	XMLName xml.Name `xml:"MDB"`
	Modules []Module `xml:"Module"`
	Tables  []Table  `xml:"Table"`
}
