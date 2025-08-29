package routes

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/statics"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	// DefaultWordsPerMinute represents average reading speed
	DefaultWordsPerMinute = 300
	// DefaultTitle is used when no title is found in markdown
	DefaultTitle = "Untitled Flatpage"
)

// FlatpagesConfig holds the root configuration for flatpages
type FlatpagesConfig struct {
	Enable bool             `mapstructure:"enable"`
	Dirs   []FlatpageConfig `mapstructure:"dirs"`
}

// FlatpageConfig holds the configuration for flatpages
type FlatpageConfig struct {
	NavName  string `mapstructure:"nav_name"`
	NavPath  string `mapstructure:"nav_path"`
	MetaDesc string `mapstructure:"meta_desc"`
	FilePath string `mapstructure:"file_path"`
	PageSize int    `mapstructure:"page_size"`
}

// Flatpage represents a markdown flatpage with its metadata
type Flatpage struct {
	Title       string
	Slug        string
	Description string
	Content     string
	UpdatedAt   string
	ReadTime    int
	NavPath     string // Added to track which nav path this page belongs to
}

// FlatpageGroup represents a group of flatpages with their configuration
type FlatpageGroup struct {
	Config FlatpageConfig
	Pages  []*Flatpage
	Total  int // 总页数
}

var (
	// allFlatpageGroups stores all loaded flatpage groups
	allFlatpageGroups = map[string]*FlatpageGroup{}
	// wordPattern is used to match Chinese characters and English words
	wordPattern = regexp.MustCompile(`[a-zA-Z]+|\p{Han}`)
	// defaultPageSize is the default number of items per page
	defaultPageSize = 10
	flatpagesConfig = &FlatpagesConfig{}
)

// InitFlatpages initializes flatpage routes and loads all markdown documents
func InitFlatpages(app *gin.Engine) error {
	if err := viper.UnmarshalKey("flatpages", flatpagesConfig); err != nil {
		return fmt.Errorf("failed to unmarshal flatpages config: %v", err)
	}

	if !flatpagesConfig.Enable {
		logging.Info(nil, "Flatpages is disabled")
		return nil
	}

	for _, cfg := range flatpagesConfig.Dirs {
		// 如果没有指定 NavPath，使用文件夹名称作为默认值
		if cfg.NavPath == "" {
			cfg.NavPath = filepath.Base(cfg.FilePath)
		}

		// 如果没有指定 NavName，使用 NavPath 作为默认值
		if cfg.NavName == "" {
			cfg.NavName = cfg.NavPath
		}

		// 确保 PageSize 有效
		if cfg.PageSize <= 0 {
			cfg.PageSize = defaultPageSize
		}

		pages, err := loadAllFlatpages(cfg.FilePath)
		if err != nil {
			logging.Warnf(nil, "Error loading flatpages from %s: %v", cfg.FilePath, err)
			continue
		}

		// Set NavPath for each page
		for _, page := range pages {
			page.NavPath = cfg.NavPath
		}

		allFlatpageGroups[cfg.NavPath] = &FlatpageGroup{
			Config: cfg,
			Pages:  pages,
			Total:  len(pages),
		}

		logging.Infof(nil, "Successfully loaded %d flatpages from %s with nav_path: %s, page_size: %d",
			len(pages), cfg.FilePath, cfg.NavPath, cfg.PageSize)

		// Register routes for this group
		fp := app.Group(cfg.NavPath)
		fp.GET("/", handleFlatpageList)
		fp.GET("/:slug", handleFlatpageDetail)
	}

	return nil
}

// handleFlatpageList handles the flatpage list page request
func handleFlatpageList(c *gin.Context) {
	navPath := strings.TrimLeft(c.Request.URL.Path, "/")
	navPath = strings.TrimRight(navPath, "/")

	group, exists := allFlatpageGroups[navPath]
	if !exists {
		c.String(http.StatusNotFound, "Flatpage group not found")
		return
	}

	meta := NewMetaData(c, webserver.CtxI18n(c, group.Config.NavName))
	if group.Config.MetaDesc != "" {
		meta.BaseDesc = webserver.CtxI18n(c, group.Config.MetaDesc)
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		logging.Warn(c, "parse offset error:"+err.Error())
	}
	pagi := goutils.PaginateByOffsetLimit(group.Total, offset, group.Config.PageSize)

	pages := group.Pages[pagi.StartIndex:pagi.EndIndex]

	data := gin.H{
		"meta":         meta,
		"allFlatpages": pages,
		"pagi":         pagi,
		"navName":      group.Config.NavName,
	}

	c.HTML(http.StatusOK, "flatpages.html", data)
}

