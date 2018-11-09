package utils

import (
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitViper(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	confile, err := os.Open("../app.yaml")
	if err != nil {
		t.Error(err)
	}
	defer confile.Close()
	tmpConfile, err := os.Create("./app.yaml")
	if err != nil {
		t.Error(err)
	}
	defer tmpConfile.Close()
	io.Copy(tmpConfile, confile)
	// 清理测试用的配置文件
	defer func() { os.Remove("./app.yaml") }()

	// 测试初始化用法
	options := []ViperOption{
		ViperOption{Name: "option1", Default: 1, Desc: "number 1"},
		ViperOption{Name: "option2", Default: true, Desc: "bool true"},
		ViperOption{Name: "option3", Default: "3", Desc: "string 3"},
		ViperOption{Name: "option.4", Default: "o4", Desc: "."},
	}
	InitViper("app", "envPrefix", options)
	if viper.GetInt("option1") != 1 {
		t.Error("get int option error")
	}
	if viper.GetBool("option2") != true {
		t.Error("get bool option error")
	}
	if viper.GetString("option3") != "3" {
		t.Error("get string option error")
	}

	// 测试设置前缀环境变量优先生效
	os.Setenv("ENVPREFIX_OPTION1", "xxx")
	if viper.GetString("option1") != "xxx" {
		t.Error("bind env prefix failed")
	}
	// 测试环境变量.替换为_
	os.Setenv("ENVPREFIX_OPTION_4", ".o_")
	if viper.GetString("option.4") != ".o_" {
		t.Error("env replacer error")
	}

	// 测试加载配置文件
	if viper.GetString("server.bind") != ":9090" {
		t.Error("read conf file error")
	}
}
