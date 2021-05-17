package ezMySQL

import (
	"errors"
	"github.com/Anveena/RoomOfRequirement/ezPasswordEncoder"
	_ "github.com/go-sql-driver/mysql"
	"sync/atomic"
	"time"
	"xorm.io/xorm"
)

type Info struct {
	Host              string
	Port              string
	Account           string
	PasswordBase64Str string
	Name              string
}

var eng *xorm.Engine
var isInit = int32(0)
var dbQueue = make(chan interface{}, 1000)
var closeSignal = make(chan bool, 10)

func InitEnv(dbInfo *Info, dbModels ...interface{}) error {
	if dbInfo.PasswordBase64Str == "" {
		return errors.New("empty password base64 str")
	}
	password, e := ezPasswordEncoder.GetPasswordFromEncodedStr(dbInfo.PasswordBase64Str)
	if e != nil {
		return e
	}
	if atomic.AddInt32(&isInit, 1) != 1 {
		return errors.New("db env MUST init once")
	}
	var err error
	eng, err = xorm.NewEngine("mysql", dbInfo.Account+
		":"+password+"@tcp("+dbInfo.Host+":"+dbInfo.Port+")/"+dbInfo.Name+"?charset=utf8")
	if err != nil {
		return errors.New("db engine init failed ,err:" + err.Error())
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	eng.SetTZLocation(loc)
	eng.SetTZDatabase(loc)
	for _, ptr := range dbModels {
		err = eng.Sync2(ptr)
		if err != nil {
			return err
		}
	}
	go start()
	return nil
}
func Engine() *xorm.Engine {
	return eng
}
func start() {
	var s *xorm.Session
Outer:
	for true {
		s = eng.NewSession()
	Inner:
		for i := 0; i < 1000; i++ {
			select {
			case model := <-dbQueue:
				if _, err := s.InsertOne(model); err != nil {
					break Inner
				}
				break
			case <-closeSignal:
				s.Close()
				break Outer
			}
		}
		s.Close()
	}
	_ = eng.Close()
	atomic.StoreInt32(&isInit, 0)
}
func StopDBQueue() {
	closeSignal <- true
}
func Insert(model interface{}) {
	dbQueue <- model
}
