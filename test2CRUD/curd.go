package test2CRUD

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB // 表示连接的数据库对象， 内部维护着一个连接池，可以安全的被多个goroutine使用

func initDB() (err error) {
	dsn := "root:123@tcp(127.0.0.1:3306)/test"
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

// 创建结构体封装数据库查询对象
type user struct {
	id   int
	name string
}

/**
* 查询单条数据
 */
func queryRowDemo() {
	sqlStr := "select id, name from user where id = ?"
	var u user
	// 确保QueryRow之后调用Scan方法，否则持有的数据库连接不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s\n", u.id, u.name)
}

/**
* 查询多条数据
 */
func queryMultiRowDemo() {
	sqlStr := "select id, name from user where id > 0"
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	// 非常重要，关闭rows释放所有连接数据
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name)
		if err != nil {
			fmt.Printf("Scan failed ,err: %v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s\n", u.id, u.name)
	}
}

/**
* 插入一条数据
 */
func insertRowDemo() {
	sqlstr := "insert into user (id, name) values (?, ?)"
	ret, err := db.Exec(sqlstr, 223, "王五")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

/**
* 更新数据
 */
func updateRowDemo() {
	sqlStr := "update user set name = ? where id = ?"
	ret, err := db.Exec(sqlStr, "张s", 222)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

/**
*	删除数据
 */

func deleteRowDemo() {
	sqlstr := "delete from user where id = ?"
	ret, err := db.Exec(sqlstr, 2)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

/**
*	预处理查询示例
 */
func prepareQueryDemo() {
	sqlstr := "select id, name from user where id > ?"
	s, err := db.Prepare(sqlstr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v", err)
		return
	}
	defer s.Close()
	rows, err := s.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s\n", u.id, u.name)
	}
}

/**
*	预处理插入示例
 */
func prepareInsertDemo() {
	sqlStr := "insert into user (id, name) values (?, ?)"
	s, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed err:%v", err)
		return
	}
	defer s.Close()
	_, err = s.Exec(123, "大王子")
	if err != nil {
		fmt.Printf("insert failed err:%v", err)
		return
	}
	_, err = s.Exec(124, "小公主")
	if err != nil {
		fmt.Printf("insert failed err:%v", err)
		return
	}

	fmt.Printf("insert succeed")
}

/**
* 事务处理示例
 */
func transactionDemo() {
	tx, err := db.Begin() // 开启一个事务
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Printf("begin trans failed , err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set name = ? where id = ?"
	ret1, err := tx.Exec(sqlStr1, "大公主", 124)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set name='王大锤' where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}
	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")

}
