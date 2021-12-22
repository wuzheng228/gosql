package test1OpenConn

import "testing"

func TestInitDatabase(t *testing.T) {
	err := initDB()
	if err != nil {
		t.Fatalf("err:%v", err)
	}
}
