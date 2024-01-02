package stack

import (
	"sync"
)

// Stack is a generic old-skool stack implementation in Go.
type Stack[T any] struct {
	lock sync.Mutex // ensures that the stack is concurrently safe
	s    []T
}

// New creates a new Stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds a value to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

// Pop removes the value at the top of the stack and returns it.
// If the stack is empty, it returns the zero value of the type T and a boolean false.
func (s *Stack[T]) Pop() (T, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]

	return res, true
}

// Peek returns the value at the top of the stack without removing it.
// If the stack is empty, it returns the zero value of the type T and a boolean false.
func (s *Stack[T]) Peek() (T, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	return s.s[l-1], true
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.s)
}

// Reset completely empties the stack.
func (s *Stack[T]) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = make([]T, 0)
}

/* Example ------------------------------------------------

    stack := Stack.New[int]()
    stack.Push(1)
    stack.Push(2)
    stack.Push(3)

    fmt.Println("Current size:", stack.Size())

    if val, ok := stack.Peek(); ok {
        fmt.Println("Top item:", val)
    }

    for stack.Size() > 0 {
        if val, ok := stack.Pop(); ok {
            fmt.Println("Pop:", val)
        }
    }


----------------------------------------------- */
