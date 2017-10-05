package engine

// timerEvent pointer because we change the index
type priorityQueue []timerEventInterface

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].getTime() < pq[j].getTime() // greater time - less priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].setIdx(i)
	pq[j].setIdx(j)
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(timerEventInterface)
	item.setIdx(n)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	event := old[n-1]
	*pq = old[0 : n-1]
	return event
}
