产生go接口代码
thrift --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift ImageHandle.thrift

生成gen-go/ImageHandle目录,里面有
	3个文件
	constants.go
	handler.go
	ttypes.go

	1一个目录
	handler-remote   //目录有一个handler-remote.go文件,是go做为服务端的代码,只有入口逻辑,业务逻辑还得单独写

	把3个文件拷贝到imagehandle目录下,覆盖同名文件
	然后再编写 go做为客户端的代码,修改imagehandle_client.go


产生cpp接口代码
thrift --gen cpp ImageHandle.thrift
