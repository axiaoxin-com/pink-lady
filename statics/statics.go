// Package statics embed 静态文件
package statics

import "embed"

// Files 静态文件资源
//
//go:embed robots.txt
//go:embed css/* font/* html/* img/* js/*
//go:embed flatpages/*
var Files embed.FS
