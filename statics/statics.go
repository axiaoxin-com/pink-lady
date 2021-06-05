// Package statics 静态文件嵌入
package statics

import "embed"

var (
	// Files 静态文件嵌入
	//go:embed files/*
	Files embed.FS
)
