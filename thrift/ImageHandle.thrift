#生成的类名
namespace go  imagehandle   //go的package 名称小写
namespace cpp ImageHandle   //c++的类名大写

/**
 * Result 保存深度学习/深度学习做物体检测返回结果
 * desc 保存返回描述符 
 * ret  保存处理过后图片
 */
struct Result {
	1: string desc,
	2: binary ret,
}

#接口服务
service Handler{
	#hsv 提取特征值
  	list<double> Feature(1:binary image),
  	#深度学习
  	Result DeepLearning(1:binary image),
	#深度学习做物体检测
  	Result ObjectDetectionDL(1:binary image)
}