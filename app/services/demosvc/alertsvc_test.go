package demosvc

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"pink-lady/app/database"
	"pink-lady/app/models/demomod"
	"pink-lady/app/router"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var utdb *gorm.DB
var utctx *gin.Context

func Setup() {
	// 初始化配置
	router.InitDependencies("../../", "config")
	utdb = database.UTDB()
	utdb.AutoMigrate(&demomod.AlertPolicy{}, &demomod.AlertFilterRule{}, &demomod.AlertTriggerRule{})
	utdb.LogMode(false)

	// 初始化gin.Context
	utctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	utctx.Request, _ = http.NewRequest("POST", "/", nil)
	log.Println("Setup!")
}

func Teardown() {
	// 关闭db
	utdb.Close()
	// Reset 变量
	utdb = nil
	utctx = nil
	os.Remove(database.UTDBFile)
	log.Println("Teardown!")
}

func CreateATestPolicy() (int64, error) {
	// 新增一条记录
	policy := &demomod.AlertPolicy{
		AppID:              -1,
		Uin:                "-1",
		Name:               fmt.Sprint("test_", time.Now().Unix(), rand.Float64()),
		MetricSetID:        -1,
		NoticeFrequencySec: 60,
		Status:             demomod.AlertPolicyStatusEnable,
		AlertGroupID:       "-1,0",
		AlertChannel:       "weixin,email",
		NoticePeriodBegin:  0,
		NoticePeriodEnd:    86399,
		URLScheme:          "http",
		CallbackURL:        "-",
		AlertFilterRules: []*demomod.AlertFilterRule{
			{
				AlertPolicyID: -1,
				Relation:      1,
				Field:         "-",
				Operating:     "=",
				Value:         "-",
			},
		},
		AlertTriggerRules: []*demomod.AlertTriggerRule{
			{
				Relation:             1,
				MetricID:             -1,
				MetricType:           1,
				Operating:            "=",
				Value:                "-",
				ContinuousCycleCount: 1,
			},
			{
				Relation:             1,
				MetricID:             -2,
				MetricType:           1,
				Operating:            "=",
				Value:                "-",
				ContinuousCycleCount: 1,
			},
		},
	}
	return CreateAlertPolicy(utctx, utdb, policy)
}

func TestCreateAlertPolicy(t *testing.T) {
	Setup()
	defer Teardown()

	// 测试新增告警策略
	policy := &demomod.AlertPolicy{
		AppID:              -1,
		Uin:                "-1",
		Name:               fmt.Sprint("test-", time.Now().Unix(), rand.Float64()),
		MetricSetID:        -1,
		NoticeFrequencySec: 60,
		Status:             demomod.AlertPolicyStatusEnable,
		AlertGroupID:       "-1,0",
		AlertChannel:       "weixin,email",
		NoticePeriodBegin:  0,
		NoticePeriodEnd:    86399,
		URLScheme:          "http",
		CallbackURL:        "-",
		AlertFilterRules: []*demomod.AlertFilterRule{
			{
				Relation:  1,
				Field:     "-",
				Operating: "=",
				Value:     "-",
			},
		},
		AlertTriggerRules: []*demomod.AlertTriggerRule{
			{
				Relation:             1,
				MetricID:             -1,
				MetricType:           1,
				Operating:            "=",
				Value:                "-",
				ContinuousCycleCount: 1,
			},
			{
				Relation:             1,
				MetricID:             -2,
				MetricType:           1,
				Operating:            "=",
				Value:                "-",
				ContinuousCycleCount: 1,
			},
		},
	}
	id, err := CreateAlertPolicy(utctx, utdb, policy)
	if err != nil {
		t.Fatal("TestCreateAlertPolicy 失败:", err)
	}
	t.Log("TestCreateAlertPolicy 添加新记录id:", id)
	if id == 0 {
		t.Fatal("TestCreateAlertPolicy 没有返回正确的策略ID")
	}
	// 测试告警过滤条件是否插入成功
	alertFilterRules := []demomod.AlertFilterRule{}
	if err := utdb.Find(&alertFilterRules, demomod.AlertFilterRule{AlertPolicyID: id}).Error; err != nil {
		t.Fatal("查询告警过滤条件失败：", err)
	}
	flen := len(alertFilterRules)
	if flen != 1 {
		t.Fatal("TestCreateAlertPolicy 告警过滤条件插入失败 flen:", flen)
	}
	// 测试告警触发条件是否插入成功
	alertTriggerRules := []demomod.AlertTriggerRule{}
	if err := utdb.Find(&alertTriggerRules, demomod.AlertTriggerRule{AlertPolicyID: id}).Error; err != nil {
		t.Fatal("TestCreateAlertPolicy 查询告警触发条件失败：", err)
	}
	tlen := len(alertTriggerRules)
	if tlen != 2 {
		t.Fatal("TestCreateAlertPolicy 告警触发条件插入失败 tlen:", tlen)
	}
}

