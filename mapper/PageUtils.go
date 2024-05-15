package mapper

type PageUtils struct {
	thisPage  int32        // 当前页 	必须 默认1
	pageSize  int32        // 每页条数 必须 默认10
	totalSize int64        // 总条数 非用户添加
	totalPage int32        // 总页数 非用户添加
	wrapper   QueryWrapper // 条件查询器
}

func BuilderPageUtils(thisPage int32, pageSize int32, wrapper QueryWrapper) PageUtils {
	return PageUtils{thisPage: thisPage, pageSize: pageSize, wrapper: wrapper}
}

func (p *PageUtils) getPageSize() int32 {
	if p.pageSize <= 0 {
		p.pageSize = 10
	}
	return p.pageSize
}
func (p *PageUtils) getThisPage() int32 {
	if p.thisPage <= 0 {
		p.thisPage = 1
	}
	return p.thisPage
}

func (p *PageUtils) getOffSet() int32 {
	pageSize := p.getPageSize()
	thisPage := p.getThisPage()
	offSet := (thisPage - 1) * pageSize
	return offSet
}

func (p *PageUtils) setTotalSize(totalSize int64) *PageUtils {
	if totalSize <= 0 {
		totalSize = 0
	}
	p.totalSize = totalSize
	return p
}

func (p *PageUtils) getTotalPage() int32 {
	totalSize := p.totalSize
	pageSize := p.getPageSize()
	totalPage := totalSize / int64(pageSize)
	if totalSize%int64(pageSize) != 0 {
		totalPage = totalPage + 1
	}
	p.totalPage = int32(totalPage)
	return int32(totalPage)
}

func (p *PageUtils) pageData() PageData {
	return builderPageData(p.getThisPage(), p.getPageSize(), p.totalSize, p.getTotalPage(), p.wrapper.resList)
}
