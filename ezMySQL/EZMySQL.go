package ezMySQL

import (
	"encoding/base64"
	"errors"
	"github.com/Anveena/RoomOfRequirement/ezCrypto"
	_ "github.com/go-sql-driver/mysql"
	"sync/atomic"
	"time"
	"xorm.io/xorm"
)

type Info struct {
	Host              string `json:"host"`
	Port              string `json:"port"`
	Account           string `json:"account"`
	PasswordBase64Str string `json:"password_base_64_str"`
	Name              string `json:"name"`
}

var eng *xorm.Engine
var isInit = int32(0)
var dbQueue = make(chan interface{}, 1000)
var closeSignal = make(chan bool, 10)

func MakePasswordBase64Str(origPwd string) (string, error) {
	origPwdData := []byte(origPwd)
	encData, err := ezCrypto.EZEncrypt(&origPwdData, "this code may be not working", 9458)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(*encData), nil
}
func getPasswordFromBase64Str(base64Str string) (string, error) {
	encData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	origData, err := ezCrypto.EZDecrypt(&encData, "this code may be not working", 9458)
	if err != nil {
		return "", err
	}
	return string(*origData), nil
}
func InitEnv(dbInfo *Info, dbModels ...interface{}) error {
	if dbInfo.PasswordBase64Str == "" {
		return errors.New("empty password base64 str")
	}
	password, e := getPasswordFromBase64Str(dbInfo.PasswordBase64Str)
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
