package test2CRUD

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	err := initDB()
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
}

func TestQueryRowDemo(t *testing.T) {
	queryRowDemo()
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

func TestPrepareQueryDemo(t *testing.T) {
	prepareQueryDemo()
}

func TestPrepareInsertDemo(t *testing.T) {
	prepareInsertDemo()
}

func TestTransactionDemo(t *testing.T) {
	transactionDemo()
}

