package model

import "encoding/xml"

type Module struct {
	XMLName xml.Name `xml:"Module"`
	Name    string   `xml:"Name,attr"`
	Code    string   `xml:"Code,attr"`
	Tables  []Table  `xml:"Table"`
}
