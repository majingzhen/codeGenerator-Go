package main

import "codeGenerator-Go/gen"

func main() {
	gen.Pdm2xml("./pdm/测试数据库模型 - 副本.PDM", true)

	gen.Xml2Code("./out/out.xml")
}
