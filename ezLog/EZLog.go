package ezLog

/*
#if (defined _WIN32) || (defined __WINDOWS_)
#include <windows.h>
#endif
void initPrint() {
#if (defined _WIN32) || (defined __WINDOWS_)
    HANDLE handle = GetStdHandle(STD_OUTPUT_HANDLE);
    DWORD mode;
    GetConsoleMode(handle, &mode);
    mode |= (DWORD)0x4;
    SetConsoleMode(handle, mode);
#endif
}
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
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
var levelShouldPrintAtTerminal = LogLvInfo

type EZLoggerModel struct {
	WhichLevelLogShouldPrint int    `json:"which_level_log_should_print"`
	DingDingUrl              string `json:"ding_ding_url"`
}

func SetUpEnv(m *EZLoggerModel) error {
	C.initPrint()
	dingURL = m.DingDingUrl
	if m.WhichLevelLogShouldPrint <= LogLvMessageToDing && dingURL == "" {
		return errors.New("log level include ding talk but ding url is empty")
	}
	levelShouldPrintAtTerminal = m.WhichLevelLogShouldPrint
	return nil
}
func T(msg ...interface{}) {
	log(LogLvTrace, fmt.Sprintln(msg...))
}
func I(msg ...interface{}) {
	log(LogLvInfo, fmt.Sprintln(msg...))
}
func W(msg ...interface{}) {
	log(LogLvWarning, fmt.Sprintln(msg...))
}
func E(msg ...interface{}) {
	log(LogLvError, fmt.Sprintln(msg...))
}
func F(msg ...interface{}) {
	log(LogLvFatal, fmt.Sprintln(msg...))
}
func SendMessageToDing(msg ...interface{}) {
	go sendToDing(fmt.Sprintln(msg...))
	log(LogLvMessageToDing, fmt.Sprintln(msg...))
}
func log(level int, msg string) {
	if level >= levelShouldPrintAtTerminal {
		_, file, line, _ := runtime.Caller(2)
		strHeader := "\n" + lvMap[level] + " "
		msg = msg[:len(msg)-1]
		msgFormatted := strings.ReplaceAll(msg, "\n", strHeader)
		finalStr := strHeader + "File:" + file +
			strHeader + "Line:" + strconv.Itoa(line) +
			strHeader + "Time:" + time.Now().String() +
			strHeader + msgFormatted
		mPrint(level, finalStr)
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
func mPrint(level int, msg string) {
	switch level {
	case LogLvTrace:
		fmt.Printf("\n %c[1;30m%s%c[0m", 0x1B, msg, 0x1B)
		break
	case LogLvInfo:
		fmt.Printf("\n %c[1;32m%s%c[0m", 0x1B, msg, 0x1B)
		break
	case LogLvWarning:
		fmt.Printf("\n %c[1;33m%s%c[0m", 0x1B, msg, 0x1B)
		break
	case LogLvError:
		fmt.Printf("\n %c[1;31m%s%c[0m", 0x1B, msg, 0x1B)
		break
	case LogLvFatal:
		fmt.Printf("\n %c[0;34m%s%c[0m", 0x1B, msg, 0x1B)
		break
	case LogLvMessageToDing:
		fmt.Printf("\n %c[1;35m%s%c[0m", 0x1B, msg, 0x1B)
		break
	}
}
