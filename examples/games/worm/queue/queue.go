package queue

import (
	"sync"
)

// Queue is a generic FIFO queue implementation in Go.
type Queue[T any] struct {
	lock sync.Mutex // ensures that the queue is concurrently safe
	q    []T
}

// New creates a new Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds a value to the back of the queue.
func (q *Queue[T]) Enqueue(v T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.q = append(q.q, v)
}

// Dequeue removes the value at the front of the queue and returns it.
// If the queue is empty, it returns the zero value of the type T and a boolean false.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.q) == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	res := q.q[0]
	q.q = q.q[1:]
	return res, true
}

// DequeueFromEnd removes the value at the end of the queue and returns it.
// If the queue is empty, it returns the zero value of the type T and a boolean false.
func (q *Queue[T]) DequeueFromEnd() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.q) == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	res := q.q[len(q.q)-1]

	if len(q.q) == 1 {
		q.q = make([]T, 0)
	} else {
		q.q = q.q[:len(q.q)-1]
	}

	return res, true
}

// Peek returns the value at the front of the queue without removing it.
// If the queue is empty, it returns the zero value of the type T and a boolean false.
func (q *Queue[T]) Peek() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.q) == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	return q.q[0], true
}

// PeekEnd returns the value at the end of the queue without removing it.
// If the queue is empty, it returns the zero value of the type T and a boolean false.
func (q *Queue[T]) PeekEnd() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.q) == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	return q.q[len(q.q)-1], true
}

// InsertAt inserts an element at a specified index in the queue.
// If the index is out of range, it will append the element at the end.
func (q *Queue[T]) InsertAt(index int, v T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if index < 0 || index > len(q.q) {
		q.q = append(q.q, v)
		return
	}
	q.q = append(q.q[:index], append([]T{v}, q.q[index:]...)...)
}

// RemoveAt removes an element at a specified index in the queue.
// Returns the removed element and true if successful, zero value and false if not.
func (q *Queue[T]) RemoveAt(index int) (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if index < 0 || index >= len(q.q) {
		var zero T
		return zero, false
	}

	if len(q.q) == 1 {
		item := q.q[index]
		q.q = make([]T, 0)
		return item, true
	}

	item := q.q[index]
	q.q = append(q.q[:index], q.q[index+1:]...)
	return item, true
}

// PeekAt returns the value at the specified index without removing it.
// Returns the peeked element and true if successful, zero value and false if not.
func (q *Queue[T]) PeekAt(index int) (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if index >= len(q.q) || len(q.q) == 0 {
		var zero T // zero value of type T
		return zero, false
	}

	return q.q[index], true
}

// Size returns the number of elements in the queue.
func (q *Queue[T]) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.q)
}

// Reset completely empties the queue.
func (q *Queue[T]) Reset() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.q = make([]T, 0)
}

/* Example ------------------------------------------------


    queue := Queue.New[int]()
    queue.Enqueue(1)
    queue.Enqueue(2)
    queue.Enqueue(3)

    fmt.Println("Current size:", queue.Size())

    if val, ok := queue.Peek(); ok {
        fmt.Println("Front item:", val)
    }

    for queue.Size() > 0 {
        if val, ok := queue.Dequeue(); ok {
            fmt.Println("Dequeue:", val)
        }
    }


----------------------------------------------- */
