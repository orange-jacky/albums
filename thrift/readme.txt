产生go接口代码
thrift --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift feature.thrift

产生cpp接口代码
thrift --gen cpp feature.thrift