func TestDescribeAlertPolicy(t *testing.T) {
	Setup()
	defer Teardown()

	// 新增一条记录
	id, err := CreateATestPolicy()
	if err != nil {
		t.Fatal("TestDescribeAlertPolicy 添加新记录失败:", err)
	}
	t.Log("TestDescribeAlertPolicy 添加新记录id:", id)

	// 按新增的id指定查询参数查询该记录
	p, err := DescribeAlertPolicy(utctx, utdb, -1, "-1", id)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicy 查询错误")
	}
	if p == nil {
		t.Fatal("TestDescribeAlertPolicy 告警策略不存在")
	}
	if p.Name == "" {
		t.Fatal("TestDescribeAlertPolicy 查询name字段是空值")
	}
	flen := len(p.AlertFilterRules)
	if flen != 1 {
		t.Fatal("TestDescribeAlertPolicy 过滤条件没有正确获取 flen:", flen)
	}
	tlen := len(p.AlertTriggerRules)
	if tlen != 2 {
		t.Fatal("TestDescribeAlertPolicy 告警条件没有正确获取 tlen:", tlen)
	}

	// 测试查询一条不存在的记录
	p, err = DescribeAlertPolicy(utctx, utdb, -1, "-1", int64(-1))
	if err == nil {
		t.Fatal("TestDescribeAlertPolicy 查询不存在的记录应该返回错误")
	}
}

func TestRules(t *testing.T) {
	Setup()
	defer Teardown()

	id, err := CreateATestPolicy()
	if err != nil {
		t.Fatal("TestRules 添加新记录失败:", err)
	}
	t.Log("TestRules 添加新记录id:", id)

	// 测试筛选条件
	frules, err := DescribeAlertFilterRules(utctx, utdb, -1, "-1", id)
	if err != nil {
		t.Fatal("TestRules 获取告警策略筛选条件错误：", err)
	}
	rlen := len(frules)
	if rlen != 1 {
		t.Fatal("TestRules 筛选规则数量错误 rlen:", rlen)
	}

	// 测试触发条件
	trules, err := DescribeAlertTriggerRules(utctx, utdb, id)
	if err != nil {
		t.Fatal("TestRules 获取告警策略触发条件错误：", err)
	}
	tlen := len(trules)
	if tlen != 2 {
		t.Fatal("TestRules 触发规则数量错误 tlen:", tlen)
	}

}

func TestDeleteAlertPolicy(t *testing.T) {
	Setup()
	defer Teardown()

	var err error

	id, err := CreateATestPolicy()
	if err != nil {
		t.Fatal("TestDeleteAlertPolicy 添加新记录失败:", err)
	}
	t.Log("TestDeleteAlertPolicy 添加新记录id:", id)

	err = DeleteAlertPolicy(utctx, utdb, -1, "-1", id)
	if err != nil {
		t.Fatal("TestDeleteAlertPolicy 删除记录失败:", err)
	}

	// 测试删除的该记录是否还存在
	policy, err := DescribeAlertPolicy(utctx, utdb, -1, "-1", id)
	if err == nil {
		t.Fatal("TestDeleteAlertPolicy 执行后获取该删除记录应该报记录不存在的错误")
	}
	if policy != nil {
		t.Fatal("TestDeleteAlertPolicy 执行删除后策略依然存在")
	}

	// 测试该记录关联的筛选条件是否还存在
	rules, err := DescribeAlertFilterRules(utctx, utdb, -1, "-1", id)
	if len(rules) != 0 {
		t.Fatal("TestDeleteAlertPolicy 筛选条件没有成功删除")
	}

	// 测试该记录关联的触发条件是否还存在
	trules, err := DescribeAlertTriggerRules(utctx, utdb, id)
	if len(trules) != 0 {
		t.Fatal("TestDeleteAlertPolicy 筛选条件没有成功删除")
	}

	// 测试删除不存在的记录
	err = DeleteAlertPolicy(utctx, utdb, -1, "-1", int64(-1))
	if err == nil {
		t.Fatal("TestDeleteAlertPolicy 删除不存在的记录应该返回错误")
	}
}

