// Package utils save the package of third party package tools and the general tool code written by yourself
// write your general tool in the package
package utils

import (
	"net"
	"net/url"
	"reflect"
	"strconv"
)

// init function here for utils package
func init() {

}

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

// StructToURLValues 将结构体转换为url.Values，key为json tag，没有json tag则使用字段名称，传入指针参数
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
