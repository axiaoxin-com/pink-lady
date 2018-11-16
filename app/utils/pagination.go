package utils

import "math"

// Pagination Paginate return it
type Pagination struct {
	// 数据总数
	ItemsCount int `json:"itemsCount"`
	// 分页总数
	PagesCount int `json:"pagesCount"`
	// 当前页码
	PageNum int `json:"pageNum"`
	// 分页大小
	PageSize int `json:"pageSize"`
	// 是否有上一页
	HasPrev bool `json:"hasPrev"`
	// 是否有下一页
	HasNext bool `json:"hasNext"`
	// 上一页页码
	PrevPageNum int `json:"prevPageNum"`
	// 下一页页码
	NextPageNum int `json:"nextPageNum"`
}

// Paginate 计算分页信息
func Paginate(itemsCount, pageNum, pageSize int) Pagination {
	if itemsCount <= 0 {
		return Pagination{}
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	pagesCount := 1
	if pageSize <= 0 {
		pageSize = -1
	} else {
		pagesCount = int(math.Ceil(float64(itemsCount) / float64(pageSize)))
	}
	hasNext := false
	nextPageNum := pageNum
	if pageNum+1 <= pagesCount {
		hasNext = true
		nextPageNum = pageNum + 1
	}
	hasPrev := false
	prevPageNum := pageNum
	if pageNum-1 > 0 {
		hasPrev = true
		prevPageNum = pageNum - 1
	}
	return Pagination{
		ItemsCount:  itemsCount,
		PagesCount:  pagesCount,
		PageNum:     pageNum,
		PageSize:    pageSize,
		HasPrev:     hasPrev,
		HasNext:     hasNext,
		PrevPageNum: prevPageNum,
		NextPageNum: nextPageNum,
	}
}
