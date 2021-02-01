package ezSIPLikeMessage

import (
	"bytes"
	"errors"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"unicode/utf8"
)

type Header struct {
	Str1 string
	Str2 string
	Str3 string
}

var instance *Parser
var once sync.Once

type typeModel struct {
	headerIdx       int
	contentStrIdx   int
	otherInfoDicIdx int
	contentLenIdx   int
	nameDic         map[string]int
}
type Parser struct {
	once          *sync.Once
	innerDic      map[string]*typeModel
	innerDicLock  *sync.RWMutex
	gbCharDecoder *encoding.Decoder
}

func DefaultParser() *Parser {
	once.Do(func() {
		instance = &Parser{
			&sync.Once{},
			make(map[string]*typeModel, 4),
			&sync.RWMutex{},
			nil,
		}
		e, err := ianaindex.MIME.Encoding("gb18030")
		if err != nil {
			panic(err.Error())
		}
		instance.gbCharDecoder = e.NewDecoder()
	})
	return instance
}
func (p *Parser) InitParser(msg interface{}) error {
	rs := &typeModel{
		-1,
		-1,
		-1,
		-1,
		map[string]int{},
	}
	structType := reflect.TypeOf(msg).Elem()
	fieldsCount := structType.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := structType.Field(i)
		if field.Type == reflect.TypeOf(&Header{}) {
			rs.headerIdx = i
			continue
		}
		if sipFieldName := field.Tag.Get("sip"); sipFieldName != "" {
			if field.Type.Name() == "string" {
				rs.nameDic[sipFieldName] = i
				if sipFieldName == "Content" {
					rs.contentStrIdx = i
				}
				if sipFieldName == "Content-Length" {
					rs.contentLenIdx = i
				}
				continue
			}
			if field.Type.String() == "map[string]string" {
				rs.otherInfoDicIdx = i
			}
		}
	}
	if rs.headerIdx < 0 {
		return errors.New("u need to specify a field which has type '*ezSIPLikeMessage.Header' to hold header")
	}
	if rs.contentStrIdx < 0 {
		return errors.New("u need to specify a string field which has tag 'Content' to hold content")
	}
	if rs.otherInfoDicIdx < 0 {
		return errors.New("u need to specify a container field which has type 'map[string]string' to hold unknown fields")
	}
	if rs.contentLenIdx < 0 {
		return errors.New("u need to specify a string field which has tag 'Content-Length' to hold content")
	}
	p.innerDicLock.Lock()
	p.innerDic[reflect.TypeOf(msg).String()] = rs
	p.innerDicLock.Unlock()
	return nil
}
func Marshal(msg interface{}) ([]byte, error) {

}
func Unmarshal(data []byte, msg interface{}) (err error) {
	instance.innerDicLock.RLock()
	tm, valid := instance.innerDic[reflect.TypeOf(msg).String()]
	instance.innerDicLock.RUnlock()
	if !valid {
		err = errors.New("u must call InitParser before u use unmarshal")
		return err
	}
	var content string
	dataLen := len(data)
	idx := 0
	headerArr := [3]string{}
	arrIdx := 0
	otherInfo := map[string]string{}
	for i := 1; i < dataLen; i++ {
		if data[i] == ' ' {
			headerArr[arrIdx] = string(data[idx:i])
			arrIdx++
			idx = i + 1
		} else if data[i] == '\n' {
			if data[i-1] == '\r' {
				i -= 1
			}
			headerArr[arrIdx] = string(data[idx:i])
			idx = i + 2
			break
		}
	}
	if idx == 0 {
		err = errors.New("unknown msg")
		return err
	}
	refValue := reflect.ValueOf(msg).Elem()
	structFieldValue := refValue.Field(tm.headerIdx)
	val := reflect.ValueOf(&Header{
		Str1: headerArr[0],
		Str2: headerArr[1],
		Str3: headerArr[2],
	})
	structFieldValue.Set(val)
	colonIdx := idx - 1
	lineBreakerIdx := idx
	for i := idx; i < dataLen; i++ {
		if data[i] == '\n' {
			//new line
			if i-lineBreakerIdx < 2 {
				//这说明下一行行是个空行
				if dataLen-i < 2 {
					break
				}
				//能走到这里.显然剩下的就是contentData
				//contentData不需要包含第一行指定xml格式的信息
				if data[i+1] == '<' &&
					data[i+2] == '?' &&
					data[i+3] == 'x' &&
					data[i+4] == 'm' {
					//从这里开始就一定不合适
					for j := i + 1; j < dataLen-2; j++ {
						if data[j] == '\n' {
							i = j
							break
						}
					}
				}
				contentData := data[i+1 : dataLen-2]
				if !utf8.Valid(contentData) {
					r := transform.NewReader(bytes.NewBuffer(contentData), instance.gbCharDecoder)
					contentData, err = ioutil.ReadAll(r)
					if err != nil {
						return err
					}
					content = string(contentData)
					break
				}
				content = string(contentData)
				break
			}
			//走到这里说明不是终结.
			//lineBreakerIdx -> colonIdx 是key
			//colonIdx -> i 是value
			if lineBreakerIdx >= colonIdx {
				err = errors.New("unknown msg")
				return err
			}
			key := string(data[lineBreakerIdx:colonIdx])
			value := strings.TrimSpace(string(data[colonIdx+1 : i]))
			if index, valid := tm.nameDic[key]; !valid {
				otherInfo[key] = value
			} else {
				structFieldValue = refValue.Field(index)
				val = reflect.ValueOf(value)
				if !structFieldValue.IsValid() || !structFieldValue.CanSet() {
					otherInfo[key] = value
				} else {
					structFieldValue.Set(val)
				}
			}
			lineBreakerIdx = i + 1
		} else if data[i] == ':' {
			if colonIdx < lineBreakerIdx {
				colonIdx = i
			}
		}
	}
	structFieldValue = refValue.Field(tm.contentStrIdx)
	val = reflect.ValueOf(content)
	structFieldValue.Set(val)

	structFieldValue = refValue.Field(tm.otherInfoDicIdx)
	val = reflect.ValueOf(otherInfo)
	structFieldValue.Set(val)
	return err
}

//func Marshal(msg SIPLikeMessage) ([]byte,error){
//	d := msg.GetSIPParaNameDic()
//	rs := msg.GetHeader().Str1 + " " + msg.GetHeader().Str2 + " " + msg.GetHeader().Str3
//	structValue := reflect.ValueOf(msg).Elem()
//	if len(msg.ContentStr) > 0 {
//		msg.ContentLength = strconv.Itoa(len(msg.ContentStr))
//	}
//	for rtspName, index := range d {
//		value := structValue.Field(index)
//		strValue := value.String()
//		if strValue != "" {
//			rs += crlf + rtspName + ": " + strValue
//		}
//	}
//	if msg.KVPairs != nil {
//		for k, v := range msg.KVPairs {
//			rs += crlf + k + ": " + v
//		}
//	}
//	rs += crlf + crlf
//	if len(msg.ContentStr) > 0 {
//		rs += msg.ContentStr
//	}
//	return rs
//}
//func (msg SIPLikeMessage) ToString() string {
//}
//func (msg *Message) ToBytes() []byte {
//	rs := []byte(msg.ToString())
//	return rs
//}
