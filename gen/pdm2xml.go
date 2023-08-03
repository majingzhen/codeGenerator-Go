package gen

import (
	"fmt"
	"github.com/beevik/etree"
	"os"
	"strconv"
	"strings"
)

var fileName = "./out/out.xml"

func addAttr(sb *strings.Builder, attr, val string) {
	sb.WriteString(" " + attr + "=\"")
	sb.WriteString(val + "\"")
}

// Pdm2xml pdm转xml
// pdmPath pdm文件路径
// isPack 数据库是否分包
func Pdm2xml(pdmPath string, isPack bool) {
	var newXML strings.Builder
	doc := etree.NewDocument()
	file, err := os.ReadFile(pdmPath)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	xml := string(file)

	for i := 0; i < 10; i++ {
		bgn := strings.Index(xml, "<a:PackageOptionsText>")
		end := strings.Index(xml, "</a:PackageOptionsText>")
		if bgn != -1 {
			xml = xml[:bgn] + xml[end+len("</a:PackageOptionsText>"):]
		}
	}
	for i := 0; i < 10; i++ {
		bgn := strings.Index(xml, "<a:ModelOptionsText>")
		end := strings.Index(xml, "</a:ModelOptionsText>")
		if bgn != -1 {
			xml = xml[:bgn] + xml[end+len("</a:ModelOptionsText>"):]
		}
	}
	for i := 0; i < 10; i++ {
		bgn := strings.Index(xml, "<a:DisplayPreferences>")
		end := strings.Index(xml, "</a:DisplayPreferences>")
		if bgn != -1 {
			xml = xml[:bgn] + xml[end+len("</a:DisplayPreferences>"):]
		}
	}
	// fmt.Println(xml)
	if err := doc.ReadFromString(xml); err != nil {
		fmt.Printf("error:%v\n", err)
	}
	var nodeList []*etree.Element
	if isPack {
		nodeList = doc.FindElements("/Model/RootObject/Children/Model/Packages/Package")
	} else {
		nodeList = doc.FindElements("/Model/RootObject/Children/Model")
	}

	s := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	newXML.WriteString(s)
	newXML.WriteString("\n")
	newXML.WriteString("<MDB>")
	newXML.WriteString("\n")
	for i := 0; i < len(nodeList); i++ {
		packageModel := nodeList[i]
		var packageCode string
		packageNodes := packageModel.ChildElements()
		for j := 0; j < len(packageNodes); j++ {
			packageNode := packageNodes[j]
			if packageNode.Tag == "Name" {
				newXML.WriteString("	 <Module>")
				addAttr(&newXML, "Name", packageNode.Text())
			}
			if packageNode.Tag == "Code" {
				packageCode = strings.ToLower(packageNode.Text())
				addAttr(&newXML, "Code", packageCode)
				newXML.WriteString(">")
				newXML.WriteString("\n")
			}
			if packageNode.Tag == "Tables" {
				tableNodes := packageNode.ChildElements()
				for t := 0; t < len(tableNodes); t++ {
					var tableCode string
					tableNode := tableNodes[t]
					tableChildNodes := tableNode.ChildElements()
					for c := 0; c < len(tableChildNodes); c++ {
						tableChildNode := tableChildNodes[c]
						if tableChildNode.Tag == "Name" {
							newXML.WriteString("    <Table")
							addAttr(&newXML, "Name", tableChildNode.Text())
						}
						if tableChildNode.Tag == "Code" {
							addAttr(&newXML, "Code", tableChildNode.Text())
							tableCode = strings.ReplaceAll(tableChildNode.Text(), "_", "")
						}
						if tableChildNode.Tag == "Comment" {
							addAttr(&newXML, "Comment", tableChildNode.Text())
						}
						if tableChildNode.Tag == "Creator" {
							addAttr(&newXML, "HtmlPkg", strings.ToLower(packageCode+"/"+tableCode))
							addAttr(&newXML, "GoPkg", strings.ToLower(packageCode+"."+tableCode))
							addAttr(&newXML, "DBCreator", tableChildNode.Text())
						}
						if tableChildNode.Tag == "TotalSavingCurrency" {
							newXML.WriteString(">")
							newXML.WriteString("\n")
						}

						// 处理 columns
						if tableChildNode.Tag == "Columns" {
							columnNodes := tableChildNode.ChildElements()
							for cc := 0; cc < len(columnNodes); cc++ {
								columnNode := columnNodes[cc]
								columnNodeAttrs := columnNode.ChildElements()
								isPk := false
								isRequire := false
								isQuery := false
								isSearch := false
								isLogic := false
								var size int
								var columnType string
								var colCode string
								for k := 0; k < len(columnNodeAttrs); k++ {
									columnNodeAttr := columnNodeAttrs[k]
									if columnNodeAttr.Tag == "ObjectID" {
										newXML.WriteString("      <Column")
									}
									if columnNodeAttr.Tag == "ID" {
										isPk = true
									}
									if !isPk && columnNodeAttr.Tag == "Mandatory" {
										isRequire = true
									}
									if columnNodeAttr.Tag == "NAME" {
										isSearch = true
										isQuery = true
									}
									if columnNodeAttr.Tag == "VALID_FLAG" {
										isLogic = true
									}
									if columnNodeAttr.Tag == "Code" {
										colCode = strings.ToUpper(columnNodeAttr.Text())
										addAttr(&newXML, "Code", colCode)
									}
									if columnNodeAttr.Tag == "Name" {
										addAttr(&newXML, "Name", columnNodeAttr.Text())
									}
									if columnNodeAttr.Tag == "Length" {
										size, _ = strconv.Atoi(columnNodeAttr.Text())
									}
									if columnNodeAttr.Tag == "Precision" {
										addAttr(&newXML, "Precision", columnNodeAttr.Text())
									}
									if columnNodeAttr.Tag == "Comment" {
										addAttr(&newXML, "Comment", columnNodeAttr.Text())
									}
									if columnNodeAttr.Tag == "DataType" {
										columnType = strings.ToUpper(columnNodeAttr.Text())
										if strings.Contains(columnType, "(") {
											columnType = columnType[0:strings.Index(columnType, "(")]
										}
										addAttr(&newXML, "Type", strings.ToUpper(columnType))
									}
								}
								// 处理常用字段长度
								if strings.Contains(columnNode.Tag, "Column") {
									if strings.Contains(columnType, "TEXT") {
										if strings.Contains(colCode, "NAME") {
											size = 32
										}
										if strings.Contains(colCode, "DESC") {
											size = 2048
										}
										if strings.Contains(colCode, "REMARK") {
											size = 256
										}
									}
									if strings.Contains(columnType, "INT") {
										if strings.Contains(colCode, "STATE") {
											size = 8
										}
										if strings.Contains(colCode, "TYPE") {
											size = 8
										}
										if strings.Contains(colCode, "DELETED") {
											size = 1 // 逻辑删除
										}
									}
									if len(colCode) > 0 && size > 0 {
										addAttr(&newXML, "Size", strconv.Itoa(size))
									}
									if isPk {
										newXML.WriteString(" Pk=\"" + strconv.FormatBool(isPk) + "\"")
									}
									if isRequire {
										newXML.WriteString(" Required=\"" + strconv.FormatBool(isRequire) + "\"")
									}
									if isQuery {
										newXML.WriteString(" NeedQuery=\"" + strconv.FormatBool(isQuery) + "\"")
									}
									if isSearch {
										newXML.WriteString(" NeedSearch=\"" + strconv.FormatBool(isSearch) + "\"")
									}
									if isLogic {
										newXML.WriteString(" LogicDelete=\"" + strconv.FormatBool(isLogic) + "\"")
									}
									newXML.WriteString("/>")
									newXML.WriteString("\n")
								}
							}
							newXML.WriteString("    </Table>")
							newXML.WriteString("\n")
						}
					}
				}
			}
		}
		newXML.WriteString("  </Module>")
		newXML.WriteString("\n")
	}
	newXML.WriteString("</MDB>")
	create, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}
	defer func(create *os.File) {
		err := create.Close()
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
	}(create)
	_, err = create.WriteString(newXML.String())
	if err != nil {
		fmt.Printf("error:%v\n", err)
	}

	fmt.Println("转换完成!")
}
