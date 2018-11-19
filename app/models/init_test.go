package models

import (
	"testing"
	"time"

	"github.com/axiaoxin/pink-lady/app/utils"
)

func TestMigrate(t *testing.T) {
	type S struct {
		BaseModel
		X int
	}
	MigrationModels = append(MigrationModels, &S{})
	err := utils.InitGormDB("sqlite3", "", "/tmp/pink-lady-unit-test.db", "", "", 0, 0, 0, true)
	if err != nil {
		t.Error("init gorm db error:", err)
	}
	defer utils.DB.Close()
	if err := Migrate(); err != nil {
		t.Error(err)
	}
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
