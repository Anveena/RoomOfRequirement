package ezLog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Anveena/RoomOfRequirement/ezFile"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LogLvTrace         = 0x1
	LogLvInfo          = 0x2
	LogLvWarning       = 0x3
	LogLvError         = 0x4
	LogLvFatal         = 0x5
	LogLvMessageToDing = 0x6
)

type dingModel struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

var lvMap = map[int]string{
	LogLvTrace:         "[Trace]  ",
	LogLvInfo:          "[Info]   ",
	LogLvWarning:       "[Warning]",
	LogLvError:         "[Error]  ",
	LogLvFatal:         "[Fatal]  ",
	LogLvMessageToDing: "[Ding]   ",
}
var dingURL = ""
var logLevel = LogLvInfo
var logFile *os.File
var logFileLock = &sync.Mutex{}
var txtMsgQue = make(chan *string, 1024)

type EZLoggerModel struct {
	LogLevel    int    `json:"log_level"`
	LogFilePath string `json:"log_file_path"`
	DingDingUrl string `json:"ding_ding_url"`
	DebugMode   bool   `json:"debug_mode"`
}

func SetUpEnv(m *EZLoggerModel) error {
	dingURL = m.DingDingUrl
	if m.LogLevel <= LogLvMessageToDing && dingURL == "" {
		return errors.New("ezlog level include ding talk but ding url is empty")
	}
	logLevel = m.LogLevel
	if !m.DebugMode {
		if err := ezFile.CreateDir(m.LogFilePath); err != nil {
			return err
		}
		go newLogFile(m.LogFilePath)
		go startLogServer()
	} else {
		go startDebugLogServer()
	}
	return nil
}
func T(msg ...interface{}) {
	ezlog(LogLvTrace, fmt.Sprintln(msg...))
}
func I(msg ...interface{}) {
	ezlog(LogLvInfo, fmt.Sprintln(msg...))
}
func W(msg ...interface{}) {
	ezlog(LogLvWarning, fmt.Sprintln(msg...))
}
func E(msg ...interface{}) {
	ezlog(LogLvError, fmt.Sprintln(msg...))
}
func F(msg ...interface{}) {
	ezlog(LogLvFatal, fmt.Sprintln(msg...))
}
func SendMessageToDing(msg ...interface{}) {
	go sendToDing(msgFmt(LogLvMessageToDing, fmt.Sprintln(msg...)))
	ezlog(LogLvMessageToDing, fmt.Sprintln(msg...))
}

func TWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvTrace, tag, fmt.Sprintln(msg...))
}
func IWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvInfo, tag, fmt.Sprintln(msg...))
}
func WWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvWarning, tag, fmt.Sprintln(msg...))
}
func EWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvError, tag, fmt.Sprintln(msg...))
}
func FWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvFatal, tag, fmt.Sprintln(msg...))
}
func SendMessageToDingWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvMessageToDing, tag, fmt.Sprintln(msg...))
	go sendToDing(msgWithTagFmt(LogLvMessageToDing, fmt.Sprintln(msg...), tag))
}
func ezlog(level int, msg string) {
	if level >= logLevel {
		rs := msgFmt(level, msg)
		txtMsgQue <- &rs
	}
}
func ezlogWithTag(level int, tag string, msg string) {
	if level >= logLevel {
		rs := msgWithTagFmt(level, tag, msg)
		txtMsgQue <- &rs
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
func startDebugLogServer() {
	runtime.LockOSThread()
	for {
		msg := <-txtMsgQue
		println(*msg)
	}
}
func startLogServer() {
	runtime.LockOSThread()
	for {
		for i := 0; i < 100; i++ {
			msg := <-txtMsgQue
			logFileLock.Lock()
			_, _ = logFile.WriteString(*msg)
			logFileLock.Unlock()
		}
		logFileLock.Lock()
		if err := logFile.Sync(); err != nil {
			println(err.Error())
		}
		logFileLock.Unlock()
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
func sendToDing(msg string) {
	m := dingModel{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: msg},
	}
	tmpData, err := json.Marshal(m)
	if err != nil {
		F(err.Error())
		return
	}
	rsp, err := http.Post(dingURL, "application/json;charset=utf-8", bytes.NewBuffer(tmpData))
	if err != nil {
		F(err.Error())
		return
	}
	rspInfo, _ := ioutil.ReadAll(rsp.Body)
	rspStr := string(rspInfo)
	if !strings.Contains(rspStr, `"errcode":0,`) {
		F("rsp from ding:" + rspStr)
	}
	_ = rsp.Body.Close()
}
