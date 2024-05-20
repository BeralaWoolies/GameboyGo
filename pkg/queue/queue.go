package queue

import (
	"fmt"
)

type Queue[T any] struct {
	data   []T
	isFull bool
	start  int
	end    int
}

func NewQueue[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		data:   make([]T, capacity),
		isFull: false,
		start:  0,
		end:    0,
	}
}

func (q *Queue[T]) String() string {
	return fmt.Sprintf(
		"[Queue full:%v size:%d start:%d end:%d data:%v]",
		q.isFull,
		len(q.data),
		q.start,
		q.end,
		q.data)
}

func (q *Queue[T]) Enqueue(elem T) error {
	if q.isFull {
		return fmt.Errorf("Queue is full")
	}

	q.data[q.end] = elem
	q.end = (q.end + 1) % len(q.data)
	q.isFull = q.end == q.start

	return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
	var res T
	if !q.isFull && q.start == q.end {
		return res, fmt.Errorf("Queue is empty")
	}

	res = q.data[q.start]
	q.start = (q.start + 1) % len(q.data)
	q.isFull = false

	return res, nil
}

func (q *Queue[T]) Peek() (T, error) {
	var res T
	if !q.isFull && q.start == q.end {
		return res, fmt.Errorf("Queue is empty")
	}

	return q.data[q.start], nil
}

func (q *Queue[T]) Size() int {
	res := q.end - q.start
	if res < 0 || (res == 0 && q.isFull) {
		res = len(q.data) - res
	}

	return res
}

func (q *Queue[T]) IsFull() bool {
	return q.isFull
}

func (q *Queue[T]) Clear() {
	q.start = 0
	q.end = 0
	q.isFull = false
}
