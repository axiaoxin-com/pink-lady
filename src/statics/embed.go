// Package statics 静态文件嵌入
package statics

import "embed"

var (
	//go:embed files/*
	Files embed.FS
)
