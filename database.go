package in_memory_db

import (
	"errors"
)

type action string

const (
	actionCreate action = "create"
	actionDelete action = "delete"
	actionUpdate action = "update"
)

type Database struct {
	storage      map[string]any
	transactions *transactionStack
}

func New() *Database {
	return &Database{
		storage:      make(map[string]any),
		transactions: new(transactionStack),
	}
}

func (db *Database) Get(key string) any {
	return db.storage[key]
}

func (db *Database) Set(key string, value any) {
	db.logSetOperation(key, value)
	db.storage[key] = value
}

func (db *Database) Delete(key string) {
	db.logDeleteOperation(key)
	delete(db.storage, key)
}

func (db *Database) BeginTransaction() {
	db.transactions.Push(new(operationStack))
}

func (db *Database) Commit() {
	_, _ = db.transactions.Pop()
}

func (db *Database) Rollback() {
	t, err := db.transactions.Pop()
	if errors.Is(err, errEmptyStack) {
		return
	}

	for err == nil {
		var op operation
		op, err = t.Pop()
		db.rollbackOperation(op)
	}
}

func (db *Database) rollbackOperation(op operation) {
	switch op.action {
	case actionCreate:
		delete(db.storage, op.key)
	case actionUpdate:
		fallthrough
	case actionDelete:
		db.storage[op.key] = op.oldValue
	}
}

func (db *Database) logSetOperation(key string, value any) {
	db.logOperation(func() (op operation, ok bool) {
		act := actionUpdate
		oldValue, exists := db.storage[key]
		if !exists {
			act = actionCreate
		}

		op = operation{key: key, newValue: value, oldValue: oldValue, action: act}
		return op, true
	})
}

func (db *Database) logDeleteOperation(key string) {
	db.logOperation(func() (op operation, ok bool) {
		oldValue, exists := db.storage[key]
		if !exists {
			return
		}

		op = operation{key: key, newValue: nil, oldValue: oldValue, action: actionDelete}
		return op, true
	})
}

func (db *Database) logOperation(getOperation func() (op operation, ok bool)) {
	t, err := db.transactions.Pop()
	if errors.Is(err, errEmptyStack) {
		return
	}

	if op, ok := getOperation(); ok {
		t.Push(op)
	}

	db.transactions.Push(t)
}
