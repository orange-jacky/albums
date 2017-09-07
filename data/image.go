package data

type Image struct {
	Filepath string `bson:"-"` //文件路径, 包含文件名
	Md5      string //文件计算md5后赋值给这个字段
}

type Images []*Image
