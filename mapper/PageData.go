package mapper

type PageData[T any] struct {
	Current    int32 `json:"current"`
	PageSize   int32 `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPage  int32 `json:"totalPage"`
	List       []T   `json:"list"`
}

func builderPageData[T any](
	thisPage int32,
	pageSize int32,
	totalSize int64,
	totalPage int32,
	list []T,
) PageData[T] {
	return PageData[T]{
		Current:    thisPage,
		PageSize:   pageSize,
		TotalCount: totalSize,
		TotalPage:  totalPage,
		List:       list,
	}
}
