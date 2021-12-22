package test3Sqlx

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type user struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func initDB() (err error) {
	dsn := "root:123@tcp(127.0.0.1:3306)/test"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

// 查询单条数据库示例
func queryRowDemo() {
	sqlStr := "select id, name from user where id=?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s\n", u.Id, u.Name)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

// 插入数据
func insertRowDemo() {
	sqlstr := "insert into user (id, name) values (?, ?)"
	ret, err := db.Exec(sqlstr, 666, "溜溜溜")
	if err != nil {
		fmt.Printf("insert failed, err:%v", err)
		return
	}
	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastInsertId failed, err: %v", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlstr := "update user set name=? where id = ?"
	ret, err := db.Exec(sqlstr, "六六六", 666)
	if err != nil {
		fmt.Printf("update user error, err: %v", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 666)
	if err != nil {
		fmt.Printf("delete err, err:%v", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

// DB.NamedExec方法用来绑定SQL语句与结构体或map中的同名字段。
func insertUserDemo() (err error) {
	sqlStr := "insert into user (id, name) values(:id, :name)"
	_, err = db.NamedExec(sqlStr, map[string]interface{}{
		"id":   777,
		"name": "七七七",
	})
	if err != nil {
		fmt.Printf("NamedExec insert failed, err : %v", err)
	}
	return
}

// NamedQuery
func namedQuery() {
	sqlStr := "select * from user where name=:name"
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "七七七"})
	if err != nil {
		fmt.Printf("namedQuery failed, err : %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan faild, err :%v", err)
			continue
		}
		fmt.Printf("user:%#v", u)
	}
	u := user{
		Name: "七七七",
	}
	// 使用结构体字段查询，根据结构体字段的db tag进行映射
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("namedQuery failed, err : %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan faild, err :%v", err)
			continue
		}
		fmt.Printf("user:%#v", u)
	}
}

// 事务操作
func transactionDemo2() (err error) {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()
	sqlStr1 := "update user set name = '八八八' where id = ?"
	rs, err := tx.Exec(sqlStr1, 777)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "update user set name='测试' where id= ?"
	rs, err = tx.Exec(sqlStr2, 1)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return err
}

// 自己拼接语句实现批量插入
func BatchInsertUsers(users []*user) error {
	// 存放(?, ?)的切片
	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users)*2)
	// 遍历users准备相关数据
	for _, u := range users {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.Id)
		valueArgs = append(valueArgs, u.Name)
	}
	stmt := fmt.Sprintf("INSERT INTO user (id, name) values %s", strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}

// 使用sqlx.In实现批量插入
// 前提是需要结构体实现driver.Valuer接口:
func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Id, u.Name}, nil
}

func BatchInsertUser2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (id, name) values (?), (?)", // (?)的数量要等于len(users)
		users...,
	)
	fmt.Println(query) //查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := db.Exec(query, args...)
	return err
}

// 使用NamedExec实现批量插入 该功能需1.3.1版本以上，并且1.3.1版本目前还有点问题，sql语句最后不能有空格和;
func BatchInsertUsers3(users []*user) error {
	_, err := db.NamedExec("INSERT INTO user (id, name) values(:id, :name)", users)
	return err
}

// sqlx.In 查询数据
func QueryByIds(ids []int) (users []user, err error) {
	// 动态填充id
	query, args, err := sqlx.In("SELECT id, name FROM user WHERE id IN (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&users, query, args...)
	return
}

func QueryAndOrderByIDs(ids []int) (users []user, err error) {
	// 动态填充id
	strIds := make([]string, 0, len(ids))
	for _, id := range ids {
		strIds = append(strIds, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("SELECT id, name FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIds, ","))
	if err != nil {
		return
	}

	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}