func TestModifyAlertPolicy(t *testing.T) {
	Setup()
	defer Teardown()

	var err error

	id, err := CreateATestPolicy()
	if err != nil {
		t.Fatal("TestModifyAlertPolicy 添加新记录失败:", err)
	}
	t.Log("TestModifyAlertPolicy 添加新记录id:", id)
	p, err := DescribeAlertPolicy(utctx, utdb, -1, "-1", id)
	p.AlertChannel = ""
	// 删除全部过滤条件
	p.AlertFilterRules = []*demomod.AlertFilterRule{}
	// 删除第一个触发条件，修改第二个触发条件，新增一个触发条件
	tr2 := p.AlertTriggerRules[1]
	tr2.Relation = 2
	ntr := &demomod.AlertTriggerRule{
		Relation:             1,
		MetricID:             -2,
		MetricType:           1,
		Operating:            "=",
		Value:                "x",
		ContinuousCycleCount: 1,
	}
	p.AlertTriggerRules = []*demomod.AlertTriggerRule{tr2, ntr}
	mPolicy, err := ModifyAlertPolicy(utctx, utdb, p)
	if err != nil {
		t.Fatal("ModifyAlertPolicy err:", err)
	}
	qPolicy, err := DescribeAlertPolicy(utctx, utdb, -1, "-1", id)
	if err != nil {
		t.Fatal("TestModifyAlertPolicy 获取被修改的id失败")
	}
	// 检查修改是否生效
	if qPolicy.AlertChannel != "" {
		t.Fatal("TestModifyAlertPolicy 修改字段未生效", qPolicy.AlertChannel)
	}
	// 检查过滤条件是否被删除
	if len(qPolicy.AlertFilterRules) != 0 {
		t.Fatal("TestModifyAlertPolicy 更新时清空过滤条件没有成功")
	}
	if len(qPolicy.AlertTriggerRules) != 2 {
		t.Fatal("TestModifyAlertPolicy 更新后的触发条件应该为2个")
	}
	if qPolicy.AlertTriggerRules[0].Relation != 2 {
		t.Fatal("TestModifyAlertPolicy 更新触发条件字段失败")
	}
	if qPolicy.AlertTriggerRules[1].Value != "x" {
		t.Fatal("TestModifyAlertPolicy 更新时新增的触发条件错误")
	}

	// 检查修改返回的结果和查询修改后的结果是否一致
	if qPolicy.AlertChannel != mPolicy.AlertChannel {
		t.Fatal("TestModifyAlertPolicy 对比更新返回的policy和查询的policy字段失败")
	}
	if qPolicy.AlertTriggerRules[0].Relation != mPolicy.AlertTriggerRules[0].Relation {
		t.Fatal("TestModifyAlertPolicy 对比更新返回的关联条件的字段失败")
	}
}

var nowCount int
var testIDs []int64

