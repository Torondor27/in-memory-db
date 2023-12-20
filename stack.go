package in_memory_db

import (
	"errors"
)

var errEmptyStack = errors.New("the stack is empty")

type operation struct {
	key      string
	newValue any
	oldValue any
	action   action
}

type operationStack struct {
	storage []operation
}

type transactionStack struct {
	storage []*operationStack
}

func (s *transactionStack) Push(v *operationStack) {
	s.storage = append(s.storage, v)
}

func (s *transactionStack) Pop() (*operationStack, error) {
	l := len(s.storage)
	if l == 0 {
		return nil, errEmptyStack
	}

	res := s.storage[l-1]
	s.storage = s.storage[:l-1]

	return res, nil
}

func (s *operationStack) Push(v operation) {
	s.storage = append(s.storage, v)
}

func (s *operationStack) Pop() (operation, error) {
	l := len(s.storage)
	if l == 0 {
		return operation{}, errEmptyStack
	}

	res := s.storage[l-1]
	s.storage = s.storage[:l-1]

	return res, nil
}
