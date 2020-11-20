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
)

const (
	LogLvDebug       = 0x1
	LogLvInfo        = 0x2
	LogLvError       = 0x3
	LogLvDingMessage = 0x4
	LogLvDingLists   = 0x5
	LogLvDingAll     = 0x6
)

var lvMap = map[int]string{
	LogLvDebug:       "[Debug] ",
	LogLvInfo:        "[Info]  ",
	LogLvError:       "[Error] ",
	LogLvDingMessage: "[Ding]  ",
	LogLvDingLists:   "[DAL!]  ",
	LogLvDingAll:     "[DAA!]  ",
}
var logLevel = LogLvInfo
var logFile *os.File
var logFileLock = &sync.Mutex{}
var debugMode = false
var dingOBJ *dingTalkModel
var appName = ""

type dingTalkModel struct {
	Enable    bool
	Mobiles   []string
	SecretKey string
	URL       string
}
type EZLoggerModel struct {
	LogLevel      int
	LogFilePath   string
	AppName       string
	DebugMode     bool
	DingTalkModel dingTalkModel
}
type dingModel struct {
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
	appName = m.AppName
	if appName == "" {
		return errors.New("app name should not be empty")
	}
	dingOBJ = &m.DingTalkModel
	if dingOBJ.Enable {
		if dingOBJ.URL == "" {
			return errors.New("ding talk is enabled but ding url is empty")
		}
	}
	logLevel = m.LogLevel
	debugMode = m.DebugMode
	if !m.DebugMode {
		if err := ezFile.CreateDir(m.LogFilePath); err != nil {
			return err
		}
		go newLogFile(m.LogFilePath)
	}
	return nil
}
func D(msg ...interface{}) {
	ezlog(LogLvDebug, fmt.Sprintln(msg...))
}
func I(msg ...interface{}) {
	ezlog(LogLvInfo, fmt.Sprintln(msg...))
}
func E(msg ...interface{}) {
	ezlog(LogLvError, fmt.Sprintln(msg...))
}
func DingMessage(msg ...interface{}) {
	ezlog(LogLvDingMessage, fmt.Sprintln(msg...))
	if dingOBJ.Enable {
		sendToDing(LogLvDingMessage, "no tag", fmt.Sprintln(msg...))
	}
}
func DingAtAll(msg ...interface{}) {
	ezlog(LogLvDingAll, fmt.Sprintln(msg...))
	if dingOBJ.Enable {
		sendToDing(LogLvDingAll, "no tag", fmt.Sprintln(msg...))
	}
}
func DingList(msg ...interface{}) {
	ezlog(LogLvDingLists, fmt.Sprintln(msg...))
	if dingOBJ.Enable {
		sendToDing(LogLvDingLists, "no tag", fmt.Sprintln(msg...))
	}
}

func DWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDebug, tag, fmt.Sprintln(msg...))
}
func IWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvInfo, tag, fmt.Sprintln(msg...))
}
func EWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvError, tag, fmt.Sprintln(msg...))
}
func DingMessageWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingMessage, tag, fmt.Sprintln(msg...))
	if dingOBJ.Enable {
		sendToDing(LogLvDingMessage, tag, fmt.Sprintln(msg...))
	}
}
func DingAtAllWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingAll, tag, fmt.Sprintln(msg...))
	if dingOBJ.Enable {
		sendToDing(LogLvDingAll, tag, fmt.Sprintln(msg...))
	}
}
func DingListWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingLists, tag, fmt.Sprintln(msg...))
	if dingOBJ.Enable {
		sendToDing(LogLvDingLists, tag, fmt.Sprintln(msg...))
	}
}
func ezlog(level int, msg string) {
	if level >= logLevel {
		rs := msgFmt(level, msg)
		if debugMode {
			println(rs)
		} else {
			logFileLock.Lock()
			_, _ = logFile.WriteString(rs)
			if err := logFile.Sync(); err != nil {
				println(err.Error())
			}
			logFileLock.Unlock()
		}
	}
}
func ezlogWithTag(level int, tag string, msg string) {
	if level >= logLevel {
		rs := msgWithTagFmt(level, tag, msg)
		if debugMode {
			println(rs)
		} else {
			logFileLock.Lock()
			_, _ = logFile.WriteString(rs)
			if err := logFile.Sync(); err != nil {
				println(err.Error())
			}
			logFileLock.Unlock()
		}
	}
}
func msgWithTagFmt(level int, tag string, msg string) string {
	_, file, line, _ := runtime.Caller(3)
	strHeader := "\n" + lvMap[level] + "[" + tag + "] "
	msgFormatted := strings.ReplaceAll(msg, "\n", strHeader)
	return strHeader + "File:" + file +
		strHeader + "Line:" + strconv.Itoa(line) +
		strHeader + "Time:" + time.Now().Format("15:04:05.999999") +
		strHeader + msgFormatted + "\n"
}
func msgFmt(level int, msg string) string {
	_, file, line, _ := runtime.Caller(3)
	strHeader := "\n" + lvMap[level] + " "
	msgFormatted := strings.ReplaceAll(msg, "\n", strHeader)
	return strHeader + "File:" + file +
		strHeader + "Line:" + strconv.Itoa(line) +
		strHeader + "Time:" + time.Now().Format("15:04:05.999999") +
		strHeader + msgFormatted + "\n"
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
	h := hmac.New(sha256.New, []byte(dingOBJ.SecretKey))
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
	if logLv == LogLvDingLists && len(dingOBJ.Mobiles) == 0 {
		logLv = LogLvDingAll
	}
	if logLv == LogLvDingAll {
		atInfo.IsAtAll = true
	} else if logLv == LogLvDingLists {
		atStr += "### **责任人**:"
		for _, people := range dingOBJ.Mobiles {
			atStr += " @" + people
		}
		atStr += "\n"
		atInfo.AtMobiles = dingOBJ.Mobiles
	}
	m := dingModel{
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
	sign, err := getHmacSHA256Base64String(timestampString + "\n" + dingOBJ.SecretKey)
	if err != nil {
		E(err.Error())
		return
	}
	go func() {
		rsp, err := http.Post(dingOBJ.URL+"&timestamp="+timestampString+"&sign="+url.QueryEscape(sign), "application/json;charset=utf-8", bytes.NewBuffer(tmpData))
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
