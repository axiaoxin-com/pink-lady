package routes

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Flatpages(app *gin.Engine) {
	navPath := viper.GetString("flatpages.nav_path")
	if navPath == "" {
		navPath = "fp"
	}

	rootDir := viper.GetString("flatpages.file_path")

	fileNames, err := LoadFlatpageFileNames(rootDir)
	if err != nil {
		logging.Error(nil, "LoadFlatpageFiles error:"+err.Error())
		return
	}

	totalCount := len(fileNames)

	fileNameIndexMap := map[string]int{}
	for idx, fileName := range fileNames {
		fileNameIndexMap[fileName] = idx
	}

	fp := app.Group(navPath)
	fp.GET("/", func(c *gin.Context) {
		meta := NewMetaData(c, webserver.CtxI18n(c, viper.GetString("flatpages.nav_name")))
		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			logging.Warn(c, "parse offset error:"+err.Error())
		}
		if offset < 0 {
			offset = 0
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			logging.Warn(c, "parse limit error:"+err.Error())
		}
		if limit > 100 || limit < 1 {
			limit = 10
		}
		pagi := goutils.PaginateByOffsetLimit(totalCount, offset, limit)
		data := gin.H{
			"meta": meta,
			"pagi": pagi,
		}
		if totalCount > 0 {
			data["fileNames"] = fileNames[pagi.StartIndex:pagi.EndIndex]
		}
		c.HTML(http.StatusOK, "flatpages_index.html", data)
		return
	})

	fp.GET("/:fileName/", func(c *gin.Context) {
		fileName, err := url.PathUnescape(c.Param("fileName"))
		if err != nil {
			logging.Error(c, "unescape flatpage fileName error:"+err.Error())
			c.Redirect(302, GetHostURL(c))
			return
		}

		fileName = strings.TrimPrefix(filepath.Join(string(filepath.Separator), fileName), string(filepath.Separator))
		fullFileName := filepath.Join(rootDir, fileName) + ".md"

		file, err := os.ReadFile(fullFileName)
		if err != nil {
			logging.Error(c, "read flatpages file error:"+err.Error())
			c.Redirect(302, GetHostURL(c))
			return
		}
		content := webserver.CtxI18n(c, string(file))

		meta := NewMetaData(c, fileName)
		data := gin.H{
			"meta":    meta,
			"content": string(content),
		}

		index := fileNameIndexMap[fileName]
		if index > 0 {
			data["nextFileName"] = fileNames[index-1]
		}
		if index < totalCount-1 {
			data["prevFileName"] = fileNames[index+1]
		}

		c.HTML(http.StatusOK, "flatpages_single.html", data)
		return
	})
}

// LoadFlatpageFileNames 加载flatpage文件名称列表（仅加载.md文件）
// 按最后修改时间降序排序
func LoadFlatpageFileNames(rootDir string) ([]string, error) {
	fileNames := []string{}

	f, err := os.Open(rootDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.ReadDir(-1)
	sort.Slice(files, func(i, j int) bool {
		iInfo, err := files[i].Info()
		if err != nil {
			logging.Error(nil, "get file info error:"+err.Error())
			return false
		}
		jInfo, err := files[j].Info()
		if err != nil {
			logging.Error(nil, "get file info error:"+err.Error())
			return false
		}
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	for _, f := range files {
		fileName := f.Name()
		ext := path.Ext(fileName)
		if strings.ToLower(ext) != ".md" {
			continue
		}
		fileName = strings.TrimSuffix(fileName, ext)
		logging.Debug(nil, "LoadFlatpageFileNames:"+fileName)
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}
