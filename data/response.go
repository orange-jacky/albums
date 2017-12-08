package data

import (
	"encoding/json"
)

//正常search返回结果
type Response struct {
	Status            int         `json:"status"`
	Cost              int64       `json:"cost"`
	Total             int         `json:"total"`
	Data              interface{} `json:"data"` //保存返回数据
	StatusDescription interface{} `json:"statusdescription"`
}

func (item *Response) Unmarshal(s []byte) error {
	return json.Unmarshal(s, item)
}
