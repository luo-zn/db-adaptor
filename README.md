# db-adaptor
数据库适配器。此库用于对接业务model层与数据库驱动层。

## 快速使用
```go
package main

import "github.com/luo-zn/db-adaptor"

var DBAdapt *db_adaptor.DbAdaptor

func main(){
    opt := db_adaptor.AdaptorOptions{Uri: "mongo://localhost:27017",DBType:"mongodb"}
    DBAdapt = db_adaptor.NewDbAdaptor(&opt)
}
```


## [例子](docs/examples.md)