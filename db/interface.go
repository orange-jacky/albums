package db

// Db 定义了一个对存储操作的接口
type Db interface {
	Connect(hosts, db string) error   //连接
	OpenDb(db string) error           //打开库
	OpenTable(table string)           //开打表
	Query(query interface{}) error    //查询
	Insert(docs ...interface{}) error //写入
	Close() error                     //关闭连接
}
