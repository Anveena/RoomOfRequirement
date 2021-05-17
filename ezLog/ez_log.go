package ezLog

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Anveena/RoomOfRequirement/ezFile"
	"github.com/Anveena/RoomOfRequirement/ezPasswordEncoder"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"xorm.io/xorm"
)

const (
	LogLvDebug       = 0x1
	LogLvInfo        = 0x2
	LogLvError       = 0x3
	LogLvDingMessage = 0x4
	LogLvDingLists   = 0x5
	LogLvDingAll     = 0x6
)

var lvHeaderMap = map[int]string{
	LogLvDebug:       "\n[Debug] ",
	LogLvInfo:        "\n[Info]  ",
	LogLvError:       "\n[Error] ",
	LogLvDingMessage: "\n[Ding]  ",
	LogLvDingLists:   "\n[DAL!]  ",
	LogLvDingAll:     "\n[DAA!]  ",
}
var logFmtStr = "%vFile:%v%vLine:%v%vTime:%v%v%v\n"

type EZLoggerModel struct {
	LogLevel     int
	AppName      string
	ConsoleModel struct {
		Enable bool
	}
	TxtFileModel struct {
		Enable bool
		Path   string
	}
	DingTalkModel struct {
		Enable                 bool
		SecretKey              string
		SecretKeyEncodedString string
		URL                    string
		URLEncodedString       string
		Mobiles                []string
	}
	MySQLModel struct {
		Enable            bool
		Host              string
		Port              string
		Account           string
		PasswordBase64Str string
		TableName         string
		DatabaseName      string
	}
}

var (
	enableConsole = false
	enableMySQL   = false
	enableTxt     = false
	enableDing    = false
	disableAll    = false
	logLevel      = LogLvInfo
	logFileLock   = &sync.Mutex{}
	logFile       *os.File
	dbEngine      *xorm.Engine
	ezLoggerModel *EZLoggerModel
	appName       = ""
	dbQueue       = make(chan *ezLogStorage, 2048)
)

type dingRequestModel struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func SetUpEnv(m *EZLoggerModel) error {
	ezLoggerModel = m
	appName = m.AppName
	logLevel = m.LogLevel
	enableConsole = m.ConsoleModel.Enable
	enableTxt = m.TxtFileModel.Enable
	enableMySQL = m.MySQLModel.Enable
	enableDing = m.DingTalkModel.Enable
	disableAll = !enableConsole && !enableTxt && !enableMySQL
	if appName == "" {
		return errors.New("app name should not be empty")
	}
	if enableDing {
		if ezLoggerModel.DingTalkModel.URLEncodedString == "" {
			return errors.New("empty url encoded str")
		}
		dingUrl, e := ezPasswordEncoder.GetPasswordFromEncodedStr(ezLoggerModel.DingTalkModel.URLEncodedString)
		if e != nil {
			return e
		}
		ezLoggerModel.DingTalkModel.URL = dingUrl
		if ezLoggerModel.DingTalkModel.SecretKeyEncodedString == "" {
			return errors.New("empty secret key encoded str")
		}
		secretKey, e := ezPasswordEncoder.GetPasswordFromEncodedStr(ezLoggerModel.DingTalkModel.SecretKeyEncodedString)
		if e != nil {
			return e
		}
		ezLoggerModel.DingTalkModel.SecretKey = secretKey
	}
	if enableTxt {
		if err := ezFile.CreateDir(ezLoggerModel.TxtFileModel.Path); err != nil {
			return err
		}
		go newLogFile(ezLoggerModel.TxtFileModel.Path)
	}
	if enableMySQL {
		if ezLoggerModel.MySQLModel.PasswordBase64Str == "" {
			return errors.New("empty password base64 str")
		}
		password, e := ezPasswordEncoder.GetPasswordFromEncodedStr(ezLoggerModel.MySQLModel.PasswordBase64Str)
		if e != nil {
			return e
		}
		var err error
		dbEngine, err = xorm.NewEngine("mysql", ezLoggerModel.MySQLModel.Account+
			":"+password+"@tcp("+ezLoggerModel.MySQLModel.Host+":"+ezLoggerModel.MySQLModel.Port+")/"+ezLoggerModel.MySQLModel.DatabaseName+"?charset=utf8")
		if err != nil {
			return errors.New("db engine init failed ,err:" + err.Error())
		}
		loc, _ := time.LoadLocation("Asia/Shanghai")
		dbEngine.SetTZLocation(loc)
		dbEngine.SetTZDatabase(loc)
		if err = dbEngine.Sync2(&ezLogStorage{}); err != nil {
			return err
		}
		go startDB()
	}
	return nil
}
func D(msg ...interface{}) {
	ezlog(LogLvDebug, msg...)
}
func I(msg ...interface{}) {
	ezlog(LogLvInfo, msg...)
}
func E(msg ...interface{}) {
	ezlog(LogLvError, msg...)
}
func DingMessage(msg ...interface{}) {
	ezlog(LogLvDingMessage, msg...)
	if ezLoggerModel.DingTalkModel.Enable {
		sendToDing(LogLvDingMessage, "no tag", fmt.Sprintln(msg...))
	}
}
func DingAtAll(msg ...interface{}) {
	ezlog(LogLvDingAll, msg...)
	if ezLoggerModel.DingTalkModel.Enable {
		sendToDing(LogLvDingAll, "no tag", fmt.Sprintln(msg...))
	}
}
func DingList(msg ...interface{}) {
	ezlog(LogLvDingLists, msg...)
	if ezLoggerModel.DingTalkModel.Enable {
		sendToDing(LogLvDingLists, "no tag", fmt.Sprintln(msg...))
	}
}

func DWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDebug, tag, msg...)
}
func IWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvInfo, tag, msg...)
}
func EWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvError, tag, msg...)
}
func DingMessageWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingMessage, tag, msg...)
	if ezLoggerModel.DingTalkModel.Enable {
		sendToDing(LogLvDingMessage, tag, fmt.Sprintln(msg...))
	}
}
func DingAtAllWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingAll, tag, msg...)
	if ezLoggerModel.DingTalkModel.Enable {
		sendToDing(LogLvDingAll, tag, fmt.Sprintln(msg...))
	}
}
func DingListWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingLists, tag, msg...)
	if ezLoggerModel.DingTalkModel.Enable {
		sendToDing(LogLvDingLists, tag, fmt.Sprintln(msg...))
	}
}
func ezlog(level int, msg ...interface{}) {
	if disableAll || level < logLevel {
		return
	}
	_, file, line, _ := runtime.Caller(2)
	if enableConsole || enableTxt {
		header := lvHeaderMap[level]
		rs := fmt.Sprintf(logFmtStr, header, file, header, line, header, time.Now().Format("15:04:05.999999"), header,
			strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
		if enableConsole {
			println(rs)
		}
		if enableTxt {
			logFileLock.Lock()
			_, _ = logFile.WriteString(rs)
			if err := logFile.Sync(); err != nil {
				println(err.Error())
			}
			logFileLock.Unlock()
		}
	}
	if enableMySQL {
		dbQueue <- &ezLogStorage{
			Level:    level,
			AppName:  appName,
			FileName: file,
			FileLine: line,
			Tag:      "no tag",
			Time:     time.Now(),
			Content:  fmt.Sprintln(msg...),
		}
	}
}
func ezlogWithTag(level int, tag string, msg ...interface{}) {
	if disableAll || level < logLevel {
		return
	}
	_, file, line, _ := runtime.Caller(2)
	if enableTxt || enableConsole {
		header := lvHeaderMap[level] + "[" + tag + "] "
		rs := fmt.Sprintf(logFmtStr, header, file, header, line, header, time.Now().Format("15:04:05.999999"), header,
			strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
		if enableConsole {
			println(rs)
		}
		if enableTxt {
			logFileLock.Lock()
			_, _ = logFile.WriteString(rs)
			if err := logFile.Sync(); err != nil {
				println(err.Error())
			}
			logFileLock.Unlock()
		}
	}
	if enableMySQL {
		dbQueue <- &ezLogStorage{
			Level:    level,
			AppName:  appName,
			FileName: file,
			FileLine: line,
			Tag:      tag,
			Time:     time.Now(),
			Content:  fmt.Sprintln(msg...),
		}
	}
}
func newLogFile(fp string) {
	runtime.LockOSThread()
	for {
		now := time.Now()
		dir := filepath.Join(fp, strconv.Itoa(now.Year()), now.Month().String())
		if logFile != nil {
			if err := logFile.Close(); err != nil {
				println(err.Error())
			}
		}
		var err error
		logFileLock.Lock()
		logFile, err = ezFile.CreateFile(dir, strconv.Itoa(now.Day())+".log", true, os.O_RDWR|os.O_CREATE|os.O_APPEND)
		logFileLock.Unlock()
		if err != nil {
			println(err.Error())
			return
		}
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}
}

func getHmacSHA256Base64String(s string) (string, error) {
	h := hmac.New(sha256.New, []byte(ezLoggerModel.DingTalkModel.SecretKey))
	if _, err := h.Write([]byte(s)); err != nil {
		return "", err
	}
	rs := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(rs), nil
}
func sendToDing(logLv int, tag string, msg string) {
	_, file, line, _ := runtime.Caller(2)
	fileSubArr := strings.SplitN(file, "/src/", 2)
	if len(fileSubArr) == 2 {
		file = fileSubArr[1]
	}
	realMsg := strings.ReplaceAll(msg, "\n", "\n>\n>")
	atInfo := struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	}{}
	atStr := ""
	if logLv == LogLvDingLists && len(ezLoggerModel.DingTalkModel.Mobiles) == 0 {
		logLv = LogLvDingAll
	}
	if logLv == LogLvDingAll {
		atInfo.IsAtAll = true
	} else if logLv == LogLvDingLists {
		atStr += "### **责任人**:"
		for _, people := range ezLoggerModel.DingTalkModel.Mobiles {
			atStr += " @" + people
		}
		atStr += "\n"
		atInfo.AtMobiles = ezLoggerModel.DingTalkModel.Mobiles
	}
	m := dingRequestModel{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: tag + "@" + appName,
			Text:  fmt.Sprintf("# %s@%s\n### **File**:%s\n### **Line**:%d\n### **Time**:%s\n### **Message**:\n>%s \n%s", tag, appName, file, line, time.Now().Format("15:04:05.999999"), realMsg, atStr),
		},
		At: atInfo,
	}
	tmpData, err := json.Marshal(m)
	if err != nil {
		E(err.Error())
		return
	}
	timestampString := strconv.Itoa(int(time.Now().Unix() * 1000))
	sign, err := getHmacSHA256Base64String(timestampString + "\n" + ezLoggerModel.DingTalkModel.SecretKey)
	if err != nil {
		E(err.Error())
		return
	}
	go func() {
		rsp, err := http.Post(ezLoggerModel.DingTalkModel.URL+"&timestamp="+timestampString+"&sign="+url.QueryEscape(sign), "application/json;charset=utf-8", bytes.NewBuffer(tmpData))
		if err != nil {
			E(err.Error())
			return
		}
		rspInfo, _ := ioutil.ReadAll(rsp.Body)
		rspStr := string(rspInfo)
		if !strings.Contains(rspStr, `"errcode":0,`) {
			E("rsp from ding:" + rspStr)
		}
		_ = rsp.Body.Close()
	}()
}

func startDB() {
	var err error
	timer := time.NewTicker(time.Second * 1)
	msgArr := make([]*ezLogStorage, 1024)
	for true {
		i := 0
	outer:
		for ; i < 1024; i++ {
			select {
			case <-timer.C:
				break outer
			case msgArr[i] = <-dbQueue:
				break
			}
		}
		if i > 0 {
			sess := dbEngine.NewSession()
			if _, err = sess.InsertMulti(msgArr[:i]); err != nil {
				DingAtAllWithTag("db error", err.Error())
			}
			_ = sess.Close()
		}
	}
}
