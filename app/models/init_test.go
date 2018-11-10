package models

import (
	"testing"
	"time"

	"github.com/axiaoxin/gin-skeleton/app/utils"
)

func TestModelMigrate(t *testing.T) {
	type S struct {
		BaseModel
		X int
	}
	Models = append(Models, &S{})
	err := utils.InitGormDB("sqlite3", "", "/tmp/gin-skeleton-unit-test.db", "", "", 0, 0, 0, true)
	if err != nil {
		t.Error("init gorm db error:", err)
	}
	defer utils.DB.Close()
	Migrate()
	if err := utils.DB.Exec("select * from s").Error; err != nil {
		t.Error(err)
	}
	s := S{X: 666}
	if err := utils.DB.Create(&s).Error; err != nil {
		t.Error(err)
	}
	s = S{}
	if err := utils.DB.Last(&s).Error; err != nil {
		t.Error(err)
	}
	if s.X != 666 {
		t.Error("x != 666")
	}
	d := s.CreatedAt.String()
	_, err = time.ParseInLocation(utils.TimeFormat, d, time.Local)
	if err != nil {
		t.Error(err)
	}
}
