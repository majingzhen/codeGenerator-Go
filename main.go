package main

import (
	"codeGenerator-Go/gen"
	"codeGenerator-Go/global"
)

func init() {
	global.InitConfig()
}

func main() {
	switch global.GVA_VP.GetString("gen_code.data_source") {
	case "xml":
		if global.GVA_VP.GetBool("gen_code.pdm_2_xml.enable") {
			gen.Pdm2xml(global.GVA_VP.GetString("gen_code.pdm_2_xml.Source_file_path"), global.GVA_VP.GetBool("gen_code.pdm_2_xml.is_model"))
			gen.Xml2Code(global.GVA_VP.GetString("gen_code.pdm_2_xml.out_path"))
		} else {
			gen.Xml2Code(global.GVA_VP.GetString("gen_code.Source_file_path"))
		}
	case "db":
		gen.DbGenCode()
	}

}
