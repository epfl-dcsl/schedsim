package blocks

import (
	"container/list"
	//	"fmt"

	"github.com/epfl-dcsl/schedsim/engine"
)

// Processor Interface describes the main processor functionality used
// in describing a topology
type Processor interface {
	engine.ActorInterface
	SetReqDrain(rd RequestDrain) // We might want to specify different drains for different processors or use the same drain for all
	SetCtxCost(cost float64)
}

// generic processor: All processors should have it as an embedded field
type genericProcessor struct {
	engine.Actor
	reqDrain RequestDrain
	ctxCost  float64
}

func (p *genericProcessor) SetReqDrain(rd RequestDrain) {
	p.reqDrain = rd
}

func (p *genericProcessor) SetCtxCost(cost float64) {
	p.ctxCost = cost
}

// RTCProcessor is a run to completion processor
type RTCProcessor struct {
	genericProcessor
	scale float64
}

// Run is the main processor loop
func (p *RTCProcessor) Run() {
	for {
		req := p.ReadInQueue()
		p.Wait(req.GetServiceTime() + p.ctxCost)
		if monitorReq, ok := req.(*MonitorReq); ok {
			monitorReq.finalLength = p.GetInQueueLen(0)
		}
		p.reqDrain.TerminateReq(req)
	}
}

// TSProcessor is a time sharing processor
type TSProcessor struct {
	genericProcessor
	quantum float64
}

// NewTSProcessor returns a new *TSProcessor
func NewTSProcessor(quantum float64) *TSProcessor {
	return &TSProcessor{quantum: quantum}
}

// Run is the main processor loop
func (p *TSProcessor) Run() {
	for {
		req := p.ReadInQueue()

		if req.GetServiceTime() <= p.quantum {
			p.Wait(req.GetServiceTime() + p.ctxCost)
			p.reqDrain.TerminateReq(req)
		} else {
			p.Wait(p.quantum + p.ctxCost)
			req.SubServiceTime(p.quantum)
			p.WriteInQueue(req)
		}
	}
}

// PSProcessor is a processor sharing processor
type PSProcessor struct {
	genericProcessor
	workerCount int
	count       int // how many concurrent requests
	reqList     *list.List
	curr        *list.Element
	prevTime    float64
}

// NewPSProcessor returns a new *PSProcessor
func NewPSProcessor() *PSProcessor {
	return &PSProcessor{workerCount: 1, reqList: list.New()}
}

// SetWorkerCount sets the number of workers in a processor sharing processor
func (p *PSProcessor) SetWorkerCount(count int) {
	p.workerCount = count
}

func (p *PSProcessor) getMinService() *list.Element {
	minS := p.reqList.Front().Value.(*Request).ServiceTime
	minI := p.reqList.Front()
	for e := p.reqList.Front(); e != nil; e = e.Next() {
		val := e.Value.(*Request).ServiceTime
		if val < minS {
			minS = val
			minI = e
		}
	}
	return minI
}

func (p *PSProcessor) getFactor() float64 {
	if p.workerCount > p.count {
		return 1.0
	}
	return float64(p.workerCount) / float64(p.count)
}

func (p *PSProcessor) updateServiceTimes() {
	currTime := engine.GetTime()
	diff := (currTime - p.prevTime) * p.getFactor()
	p.prevTime = currTime
	for e := p.reqList.Front(); e != nil; e = e.Next() {
		req := e.Value.(engine.ReqInterface)
		req.SubServiceTime(diff)
	}
}

// Run is the main processor loop
func (p *PSProcessor) Run() {
	var d float64
	d = -1
	for {
		intr, newReq := p.WaitInterruptible(d)
		//update times
		p.updateServiceTimes()
		if intr {
			req := p.curr.Value.(engine.ReqInterface)
			p.reqDrain.TerminateReq(req)
			p.reqList.Remove(p.curr)
			p.count--
		} else {
			p.count++
			p.reqList.PushBack(newReq)
		}
		if p.count > 0 {
			p.curr = p.getMinService()
			d = p.curr.Value.(engine.ReqInterface).GetServiceTime() / p.getFactor()
		} else {
			d = -1
		}
	}
}
