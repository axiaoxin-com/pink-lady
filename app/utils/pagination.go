package utils

import "math"

// Pagination Paginate return it
// 异常数据时分页总数为0，当前页码、上下页码均不判断逻辑，只管数值增减
type Pagination struct {
	// 数据总数
	TotalCount int `json:"TotalCount"`
	// 分页总数
	PagesCount int `json:"PagesCount"`
	// 当前页码
	PageNum int `json:"PageNum"`
	// 分页大小
	PageSize int `json:"PageSize"`
	// 是否有上一页
	HasPrev bool `json:"HasPrev"`
	// 是否有下一页
	HasNext bool `json:"HasNext"`
	// 上一页页码
	PrevPageNum int `json:"PrevPageNum"`
	// 下一页页码
	NextPageNum int `json:"NextPageNum"`
}

// PaginateByPageNumSize 按pagenum,pagesize计算分页信息
func PaginateByPageNumSize(totalCount, pageNum, pageSize int) *Pagination {
	if totalCount <= 0 {
		return nil
	}
	pagesCount := 0
	if pageSize > 0 {
		pagesCount = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	}
	hasNext := true
	nextPageNum := pageNum + 1
	if nextPageNum >= pagesCount {
		hasNext = false
	}
	hasPrev := true
	prevPageNum := pageNum - 1
	if prevPageNum <= 0 {
		hasPrev = false
	}
	return &Pagination{
		TotalCount:  totalCount,
		PagesCount:  pagesCount,
		PageNum:     pageNum,
		PageSize:    pageSize,
		HasPrev:     hasPrev,
		HasNext:     hasNext,
		PrevPageNum: prevPageNum,
		NextPageNum: nextPageNum,
	}
}

// PaginateByOffsetLimit 按offset,limit计算分页信息
func PaginateByOffsetLimit(totalCount, offset, limit int) *Pagination {
	if totalCount <= 0 {
		return nil
	}
	pageNum := 1
	if offset <= 0 {
		pageNum = 1
	} else {
		pageNum = offset/limit + 1
	}
	pagesCount := 0
	if limit > 0 {
		pagesCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}
	hasNext := true
	nextPageNum := pageNum + 1
	if limit == 0 || offset+limit >= totalCount || nextPageNum >= pagesCount {
		hasNext = false
	}
	hasPrev := true
	prevPageNum := pageNum - 1
	if limit == 0 || offset+limit <= 0 || pageNum == 1 {
		hasPrev = false
	}
	return &Pagination{
		TotalCount:  totalCount,
		PagesCount:  pagesCount,
		PageNum:     pageNum,
		PageSize:    limit,
		HasPrev:     hasPrev,
		HasNext:     hasNext,
		PrevPageNum: prevPageNum,
		NextPageNum: nextPageNum,
	}
}
