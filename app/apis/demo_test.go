package apis

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"pink-lady/app/database"
	"pink-lady/app/router"
	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var utrouter *gin.Engine

func Setup() {
	// 注册路由
	utrouter = router.SetupRouter("..", "config")
	RegisterRoutes(utrouter)
}

func Teardown() {
	// 关闭数据库连接
	database.InstanceMap.Close()
	utrouter = nil
}

func TestAlertAPI(t *testing.T) {
	Setup()
	defer Teardown()

	// 测试创建接口
	paramCreate := fmt.Sprintf(`{
		"appid": 1,
		"uin": "axiaoxin",
		"alert_policy": {
			"alert_channel": "weixin,sms",
			"alert_filter_rules": [
			{
				"field": "field",
				"operating": "=",
				"relation": 1,
				"value": "value"
			}
			],
			"alert_group_id": "-1,0",
			"alert_trigger_rules": [
			{
				"continuous_cycle_count": 1,
				"metric_id": -1,
				"metric_type": 1,
				"operating": "=",
				"relation": 1,
				"value": "1"
			}
			],
			"callback_url": "axiaoxin.com",
			"latest_alert_time": "string",
			"metric_set_id": -1,
			"name": "ut-%d",
			"notice_frequency_sec": 60,
			"notice_period_begin": 0,
			"notice_period_end": 86399,
			"status": 1,
			"uin": "string",
			"url_scheme": "http"
		}
	}`, time.Now().UnixNano()+rand.Int63())

	respRecorder := utils.PerformRequest(utrouter, "POST", "/demo/alert-policy", []byte(paramCreate))
	body := respRecorder.Body.Bytes()
	data := jsoniter.Get(body)

	if respRecorder.Code != 200 {
		t.Fatal("接口响应错误：", respRecorder.Code, string(body))
	}
	if data.Get("code").ToInt() != 0 {
		t.Fatal("接口返回了错误信息，没有成功创建告警策略:", string(body))
	}
	policyID := data.Get("data").ToInt64()
	if policyID == 0 {
		t.Fatal("接口没有返回正确id:", string(body))
	}

	// 测试查询接口
	respRecorder = utils.PerformRequest(utrouter, "GET", fmt.Sprint("/demo/alert-policy/1/axiaoxin/", policyID), nil)
	body = respRecorder.Body.Bytes()
	data = jsoniter.Get(body)

	if respRecorder.Code != 200 {
		t.Fatal("接口响应错误：", respRecorder.Code, string(body))
	}
	if data.Get("code").ToInt() != 0 {
		t.Fatal("接口返回了错误信息，没有成功获取告警策略:", string(body))
	}
	policyName := data.Get("data", "name").ToString()
	if policyName == "" {
		t.Fatal("接口返回的结果没有name字段:", string(body))
	}

	// 测试修改接口
	paramModify := fmt.Sprintf(`{
		"appid": 1,
		"uin": "axiaoxin",
		"alert_policy": {
			"id": %d,
			"alert_channel": "weixin",
			"alert_filter_rules": [],
			"alert_group_id": "-1,0",
			"alert_trigger_rules": [
			{
				"continuous_cycle_count": 1,
				"metric_id": -1,
				"metric_type": 1,
				"operating": "=",
				"relation": 1,
				"value": "1"
			}
			],
			"callback_url": "axiaoxin.com",
			"latest_alert_time": "string",
			"metric_set_id": -1,
			"name": "%s",
			"notice_frequency_sec": 60,
			"notice_period_begin": 0,
			"notice_period_end": 86399,
			"status": 1,
			"uin": "string",
			"url_scheme": "http"
		}
	}`, policyID, policyName)

	respRecorder = utils.PerformRequest(utrouter, "PUT", "/demo/alert-policy", []byte(paramModify))
	body = respRecorder.Body.Bytes()
	data = jsoniter.Get(body)

	if respRecorder.Code != 200 {
		t.Fatal("接口响应错误：", respRecorder.Code, string(body))
	}
	if data.Get("code").ToInt() != 0 {
		t.Fatal("接口返回了错误信息，没有成功获取告警策略:", string(body))
	}
	if data.Get("data", "alert_channel").ToString() != "weixin" {
		t.Fatal("修改策略表的alert_channel字段失败:", string(body))
	}
	if data.Get("data", "alert_filter_rules").Size() != 0 {
		t.Fatal("修改策略清空过滤条件未生效:", string(body))
	}
	if data.Get("data", "alert_trigger_rules").Size() != 1 {
		t.Fatal("修改策略触发条件未生效:", string(body))
	}

	// 测试列表接口
	// 搜索条件全部为空值（不搜索）
	respRecorder = utils.PerformRequest(utrouter, "GET", "/demo/alert-policy?appid=1&uin=axiaoxin", nil)
	body = respRecorder.Body.Bytes()
	data = jsoniter.Get(body)

	if respRecorder.Code != 200 {
		t.Fatal("接口响应错误：", respRecorder.Code, string(body))
	}
	if data.Get("code").ToInt() != 0 {
		t.Fatal("接口返回了错误信息，没有成功获取告警策略:", string(body))
	}
	if data.Get("data", "total_count").ToInt() == 0 {
		t.Fatal("列表total_count字段错误:", string(body))
	}
	if data.Get("data", "alert_policies").Size() == 0 {
		t.Fatal("列表中没有数据:", string(body))
	}
	// 指定当前ID搜索
	respRecorder = utils.PerformRequest(utrouter, "GET", fmt.Sprint("/demo/alert-policy?appid=1&uin=axiaoxin&id=", policyID), nil)
	body = respRecorder.Body.Bytes()
	data = jsoniter.Get(body)
	if data.Get("data", "alert_policies").Size() != 1 {
		t.Fatal("按ID搜索返回错误结果:", string(body))
	}

	// 测试删除接口
	respRecorder = utils.PerformRequest(utrouter, "DELETE", fmt.Sprint("/demo/alert-policy/1/axiaoxin/", policyID), nil)
	body = respRecorder.Body.Bytes()
	data = jsoniter.Get(body)

	if respRecorder.Code != 200 {
		t.Fatal("接口响应错误：", respRecorder.Code, string(body))
	}
	if data.Get("code").ToInt() != 0 {
		t.Fatal("接口返回了错误信息，没有成功获取告警策略:", string(body))
	}
	if data.Get("data").ToBool() != true {
		t.Fatal("删除失败没有返回true:", string(body))
	}

}
