package ezNetworking

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//todo 重写
func PostHttpReq(url string, reqJsonData interface{}) (*http.Response, error) {
	var bf *bytes.Buffer
	if reqJsonData == nil {
		bf = nil
	} else {
		reqData, err := json.Marshal(reqJsonData)
		if err != nil {
			return nil, err
		}
		bf = bytes.NewBuffer(reqData)
	}
	client := &http.Client{}
	req, err := http.NewRequest("post", url, bf)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
