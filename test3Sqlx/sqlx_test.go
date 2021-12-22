package test3Sqlx

import (
	"fmt"
	"testing"
)

func TestInitDB(t *testing.T) {
	initDB()
}

func TestQueryMultiRowDemo(t *testing.T) {
	queryMultiRowDemo()
}

func TestInsertRowDemo(t *testing.T) {
	insertRowDemo()
}

func TestUpdateRowDemo(t *testing.T) {
	updateRowDemo()
}

func TestDeleteRowDemo(t *testing.T) {
	deleteRowDemo()
}

func TestInsertUserDemo(t *testing.T) {
	insertUserDemo()
}

func TestNamedQuery(t *testing.T) {
	namedQuery()
}

func TestTransactionDemo2(T *testing.T) {
	transactionDemo2()
}

func TestBatchInsertUsers(t *testing.T) {
	var u1 *user = &user{Id: 999, Name: "九九九"}
	var u2 *user = &user{Id: 998, Name: "九九八"}
	users := []*user{u1, u2}
	err := BatchInsertUsers(users)
	if err != nil {
		t.Fatalf("bathInsert failed, err: %v", err)
	}
}

func TestBatchInsertUsers2(t *testing.T) {
	users := []interface{}{
		user{Id: 997, Name: "997"},
		user{Id: 996, Name: "996"},
	}
	err := BatchInsertUser2(users)
	if err != nil {
		t.Fatalf("bathInsertFaild, err:%v", err)
	}
}

func TestBatchInsertUsers3(t *testing.T) {
	var u1 *user = &user{Id: 995, Name: "九九五"}
	var u2 *user = &user{Id: 994, Name: "994"}
	users := []*user{u1, u2}
	err := BatchInsertUsers3(users)
	if err != nil {
		t.Fatalf("bathInsertFaild, err:%v", err)
	}
}

func TestQueryByIds(t *testing.T) {
	ids := []int{222, 3, 1, 124, 123}
	user, err := QueryByIds(ids)
	if err != nil {
		t.Fatalf("err: %v", err)
		return
	}
	fmt.Printf("%#v", user)
}

func TestQueryAndOrderByIDs(t *testing.T) {
	ids := []int{222, 3, 1, 124, 123}
	user, err := QueryAndOrderByIDs(ids)
	if err != nil {
		t.Fatalf("err: %v", err)
		return
	}
	fmt.Printf("%#v", user)
}
