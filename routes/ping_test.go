package routes

import (
	"testing"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/pink-lady/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	viper.Set("basic_auth.username", "admin")
	viper.Set("basic_auth.password", "admin")
	viper.Set("env", "localhost")
	viper.Set("mysql.localhost.dbname.dsn", "root:roooooot@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC")
	defer viper.Reset()

	models.Init()
	InitRouter(r)
	recorder, err := goutils.RequestHTTPHandler(
		r,
		"GET",
		"/x/ping",
		nil,
		map[string]string{"Authorization": "Basic YWRtaW46YWRtaW4="},
	)
	assert.Nil(t, err)
	assert.Equal(t, recorder.Code, 200)
}
