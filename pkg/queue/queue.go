/*
The following code was adapted from https://github.com/eapache/queue/blob/main/v2/queue.go
*/
package queue

// MIN_QUEUE_LEN is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const MIN_QUEUE_LEN = 8

// Queue represents a single instance of the queue data structure.
type Queue[V any] struct {
	Buf   []*V
	Head  int
	Tail  int
	count int
}

// New constructs and returns a new Queue.
func New[V any]() *Queue[V] {
	return &Queue[V]{
		Buf: make([]*V, MIN_QUEUE_LEN),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue[V]) Length() int {
	return q.count
}

func (q *Queue[V]) IsEmpty() bool {
	return q.Length() == 0
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue[V]) resize() {
	newBuf := make([]*V, q.count<<1)

	if q.Tail > q.Head {
		copy(newBuf, q.Buf[q.Head:q.Tail])
	} else {
		n := copy(newBuf, q.Buf[q.Head:])
		copy(newBuf[n:], q.Buf[:q.Tail])
	}

	q.Head = 0
	q.Tail = q.count
	q.Buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *Queue[V]) Add(elem V) {
	if q.count == len(q.Buf) {
		q.resize()
	}

	q.Buf[q.Tail] = &elem
	// bitwise modulus
	q.Tail = (q.Tail + 1) & (len(q.Buf) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Queue[V]) Peek() V {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return *(q.Buf[q.Head])
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue[V]) Get(i int) V {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}
	// bitwise modulus
	return *(q.Buf[(q.Head+i)&(len(q.Buf)-1)])
}

// Replace sets the element index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue[V]) Replace(i int, elem V) {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		panic("queue: Replace() called with index out of range")
	}

	q.Buf[(q.Head+i)&(len(q.Buf)-1)] = &elem
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *Queue[V]) Remove() V {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	ret := q.Buf[q.Head]
	q.Buf[q.Head] = nil
	// bitwise modulus
	q.Head = (q.Head + 1) & (len(q.Buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.Buf) > MIN_QUEUE_LEN && (q.count<<2) == len(q.Buf) {
		q.resize()
	}
	return *ret
}

func (q *Queue[V]) Clear() {
	q.Head = 0
	q.Tail = 0
	q.count = 0
}
