package data

//图片属性
type Metadata struct {
	User       string //用户名
	Album      string //相册名
	Name       string //图片名
	Type       string //图片扩展名
	Updatetime int64  //图片写入时间,精确到毫秒

}

//查找图片的属性
type Attr struct {
	Id          string
	Md5         string
	Filename    string
	ContentType string
}

// 图片数据
type Imagedata struct {
	Metadata
	Attr
}

//图片数组
type Images []*Imagedata

// 图片特征数据
type Featuredata struct {
	Metadata
	Features []float64 //图片特征向量
	Attr
}

// 图片特征数组
type Features []*Featuredata
