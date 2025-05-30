package routes

import (
	"html/template"
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
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/viper"
)

const (
	// DefaultWordsPerMinute represents average reading speed
	DefaultWordsPerMinute = 300
	// DefaultNavPath is the default URL path for flatpages
	DefaultNavPath = "fp"
	// DefaultTitle is used when no title is found in markdown
	DefaultTitle = "Untitled Flatpage"
)

// FlatpageConfig holds the configuration for flatpages
type FlatpageConfig struct {
	FilePath string
	NavPath  string
}

// Flatpage represents a markdown flatpage with its metadata
type Flatpage struct {
	Title       string
	Slug        string
	Description string
	Content     template.HTML
	UpdatedAt   string
	ReadTime    int
}

var (
	// allFlatpages stores all loaded flatpage documents
	allFlatpages = []*Flatpage{}
	// wordPattern is used to match Chinese characters and English words
	wordPattern = regexp.MustCompile(`[a-zA-Z]+|\p{Han}`)
)

// InitFlatpages initializes flatpage routes and loads all markdown documents
func InitFlatpages(app *gin.Engine) error {
	cfg := FlatpageConfig{
		FilePath: viper.GetString("flatpages.file_path"),
		NavPath:  viper.GetString("flatpages.nav_path"),
	}
	if cfg.NavPath == "" {
		cfg.NavPath = DefaultNavPath
	}

	pages, err := loadAllFlatpages(cfg.FilePath)
	if err != nil {
		return err
	}
	allFlatpages = pages
	logging.Infof(nil, "Successfully loaded %d flatpages from %s", len(pages), cfg.FilePath)

	fp := app.Group(cfg.NavPath)
	fp.GET("/", handleFlatpageList)
	fp.GET("/:slug", handleFlatpageDetail)
	return nil
}

// handleFlatpageList handles the flatpage list page request
func handleFlatpageList(c *gin.Context) {
	meta := NewMetaData(c, webserver.CtxI18n(c, viper.GetString("flatpages.nav_name")))
	meta.BaseDesc = ""

	total := len(allFlatpages)
	limit := 10

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		logging.Warn(c, "parse offset error:"+err.Error())
	}
	pagi := goutils.PaginateByOffsetLimit(total, offset, limit)
	data := gin.H{
		"meta":         meta,
		"allFlatpages": allFlatpages[pagi.StartIndex:pagi.EndIndex],
		"pagi":         pagi,
	}

	c.HTML(http.StatusOK, "flatpages.html", data)
}

// handleFlatpageDetail handles individual flatpage request
func handleFlatpageDetail(c *gin.Context) {
	slug := c.Param("slug")
	currentPage, prevPage, nextPage := findFlatpageBySlug(slug)

	if currentPage == nil {
		c.String(http.StatusNotFound, "Flatpage not found")
		return
	}

	meta := NewMetaData(c, webserver.CtxI18n(c, currentPage.Title)+"-"+webserver.CtxI18n(c, viper.GetString("flatpages.nav_name")))
	meta.BaseDesc = currentPage.Description

	data := gin.H{
		"meta":     meta,
		"flatpage": currentPage,
		"prev":     prevPage,
		"next":     nextPage,
	}

	c.HTML(http.StatusOK, "flatpage.html", data)
}

// findFlatpageBySlug finds a flatpage and its adjacent pages by slug
func findFlatpageBySlug(slug string) (current, prev, next *Flatpage) {
	for i, page := range allFlatpages {
		if page.Slug == slug {
			current = allFlatpages[i]
			if i > 0 {
				prev = allFlatpages[i-1]
			}
			if i < len(allFlatpages)-1 {
				next = allFlatpages[i+1]
			}
			break
		}
	}
	return
}

// loadAllFlatpages reads and parses all markdown files from the specified directory
func loadAllFlatpages(flatpagesPath string) ([]*Flatpage, error) {
	entries, err := os.ReadDir(flatpagesPath)
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
func loadSingleFlatpage(basePath string, entry os.DirEntry) (*Flatpage, error) {
	filePath := filepath.Join(basePath, entry.Name())
	content, err := os.ReadFile(filePath)
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

	// Parse markdown content
	parsedContent := parseMarkdownToHTML(content)

	return &Flatpage{
		Title:       title,
		Slug:        slug,
		Description: description,
		Content:     parsedContent,
		UpdatedAt:   modTime.Format(time.DateOnly),
		ReadTime:    calculateReadTime(string(content)),
	}
}

// parseMarkdownToHTML converts markdown content to HTML
func parseMarkdownToHTML(content []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return template.HTML(markdown.Render(doc, renderer))
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
