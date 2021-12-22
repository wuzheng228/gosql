package test1OpenConn

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB // 表示连接的数据库对象， 内部维护着一个连接池，可以安全的被多个goroutine使用

func initDB() (err error) {
	dsn := "root:123@tcp(127.0.0.1:3306)/tes22t"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("打开失败，db:%#v", db)
		return err
	}
	// 尝试与数据库建立连接
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func Open() {
	// DSN:data source name
	dsn := "root:123@tcp(127.0.0.1:3306)/t1231est"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Printf("打开失败，db:%#v", db)
		panic(err)
	}
	defer db.Close() // 写在error判断的下面
	fmt.Printf("成功打开数据库，db:%#v", db)
}
