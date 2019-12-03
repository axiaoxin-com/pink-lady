// Package utils save the package of third party package tools and the general tool code written by yourself
// write your general tool in the package
package utils

import (
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// GetLocalIP 获取当前IP
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// StructToURLValues 将结构体指针对象转换为url.Values，key为json tag，value为结构体字段值，没有json tag则使用字段名称
func StructToURLValues(i interface{}) (values url.Values) {
	values = url.Values{}

	iv := reflect.ValueOf(i).Elem() // Elem() 则i必须传指针，不使用Elem() 则不传递指针
	it := iv.Type()
	for i := 0; i < iv.NumField(); i++ {
		vf := iv.Field(i)
		tf := it.Field(i)

		k := tf.Tag.Get("json")
		if k == "" {
			k = tf.Name
		}

		v := ""
		switch vf.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(vf.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(vf.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(vf.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(vf.Float(), 'f', 4, 64)
		case []byte:
			v = string(vf.Bytes())
		case string:
			v = vf.String()
		}

		values.Set(k, v)
	}
	return
}

// CopyFile 复制文件
func CopyFile(sourceFile, destinationFile string) (err error) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Println("[Error] ", err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	return
}

// RemoveAllWhiteSpace 删除字符串中所有的空白符
func RemoveAllWhiteSpace(s string) string {
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}
