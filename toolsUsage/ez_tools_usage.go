package main

import (
	"fmt"
	"github.com/Anveena/RoomOfRequirement/ezConfig"
	"github.com/Anveena/RoomOfRequirement/ezCrypto"
	"github.com/Anveena/RoomOfRequirement/ezFile"
	"github.com/Anveena/RoomOfRequirement/ezHash"
	"github.com/Anveena/RoomOfRequirement/ezLog"
	"github.com/Anveena/RoomOfRequirement/ezMySQL"
	"github.com/Anveena/RoomOfRequirement/ezRandom"
	"github.com/Anveena/RoomOfRequirement/ezXMLTreeMaker"
	"os"
	"time"
)

func main() {
	var obj ezLog.EZLoggerModel
	err := ezConfig.ReadConf(&obj)
	if err != nil {
		println(err.Error())
		return
	}
	err = ezLog.SetUpEnv(&obj)
	if err != nil {
		println(err.Error())
		return
	}
	ezLog.T(ezRandom.RandomString(ezRandom.OctNumberOnly, 16), ezRandom.OctNumberOnly)
	ezLog.I(ezRandom.RandomString(ezRandom.HexNumberOnly, 16), ezRandom.HexNumberOnly)
	ezLog.W(ezRandom.RandomString(ezRandom.LowercaseLetterOnly, 16), ezRandom.LowercaseLetterOnly)
	ezLog.E(ezRandom.RandomString(ezRandom.CapitalLetterOnly, 16), ezRandom.CapitalLetterOnly)
	ezLog.F(ezRandom.RandomString(ezRandom.NumberAndCapitalLetter, 16), ezRandom.NumberAndCapitalLetter)
	ezLog.SendMessageToDing(ezRandom.RandomString(ezRandom.NumberAndLowercaseLetter, 16), ezRandom.NumberAndLowercaseLetter)
	ezLog.I(ezMySQL.MakePasswordBase64Str("jyydb_2015!"))
	f, err := ezFile.CreateFile("/Users/panys/Desktop/", "wzz.txt", true, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		ezLog.F(err.Error())
		return
	}
	_, err = f.Write([]byte("宇智波"))
	if err != nil {
		ezLog.F(err.Error())
		return
	}
	f.Truncate(0)
	_, err = f.Write([]byte("多由也duoyouye"))
	if err != nil {
		ezLog.F(err.Error())
		return
	}
	if err = f.Close(); err != nil {
		ezLog.F(err.Error())
		return
	}
	orikey := "大番薯"
	aesKey := ezCrypto.MakeMD5Key(orikey, time.Now().UnixNano()%251)
	data := []byte("鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	encData, err := ezCrypto.AESCBCEncrypt(&data, &aesKey)
	if err != nil {
		ezLog.F(err.Error())
		return
	}
	oriData, err := ezCrypto.AESCBCDecrypt(encData, &aesKey)
	if err != nil {
		ezLog.F(err.Error())
		return
	}
	ezLog.I("aes decrypt result:", string(*oriData))
	root := ezXMLTreeMaker.NewXMLTree("root", "")
	sub := ezXMLTreeMaker.NewXMLTree("sub", "大番薯")
	sub.SetAttr("身高", "165")
	sub.SetAttr("体重", "65")
	root.AddNode(sub)
	sub.SetValue("秦先生")
	sub.SetAttr("体重", "85")
	root.AddNode(sub)
	ezLog.I(root.StrValue())
	para := ezHash.Parameters{
		Width:      32,
		Polynomial: 0x04C11DB7,
		ReflectIn:  false,
		ReflectOut: false,
		Init:       0xFFFFFFFF,
		FinalXor:   0,
	}
	rs := ezHash.CalculateCRC(&para, []byte{0x00, 0xB0, 0x0D, 0x00, 0x01, 0xC1, 0x00, 0x00, 0x00, 0x01, 0xF0, 0x00})
	ezLog.I("crc value", fmt.Sprintf("%x\n", rs))
}
