package mapper

type PageData struct {
	Current    int32       `json:"current"`
	PageSize   int32       `json:"pageSize"`
	TotalCount int64       `json:"totalCount"`
	TotalPage  int32       `json:"totalPage"`
	List       interface{} `json:"list"`
}

func builderPageData(thisPage int32, pageSize int32, totalSize int64, totalPage int32, list interface{}) PageData {
	return PageData{Current: thisPage, PageSize: pageSize, TotalCount: totalSize, TotalPage: totalPage, List: list}
}