// handleFlatpageDetail handles individual flatpage request
func handleFlatpageDetail(c *gin.Context) {
	pathParts := strings.Split(strings.TrimLeft(c.Request.URL.Path, "/"), "/")
	navPath := strings.Join(pathParts[:len(pathParts)-1], "/")
	group, exists := allFlatpageGroups[navPath]
	if !exists {
		c.String(http.StatusNotFound, "Flatpage group not found")
		return
	}

	slug := c.Param("slug")
	currentPage, prevPage, nextPage := findFlatpageBySlug(group.Pages, slug)

	if currentPage == nil {
		c.String(http.StatusNotFound, "Flatpage not found")
		return
	}

	navName := webserver.CtxI18n(c, group.Config.NavName)
	meta := NewMetaData(c, webserver.CtxI18n(c, currentPage.Title)+"-"+navName)
	meta.BaseDesc = webserver.CtxI18n(c, currentPage.Description)

	data := gin.H{
		"meta":     meta,
		"flatpage": currentPage,
		"prev":     prevPage,
		"next":     nextPage,
		"navName":  navName,
	}

	c.HTML(http.StatusOK, "flatpage.html", data)
}

// findFlatpageBySlug finds a flatpage and its adjacent pages by slug
func findFlatpageBySlug(pages []*Flatpage, slug string) (current, prev, next *Flatpage) {
	for i, page := range pages {
		if page.Slug == slug {
			current = pages[i]
			if i > 0 {
				prev = pages[i-1]
			}
			if i < len(pages)-1 {
				next = pages[i+1]
			}
			break
		}
	}
	return
}

// loadAllFlatpages reads and parses all markdown files from the specified directory
func loadAllFlatpages(flatpagesPath string) ([]*Flatpage, error) {
	var entries []fs.DirEntry
	var err error

	// 如果路径以 flatpages/ 开头，从 embed 读取
	if strings.HasPrefix(flatpagesPath, "flatpages/") {
		entries, err = statics.Files.ReadDir(flatpagesPath)
	} else {
		// 否则从本地文件系统读取
		entries, err = os.ReadDir(flatpagesPath)
	}
	if err != nil {
		return nil, err
	}

	var pages []*Flatpage
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			page, err := loadSingleFlatpage(flatpagesPath, entry)
			if err != nil {
				logging.Errorf(nil, "Error loading flatpage %s: %v", entry.Name(), err)
				continue
			}
			pages = append(pages, page)
		}
	}

	// Sort pages by updated time (newest first)
	sort.Slice(pages, func(i, j int) bool {
		timeI, _ := time.Parse(time.DateOnly, pages[i].UpdatedAt)
		timeJ, _ := time.Parse(time.DateOnly, pages[j].UpdatedAt)
		return timeI.After(timeJ)
	})

	return pages, nil
}

// loadSingleFlatpage loads and parses a single markdown file
func loadSingleFlatpage(basePath string, entry fs.DirEntry) (*Flatpage, error) {
	filePath := filepath.Join(basePath, entry.Name())

	var content []byte
	var err error

	// 如果路径以 flatpages/ 开头，从 embed 读取
	if strings.HasPrefix(basePath, "flatpages/") {
		content, err = statics.Files.ReadFile(filePath)
	} else {
		// 否则从本地文件系统读取
		content, err = os.ReadFile(filePath)
	}
	if err != nil {
		return nil, err
	}

	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	return parseMarkdownFlatpage(content, entry.Name(), info.ModTime()), nil
}

// parseMarkdownFlatpage extracts information from markdown content
func parseMarkdownFlatpage(content []byte, filename string, modTime time.Time) *Flatpage {
	lines := strings.Split(string(content), "\n")
	title := DefaultTitle
	description := ""

	// Extract title from first h1 heading
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			break
		}
	}

	// Extract description from the first blockquote
	for _, line := range lines {
		if strings.HasPrefix(line, "> ") {
			description = strings.TrimPrefix(line, "> ")
			break
		}
	}

	// Generate slug from filename
	slug := strings.TrimSuffix(filename, filepath.Ext(filename))

	return &Flatpage{
		Title:       title,
		Slug:        slug,
		Description: description,
		Content:     string(content),
		UpdatedAt:   modTime.Format(time.DateOnly),
		ReadTime:    calculateReadTime(string(content)),
	}
}

// calculateReadTime estimates reading time in minutes
func calculateReadTime(content string) int {
	wordCount := len(wordPattern.FindAllString(content, -1))
	readTime := wordCount / DefaultWordsPerMinute
	if readTime < 1 {
		return 1
	}
	return readTime
}
