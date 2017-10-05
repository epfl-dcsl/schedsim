package engine

import (
	"math/rand"
)

// Actor is the basic simulation element. Every element (generator or processor)
// should have an actor as a nested struct.
type Actor struct {
	toModel   chan interface{}
	wakeUpCh  chan int
	inQueues  []QueueInterface
	outQueues []QueueInterface
}

func (a *Actor) init(ch chan interface{}) {
	a.toModel = ch
	a.wakeUpCh = make(chan int)
}

// AddInQueue adds another input queue.
// Input queues should be added in decreasing priority
func (a *Actor) AddInQueue(q QueueInterface) {
	a.inQueues = append(a.inQueues, q)
}

// AddOutQueue adds another output queue.
// Output queues should be added in decreasing priority
func (a *Actor) AddOutQueue(q QueueInterface) {
	mdl.queues[q] = true
	a.outQueues = append(a.outQueues, q)
}

// GetInQueueLen returns the length of a given (idx) input queue
func (a *Actor) GetInQueueLen(idx int) int {
	return a.inQueues[idx].Len()
}

// GetOutQueueLen returns the length of a given (idx) output queue
func (a *Actor) GetOutQueueLen(idx int) int {
	return a.outQueues[idx].Len()
}

// GetAllOutQueueLens returns the length of all output queues
func (a *Actor) GetAllOutQueueLens() []int {
	res := make([]int, len(a.outQueues))
	for i, q := range a.outQueues {
		res[i] = q.Len()
	}
	return res
}

// GetAllInQueueLens returns the length of all int queues
func (a *Actor) GetAllInQueueLens() []int {
	res := make([]int, len(a.inQueues))
	for i, q := range a.inQueues {
		res[i] = q.Len()
	}
	return res
}

// GetOutQueueCount returns how many out queues exist
func (a *Actor) GetOutQueueCount() int {
	return len(a.outQueues)
}

// GetInQueueCount returns how many in queues exist
func (a *Actor) GetInQueueCount() int {
	return len(a.inQueues)
}

// Wait blocks the actor for a specific duration d
func (a *Actor) Wait(d float64) {
	e := timerEvent{time: d + mdl.getTime(), wakeUpCh: a.wakeUpCh}
	a.toModel <- e
	<-a.wakeUpCh // block
}

// WaitInterruptible blocks the actor for a d interval, unless there is an
// incoming request in the first input queue.
// Returns true, nil if woken up by the timeout or false, ReqInterface
// if woken up by the incoming req. If red is negative just read input queue
func (a *Actor) WaitInterruptible(d float64) (bool, ReqInterface) {
	if a.inQueues[0].Len() > 0 {
		return false, a.inQueues[0].Dequeue()
	}

	// Negative timeout - no timeout
	if d < 0 {
		return false, a.ReadInQueue()
	}
	timeoutTime := d + mdl.getTime()
	lEvent := linkedEvent{
		timerEvent: timerEvent{time: timeoutTime, wakeUpCh: a.wakeUpCh},
		blockEvent: blockEvent{wakeUpCh: a.wakeUpCh, queues: a.inQueues},
	}
	a.toModel <- lEvent
	<-a.wakeUpCh

	if a.inQueues[0].Len() > 0 {
		return false, a.inQueues[0].Dequeue()
	}
	if mdl.getTime() == timeoutTime {
		return true, nil
	}

	return false, nil
}

// ReadInQueue tries to read the first input queue. If there is a ReqInterface
// available it returns, otherwise the actor blocks
func (a *Actor) ReadInQueue() ReqInterface {
	if a.inQueues[0].Len() > 0 {
		return a.inQueues[0].Dequeue()
	}

	bEvent := blockEvent{wakeUpCh: a.wakeUpCh, queues: a.inQueues}
	a.toModel <- bEvent
	<-a.wakeUpCh
	return a.ReadInQueue()
}

func (a *Actor) ReadInQueueI(idx int) ReqInterface {
	if a.inQueues[idx].Len() > 0 {
		return a.inQueues[idx].Dequeue()
	}

	bEvent := blockEvent{wakeUpCh: a.wakeUpCh, queues: a.inQueues}
	a.toModel <- bEvent
	<-a.wakeUpCh
	return a.ReadInQueueI(idx)
}

// ReadInQueues tries to read from all the queues in descending priority
// and blocks only if all the queues are empty. In returns the element of the
// first queue found non-empty
// Returns the ReqInterface and from which queue it was read
func (a *Actor) ReadInQueues() (ReqInterface, int) {
	for i, q := range a.inQueues {
		if q.Len() > 0 {
			return q.Dequeue(), i
		}
	}

	bEvent := blockEvent{wakeUpCh: a.wakeUpCh, queues: a.inQueues}
	a.toModel <- bEvent
	<-a.wakeUpCh

	return a.ReadInQueues()
}

type queueIdx struct {
	idx int
	q   QueueInterface
}

// ReadInQueuesRand reads from a randomly chosen input queue
// Returns the ReqInterface and from which queue it was read
func (a *Actor) ReadInQueuesRand() (ReqInterface, int) {
	var available []queueIdx
	for i, q := range a.inQueues {
		if q.Len() > 0 {
			available = append(available, queueIdx{i, q})
		}
	}
	if len(available) > 0 {
		q := available[rand.Intn(len(available))]
		return q.q.Dequeue(), q.idx
	}

	bEvent := blockEvent{wakeUpCh: a.wakeUpCh, queues: a.inQueues}
	a.toModel <- bEvent
	<-a.wakeUpCh
	return a.ReadInQueues()
}

// ReadInQueuesRandLocalPr tries to read from the first input queue and if it's
// empty reads randomly from rest incoming queues
// Returns the ReqInterface and from which queue it was read
func (a *Actor) ReadInQueuesRandLocalPr() (ReqInterface, int) {
	if a.inQueues[0].Len() > 0 {
		return a.inQueues[0].Dequeue(), 0
	}
	var available []queueIdx
	for i, q := range a.inQueues {
		if q.Len() > 0 {
			available = append(available, queueIdx{i, q})
		}
	}
	if len(available) > 0 {
		q := available[rand.Intn(len(available))]
		return q.q.Dequeue(), q.idx
	}

	bEvent := blockEvent{wakeUpCh: a.wakeUpCh, queues: a.inQueues}
	a.toModel <- bEvent
	<-a.wakeUpCh
	return a.ReadInQueuesRandLocalPr()
}

// WriteOutQueue writes a ReqInterface to the first output queue
func (a *Actor) WriteOutQueue(el ReqInterface) {
	a.outQueues[0].Enqueue(el)
}

// WriteInQueue writes a ReqInterface to the first input queue
// It can be used for feedback loops
func (a *Actor) WriteInQueue(el ReqInterface) {
	a.inQueues[0].Enqueue(el)
}

// WriteOutQueueI writes a ReqInterface to the given i out queue
func (a *Actor) WriteOutQueueI(el ReqInterface, i int) {
	a.outQueues[i].Enqueue(el)
}

// WriteInQueueI writes a ReqInterface to the given i in queue
func (a *Actor) WriteInQueueI(el ReqInterface, i int) {
	a.inQueues[i].Enqueue(el)
}
