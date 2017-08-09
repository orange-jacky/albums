package data

///////////////////////////////////////////////////////////////////////////////////
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
	Filename    string //图片访问地址
	ContentType string
}

///////////////////////////////////////////////////////////////////////////////////
// 图片数据,单独保存图片
type Imagedata struct {
	Metadata
	Attr
}

//图片数组
type Images []*Imagedata

///////////////////////////////////////////////////////////////////////////////////
// 图片特征数据
type Featuredata struct {
	Metadata `bson:",inline"`
	Features []float64 `bson:"features" json:"-"` //图片特征向量
	Attr     `bson:",inline"`
}

// 图片特征数组
type Features []*Featuredata

///////////////////////////////////////////////////////////////////////////////////
//保存用户与图片url关联属性
type UserImage struct {
	Metadata `bson:",inline"`
	Attr     `bson:",inline"`
}

type UserImages []*UserImage

///////////////////////////////////////////////////////////////////////////////////
//用户collection
type User struct {
	UserName string
	Password string
}
