package in_memory_db

import (
	"testing"
)

func Test_Example1(t *testing.T) {
	db := New()
	db.Set("key1", "value1")
	db.BeginTransaction()
	db.Set("key1", "value2")
	db.Commit()
	assertGet(t, db.Get("key1"), "value2")
}

func Test_Example2(t *testing.T) {
	db := New()
	db.Set("key1", "value1")
	db.BeginTransaction()
	assertGet(t, db.Get("key1"), "value1")
	db.Set("key1", "value2")
	assertGet(t, db.Get("key1"), "value2")
	db.Rollback()
	assertGet(t, db.Get("key1"), "value1")
}

func Test_Example3(t *testing.T) {
	db := New()
	db.Set("key1", "value1")
	db.BeginTransaction()
	db.Set("key1", "value2")
	assertGet(t, db.Get("key1"), "value2")
	db.BeginTransaction()
	assertGet(t, db.Get("key1"), "value2")
	db.Delete("key1")
	db.Commit()
	assertGet(t, db.Get("key1"), nil)
	db.Commit()
	assertGet(t, db.Get("key1"), nil)
}

func Test_Example4(t *testing.T) {
	db := New()
	db.Set("key1", "value1")
	db.BeginTransaction()
	db.Set("key1", "value2")
	assertGet(t, db.Get("key1"), "value2")
	db.BeginTransaction()
	assertGet(t, db.Get("key1"), "value2")
	db.Delete("key1")
	db.Rollback()
	assertGet(t, db.Get("key1"), "value2")
	db.Commit()
	assertGet(t, db.Get("key1"), "value2")
}

func Test_CustomExample(t *testing.T) {
	db := New()
	db.Set("key1", "test")

	//<<< transaction 1
	db.BeginTransaction()
	db.Set("key1", "test1")
	//>>> transaction 2
	db.BeginTransaction()
	db.Set("key1", "test2")
	//<<< transaction 3
	db.BeginTransaction()
	db.Delete("key1")
	assertGet(t, db.Get("key1"), nil)
	db.Rollback()
	//>>> transaction 3
	assertGet(t, db.Get("key1"), "test2")
	db.Rollback()
	//>>> transaction 2
	assertGet(t, db.Get("key1"), "test1")
	db.Rollback()
	//>>> transaction 1

	//<<< transaction 4
	db.BeginTransaction()
	db.Delete("key1")
	db.Commit()
	//>>> transaction 4

	assertGet(t, db.Get("key1"), nil)
}

func Test_CustomExample2(t *testing.T) {
	db := New()
	db.Set("key1", "test")

	db.BeginTransaction()
	db.Set("key1", "test1")
	db.Set("key1", "test2")
	db.Set("key1", "test3")
	db.Delete("key1")
	db.Set("key1", "test4")
	db.Set("key2", "lorem")
	db.Delete("key2")
	db.Set("key2", "ipsum")
	db.Rollback()

	assertGet(t, db.Get("key1"), "test")
	assertGet(t, db.Get("key2"), nil)
}

func assertGet(t *testing.T, got, want any) {
	if got != want {
		t.Errorf("Invalid value. want: %v, got: %v", want, got)
	}
}
