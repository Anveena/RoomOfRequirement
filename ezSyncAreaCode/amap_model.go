package ezSyncAreaCode

import "time"

type aMapRspModel struct {
	Status    string      `json:"status"`
	Info      string      `json:"info"`
	InfoCode  string      `json:"infocode"`
	Count     string      `json:"count"`
	Districts []aMapModel `json:"districts"`
}
type aMapModel struct {
	AdCode    string      `json:"adcode"`
	Name      string      `json:"name"`
	Center    string      `json:"center"`
	Level     string      `json:"level"`
	Districts []aMapModel `json:"districts"`
}
type areaDBModel struct {
	Code       string    `xorm:"pk char(8) comment('区域编码')"`
	Name       string    `xorm:"char(255) comment('区域名称')"`
	Center     string    `xorm:"char(64) comment('区域中心点坐标')"`
	ParentCode string    `xorm:"char(8) index comment('区域编码')"`
	ParentName string    `xorm:"char(255) comment('区域名称')"`
	CreateTime time.Time `xorm:"updated"`
}
