package ezLog

import (
	"fmt"
	"strings"
	"time"
)

type ezLogStorage struct {
	ID       int64     `xorm:"pk 'id' autoincr comment('-')"`
	Level    int       `xorm:"index(lv_and_app) index(lv_and_app_and_time) comment('日志等级:1-debug,2-info,3-err,4-ding,5-ding_list,6-ding_all')"`
	AppName  string    `xorm:"char(32) index(lv_and_app) index(lv_and_app_and_time) comment('App名字')"`
	FileName string    `xorm:"char(255) comment('代码文件')"`
	FileLine int       `xorm:"comment('代码行')"`
	Tag      string    `xorm:"char(64) index index(lv_and_app) index(lv_and_app_and_time) comment('日志标签')"`
	Time     time.Time `xorm:"index(lv_and_app_and_time) comment('日志时间')"`
	Content  string    `xorm:"text comment('具体日志')"`
}

func (e ezLogStorage) TableName() string {
	t := time.Now()
	return fmt.Sprintf("logs_of_%v_%v_%v", t.Year(), strings.ToLower(t.Month().String()), t.Day())
}
