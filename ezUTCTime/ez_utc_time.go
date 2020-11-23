package ezUTCTime

import (
	"github.com/beevik/ntp"
	"time"
)

var protobufTime int64

func SyncTimeFromAliyun() error {
	t1 := time.Now().Unix()
	t2, err := ntp.Time("ntp1.aliyun.com")
	if err != nil {
		return err
	}
	protobufTime = t2.Unix() - t1
	return nil
}
func GetAliyunTimestamp() uint64 {
	return uint64(time.Now().Unix()+protobufTime) * 1000
}
