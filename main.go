package main

import (
	"codeGenerator-Go/before"
	"codeGenerator-Go/gen"
	"codeGenerator-Go/global"
)

func init() {
	before.Viper()
}

func main() {
	if global.GVA_VP.GetString("gen_code.data_source") == "xml" {
		if global.GVA_VP.GetBool("gen_code.pdm_2_xml.enable") {
			gen.Pdm2xml(global.GVA_VP.GetString("gen_code.pdm_2_xml.Source_file_path"), global.GVA_VP.GetBool("gen_code.pdm_2_xml.is_model"))
			gen.Xml2Code(global.GVA_VP.GetString("gen_code.pdm_2_xml.out_path"))
		} else {
			gen.Xml2Code(global.GVA_VP.GetString("gen_code.Source_file_path"))
		}
	}
}
