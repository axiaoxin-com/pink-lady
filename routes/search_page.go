package routes

import (
	"net/http"
	"strings"

	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

// SearchResult 搜索结果
type SearchResult struct {
	Title       string
	Description string
	URL         string
	Type        string
}

// Search 处理搜索请求
func Search(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	meta := NewMetaData(c, webserver.CtxI18n(c, "搜索结果")+" - "+keyword)

	// 搜索结果列表
	var results []SearchResult

	// 搜索flatpages中的内容
	if flatpagesConfig != nil && flatpagesConfig.Enable {
		for _, group := range flatpagesConfig.Dirs {
			for _, page := range allFlatpageGroups[group.NavPath].Pages {
				if strings.Contains(strings.ToLower(page.Title), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(page.Description), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(page.Content), strings.ToLower(keyword)) {
					results = append(results, SearchResult{
						Title:       page.Title,
						Description: page.Description,
						URL:         "/" + group.NavPath + "/" + page.Slug,
						Type:        webserver.CtxI18n(c, "文档"),
					})
				}
			}
		}
	}

	data := gin.H{
		"meta":       meta,
		"keyword":    keyword,
		"results":    results,
		"hasResults": len(results) > 0,
	}

	c.HTML(http.StatusOK, "search.html", data)
}
