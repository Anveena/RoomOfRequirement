package ezSyncAreaCode

import (
	"encoding/json"
	"errors"
	"github.com/Anveena/RoomOfRequirement/ezMySQL"
	"github.com/Anveena/RoomOfRequirement/ezPasswordEncoder"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
	"xorm.io/xorm"
)

type sortHelper []*areaDBModel

func (s sortHelper) Len() int {
	return len(s)
}
func (s sortHelper) Less(i, j int) bool {
	iCode, _ := strconv.Atoi(s[i].Code)
	jCode, _ := strconv.Atoi(s[j].Code)
	return iCode < jCode
}
func (s sortHelper) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func SyncFromAMap(provinceName string, aMapKey string, dbInfo *ezMySQL.Info) error {
	eng, err := initDBEnv(dbInfo)
	if err != nil {
		return err
	}
	defer func() {
		_ = eng.Close()
	}()
	toWriteArr, err := getAMapData(provinceName, aMapKey)
	if err != nil {
		return err
	}
	sort.Sort(toWriteArr)
	_, err = eng.Insert(toWriteArr)
	if err != nil {
		return err
	}
	return nil
}
func getAMapData(provinceName string, aMapKey string) (sortHelper, error) {
	aMapURL := "https://restapi.amap.com/v3/config/district?key=" + aMapKey + "&subdistrict=3&keywords=" + provinceName
	println(aMapURL)
	rsp, err := http.Get(aMapURL)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, errors.New("http status is not 200")
	}
	bodyData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	var rspModel aMapRspModel
	if err = json.Unmarshal(bodyData, &rspModel); err != nil {
		return nil, err
	}
	var toWriteArea sortHelper
	for _, p := range rspModel.Districts {
		//省级
		pCode := p.AdCode + "00"
		toWriteArea = append(toWriteArea, &areaDBModel{
			Code:       pCode,
			Name:       p.Name,
			Center:     p.Center,
			ParentCode: "00000000",
			ParentName: "中国",
			CreateTime: time.Time{},
		})
		if p.Level != "province" {
			return nil, errors.New("在province级里面发现了" + p.Name + p.Level)
		}
		for _, c := range p.Districts {
			//市级
			cCode := c.AdCode + "00"
			toWriteArea = append(toWriteArea, &areaDBModel{
				Code:       cCode,
				Name:       c.Name,
				Center:     c.Center,
				ParentCode: pCode,
				ParentName: p.Name,
				CreateTime: time.Time{},
			})
			cSubStreet := -1
			for j, d := range c.Districts {
				//区级
				if d.Level == "street" {
					if cSubStreet == 1 {
						//https://lbs.amap.com/api/webservice/guide/api/district:
						//目前部分城市和省直辖县因为没有区县的概念，故在市级下方直接显示街道。例如：广东-东莞、海南-文昌市。
						return nil, errors.New("both street and district under " + c.Name)
					} else if cSubStreet == -1 {
						//这里还没有进来过 新弄一个区
						toWriteArea = append(toWriteArea, &areaDBModel{
							Code:       c.AdCode[:len(c.AdCode)-2] + "0100",
							Name:       c.Name + "区",
							Center:     c.Center,
							ParentCode: cCode,
							ParentName: c.Name,
							CreateTime: time.Time{},
						})
					}
					cSubStreet = 0
					dCode := c.AdCode[:len(c.AdCode)-2] + "01"
					if j < 9 {
						dCode = dCode + "0" + strconv.Itoa(j+1)
					} else {
						dCode = dCode + strconv.Itoa(j+1)
					}
					toWriteArea = append(toWriteArea, &areaDBModel{
						Code:       dCode,
						Name:       d.Name,
						Center:     d.Center,
						ParentCode: c.AdCode[:len(c.AdCode)-2] + "0100",
						ParentName: c.Name + "区",
						CreateTime: time.Time{},
					})
				} else {
					if cSubStreet == 0 {
						//https://lbs.amap.com/api/webservice/guide/api/district:
						//目前部分城市和省直辖县因为没有区县的概念，故在市级下方直接显示街道。例如：广东-东莞、海南-文昌市。
						return nil, errors.New("both street and district under " + c.Name)
					}
					cSubStreet = 1
					dCode := d.AdCode + "00"
					toWriteArea = append(toWriteArea, &areaDBModel{
						Code:       dCode,
						Name:       d.Name,
						Center:     d.Center,
						ParentCode: cCode,
						ParentName: c.Name,
						CreateTime: time.Time{},
					})
					for i, s := range d.Districts {
						//街道级
						sCode := s.AdCode
						if i < 9 {
							sCode = sCode + "0" + strconv.Itoa(i+1)
						} else {
							sCode = sCode + strconv.Itoa(i+1)
						}
						toWriteArea = append(toWriteArea, &areaDBModel{
							Code:       sCode,
							Name:       s.Name,
							Center:     s.Center,
							ParentCode: dCode,
							ParentName: d.Name,
							CreateTime: time.Time{},
						})
					}
				}
			}
		}
	}
	return toWriteArea, nil
}
func initDBEnv(dbInfo *ezMySQL.Info) (*xorm.Engine, error) {
	if dbInfo.PasswordBase64Str == "" {
		return nil, errors.New("empty password base64 str")
	}
	password, err := ezPasswordEncoder.GetPasswordFromEncodedStr(dbInfo.PasswordBase64Str)
	if err != nil {
		return nil, err
	}

	eng, err := xorm.NewEngine("mysql", dbInfo.Account+
		":"+password+"@tcp("+dbInfo.Host+":"+dbInfo.Port+")/"+dbInfo.Name+"?charset=utf8")
	if err != nil {
		return nil, errors.New("db engine init failed ,err:" + err.Error())
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	eng.SetTZLocation(loc)
	eng.SetTZDatabase(loc)
	err = eng.Sync2(&areaDBModel{})
	if err != nil {
		return nil, err
	}
	return eng, nil
}
