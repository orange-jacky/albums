package data

type ImageInfo struct {
	User       string    //用户名
	Album      string    //相册名
	Filepath   string    `bson:"-" json:"-"` //文件路径, 包含文件名
	Filename   string    //图片名
	Type       string    //图片扩展名
	Updatetime int64     //图片写入时间,精确到毫秒
	Md5        string    `json:"-"`                 //图片的md5
	Url        string    `bson:"-"`                 //图片给访问路径
	Features   []float64 `bson:"features" json:"-"` //图片特征向量
}

type ImageInfos []*ImageInfo
