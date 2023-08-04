package common

type PageInfoV2 struct {
	Limit     int         `json:"limit" form:"limit"`         // 每页大小
	Offset    int         `json:"offset" form:"offset"`       // 当前页码
	Page      int         `json:"page" form:"page"`           // 页码
	Total     int64       `json:"total" form:"total"`         // 总页
	PageSize  int         `json:"pageSize" form:"pageSize"`   // 每页大小
	Keyword   string      `json:"keyword" form:"keyword"`     // 查询关键字
	QueryInfo interface{} `json:"queryInfo" form:"queryInfo"` // 查询信息，例如date1>0这样的数据
	FormList  interface{} `json:"formList" form:"formList"`   // 返回列表
}

type Id struct {
	ID string `json:"id" form:"id"` // 主键ID
}

type Ids struct {
	Ids []string `json:"ids" form:"ids"`
}
