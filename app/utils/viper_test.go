package utils

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitViper(t *testing.T) {
	// 创建config文件
	err := CopyFile("../config.toml.example", "/tmp/config-plut.toml")
	if err != nil {
		t.Fatal(err)
	}
	// 清理测试用的配置文件
	defer func() { os.Remove("/tmp/config-plut.toml") }()

	// 测试初始化用法
	options := []ViperOption{
		{Name: "option1", Default: 1, Desc: "number 1"},
		{Name: "option2", Default: true, Desc: "bool true"},
		{Name: "option3", Default: "3", Desc: "string 3"},
		{Name: "option.4", Default: "o4", Desc: "."},
	}
	if err := InitViper([]string{"/tmp/"}, "config-plut", "envPrefix", options...); err != nil {
		t.Error(err)
	}
	if viper.GetInt("option1") != 1 {
		t.Fatal("get int option error")
	}
	if viper.GetBool("option2") != true {
		t.Fatal("get bool option error")
	}
	if viper.GetString("option3") != "3" {
		t.Fatal("get string option error")
	}

	// 测试设置前缀环境变量优先生效
	os.Setenv("ENVPREFIX_OPTION1", "xxx")
	if viper.GetString("option1") != "xxx" {
		t.Fatal("bind env prefix failed")
	}
	// 测试环境变量.替换为_
	os.Setenv("ENVPREFIX_OPTION_4", ".o_")
	if viper.GetString("option.4") != ".o_" {
		t.Fatal("env replacer error")
	}

	// 测试加载配置文件
	if viper.GetString("server.bind") == "" {
		t.Fatal("read conf file error")
	}
}

func TestNewViperOption(t *testing.T) {
	o := NewViperOption("name", "default value", "desc")
	if o.Name != "name" || o.Default != "default value" || o.Desc != "desc" {
		t.Fatal(o)
	}
}
