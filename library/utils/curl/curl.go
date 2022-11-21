package curl

import (
	"github.com/idoubi/goz"
)

type RetData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

//curl  Get请求
//Query: map[string]interface{}{
//			"key1": "value1",
//			"key2": []string{"value21", "value22"},
//			"key3": "333",
//		},
func Get(url string, query map[string]interface{}, header map[string]interface{}) (ret []byte, err error) {
	resp, err := goz.NewClient().Get(url, goz.Options{Query: query, Headers: header})
	if err != nil {
		return
	}

	body, err := resp.GetBody()
	if err != nil {
		return
	}

	ret = body
	return
}

// post请求
// FormParams: map[string]interface{}{
//        "key1": "value1",
//        "key2": []string{"value21", "value22"},
//        "key3": "333",
//    },
func Post(url string, data interface{}) (ret []byte, err error) {
	resp, err := goz.NewClient().Post(url, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: data,
	})
	if err != nil {
		return
	}

	body, err := resp.GetBody()
	if err != nil {
		return
	}

	ret = body
	return
}