func TestDescribeAlertPolicies1(t *testing.T) {
	Setup()

	// 新创建10条新记录确保一定有数据用于测试列表操作 （可能已插入过数据）
	for i := 0; i < 10; i++ {
		id, err := CreateATestPolicy()
		if err != nil {
			t.Fatal("TestDescribeAlertPolicies 添加新记录失败:", err)
		}
		testIDs = append(testIDs, id)
	}
	// 获取当前数据总数
	utdb.Model(&demomod.AlertPolicy{}).Where("appid = ? AND uin = ?", -1, "-1").Count(&nowCount)

	// 测试过滤条件使用空值查询
	offset := 0
	limit := 0
	order := ""
	id := int64(0)
	name := ""
	list, count, err := DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count != nowCount {
		t.Fatal("TestDescribeAlertPolicies 过滤条件为空值count应为全部数据的总数", count, nowCount)
	}
	if len(list) != 0 {
		t.Fatal("TestDescribeAlertPolicies limit为0返回list应该为空")
	}
}

func TestDescribeAlertPolicies2(t *testing.T) {
	// 测试limit
	id := int64(0)
	name := ""
	offset := 0
	limit := 5
	order := ""
	list, count, err := DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count != nowCount {
		t.Fatal("TestDescribeAlertPolicies 过滤条件为空值count应为全部数据的总数", count, nowCount)
	}
	if len(list) != limit {
		t.Fatal("TestDescribeAlertPolicies limit为10 list返回返回10条数据")
	}
	// 默认排序按id 从小到大
	if list[0].ID > list[1].ID {
		t.Fatal("默认排序应为ID从小到大")
	}
	// offset翻页
	offset = 5
	limit = 5
	olist, count, err := DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count != nowCount {
		t.Fatal("TestDescribeAlertPolicies 过滤条件为空值count应为全部数据的总数", count, nowCount)
	}
	if len(list) != limit {
		t.Fatal("TestDescribeAlertPolicies limit应该和返回的数量一致")
	}
	// 默认排序按id 从小到大
	if olist[0].ID > olist[1].ID {
		t.Fatal("默认排序应为ID从小到大", olist)
	}

	// 测试排序按id 从大到小 拿到的最后10条记录应该全为新加的10条记录
	order = "id desc"
	offset = 0
	limit = 5
	list, count, err = DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count != nowCount {
		t.Fatal("TestDescribeAlertPolicies 过滤条件为空值count应为全部数据的总数", count, nowCount)
	}
	if len(list) != limit {
		t.Fatal("TestDescribeAlertPolicies limit应该和返回的数量一致")
	}
	if list[0].ID < list[1].ID {
		t.Fatal("排序应为ID从大到小")
	}
	// offset 翻页
	offset = 5
	olist, count, err = DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if len(olist) != limit {
		t.Fatal("TestDescribeAlertPolicies limit应该和返回的数量一致")
	}
	if olist[0].ID+1 != list[limit-1].ID {
		t.Fatal("TestDescribeAlertPolicies offset之后的第一个记录的ID应该和offset之前的最后一个记录的ID连续")
	}
}

func TestDescribeAlertPolicies3(t *testing.T) {
	// 测试取消offset和limit
	id := int64(0)
	order := ""
	name := ""
	offset := -1
	limit := -1
	list, count, err := DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count != nowCount {
		t.Fatal("TestDescribeAlertPolicies 过滤条件为空值count应为全部数据的总数", count, nowCount)
	}
	if len(list) != count {
		t.Fatal("TestDescribeAlertPolicies list中的数据应该为全部数据")
	}

	// 测试按id精确搜索
	id = testIDs[9]
	list, count, err = DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count != 1 {
		t.Fatal("TestDescribeAlertPolicies 指定id搜索count应为1", count)
	}
	if len(list) != 1 {
		t.Fatal("TestDescribeAlertPolicies list中的数据应该只有一条记录")
	}

	// 测试按name模糊搜索
	id = 0
	name = "test"
	offset = 0
	limit = 5
	list, count, err = DescribeAlertPolicies(utctx, utdb, -1, "-1", offset, limit, order, id, name)
	if err != nil {
		t.Fatal("TestDescribeAlertPolicies 空值过滤条件时错误:", err)
	}
	if count < 10 {
		t.Fatal("TestDescribeAlertPolicies 模糊查找带_的记录至少10条", count)
	}
	if len(list) != limit {
		t.Fatal("TestDescribeAlertPolicies list中的数据肯定应该和limit相同")
	}
}
