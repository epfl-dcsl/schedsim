package blocks

import (
	//"container/heap"
	"container/list"
	//"sort"
	//"fmt"
	"github.com/epfl-dcsl/schedsim/engine"
)

var count = 0

// Queue is a imple FIFO queue
type Queue struct {
	l  *list.List
	id int
}

// NewQueue returns a new *Queue
func NewQueue() *Queue {
	q := &Queue{}
	q.l = list.New()
	q.id = count
	count++
	return q
}

// Enqueue enqueues a new ReqInterface at the queue
func (q *Queue) Enqueue(el engine.ReqInterface) {
	//fmt.Printf("time: %v, queue: %v, len: %v\n", engine.GetTime(), q.id, q.Len())
	q.l.PushBack(el)
}

// Dequeue dequeues the last ReqInterface from the queue
func (q *Queue) Dequeue() engine.ReqInterface {
	el := q.l.Front()
	q.l.Remove(el)
	return el.Value.(engine.ReqInterface)
}

// Len returns the queue length
func (q *Queue) Len() int {
	return q.l.Len()
}
