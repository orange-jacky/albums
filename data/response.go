package data

//正常search返回结果
type Response struct {
	Status int         `json:"status"`
	Cost   int         `json:"cost"`
	Data   interface{} `json:"data"` //保存返回数据
}
