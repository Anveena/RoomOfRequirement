package main

import (
	"fmt"
	"github.com/Anveena/RoomOfRequirement/ezConfig"
	"github.com/Anveena/RoomOfRequirement/ezLog"
	"github.com/Anveena/RoomOfRequirement/ezPasswordEncoder"
	"github.com/Anveena/RoomOfRequirement/ezRandom"
	"github.com/Anveena/RoomOfRequirement/ezUTCTime"
	"strings"
	"sync"
	"time"
)

func main() {
	ezUTCTime.SyncTimeFromAliyun()
	println(ezUTCTime.GetAliyunTime().String())

	ta := time.Now()
	println(fmt.Sprintf("logs_of_%v_%v_%v", ta.Year(), strings.ToLower(ta.Month().String()), ta.Day()))
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
	t := time.Now()
	for i := 0; i < 2000; i++ {
		ezLog.E("鸟宿池边树", i)
		ezLog.DWithTag("测试", "僧推月下门", i)
	}
	println(time.Now().Sub(t).String())
	//8638000 5444000 5748000
	println(time.Now().Sub(t).String())
	ezLog.D(ezRandom.RandomString(ezRandom.OctNumberOnly, 16), ezRandom.OctNumberOnly)
	ezLog.E(ezRandom.RandomString(ezRandom.CapitalLetterOnly, 16), ezRandom.CapitalLetterOnly)
	ezLog.I(ezRandom.RandomString(ezRandom.NumberAndCapitalLetter, 16), ezRandom.NumberAndCapitalLetter)
	//ezLog.DingMessage("理论上这是不@的\n鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	ezLog.DingList("理论上这是@我的\n鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	//ezLog.DingAtAll("理论上这是@所有人的\n鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	//ezLog.DingMessageWithTag("test tag","理论上这是不@的\n鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	ezLog.DingListWithTag("test tag", "理论上这是@我的\n鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	//ezLog.DingAtAllWithTag("test tag","理论上这是@所有人的\n鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	//ezLog.SendMessageToDing(ezRandom.RandomString(ezRandom.NumberAndLowercaseLetter, 16), ezRandom.NumberAndLowercaseLetter)
	ezLog.I(ezPasswordEncoder.EncodePassword("jyydb_2015!"))
	//f, err := ezFile.CreateFile("/Users/panys/Desktop/", "wzz.txt", true, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	//if err != nil {
	//	ezLog.E(err.Error())
	//	return
	//}
	//_, err = f.Write([]byte("宇智波"))
	//if err != nil {
	//	ezLog.E(err.Error())
	//	return
	//}
	//f.Truncate(0)
	//_, err = f.Write([]byte("多由也duoyouye"))
	//if err != nil {
	//	ezLog.E(err.Error())
	//	return
	//}
	//if err = f.Close(); err != nil {
	//	ezLog.E(err.Error())
	//	return
	//}
	//orikey := "大番薯"
	//aesKey := ezCrypto.MakeMD5Key(orikey, time.Now().UnixNano()%251)
	//data := []byte("鹅鹅鹅\n曲项向天歌\n白毛浮绿水\n红掌拨清波")
	//encData, err := ezCrypto.AESCBCEncrypt(&data, &aesKey)
	//if err != nil {
	//	ezLog.E(err.Error())
	//	return
	//}
	//oriData, err := ezCrypto.AESCBCDecrypt(encData, &aesKey)
	//if err != nil {
	//	ezLog.E(err.Error())
	//	return
	//}
	//ezLog.I("aes decrypt result:", string(*oriData))
	//root := ezXMLTreeMaker.NewXMLTree("root", "")
	//sub := ezXMLTreeMaker.NewXMLTree("sub", "大番薯")
	//sub.SetAttr("身高", "165")
	//sub.SetAttr("体重", "65")
	//root.AddNode(sub)
	//sub.SetValue("秦先生")
	//sub.SetAttr("体重", "85")
	//root.AddNode(sub)
	//ezLog.I(root.StrValue())
	//para := ezHash.Parameters{
	//	Width:      32,
	//	Polynomial: 0x04C11DB7,
	//	ReflectIn:  false,
	//	ReflectOut: false,
	//	Init:       0xFFFFFFFF,
	//	FinalXor:   0,
	//}
	//rs := ezHash.CalculateCRC(&para, []byte{0x00, 0xB0, 0x0D, 0x00, 0x01, 0xC1, 0x00, 0x00, 0x00, 0x01, 0xF0, 0x00})
	//ezLog.I("crc value", fmt.Sprintf("%x\n", rs))
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
