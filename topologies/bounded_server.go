package topologies

import (
	"fmt"

	"github.com/epfl-dcsl/schedsim/blocks"
	"github.com/epfl-dcsl/schedsim/engine"
)

func BoundedQueue(lambda, mu, duration float64, bufferSize int) {

	engine.InitSim()

	//Init the statistics
	stats := &blocks.AllKeeper{}
	stats.SetName("Main Stats")
	engine.InitStats(stats)

	droppedStats := &blocks.AllKeeper{}
	droppedStats.SetName("Dropped Stats")
	engine.InitStats(droppedStats)

	// Add generator
	var g blocks.Generator
	g = blocks.NewMDRandGenerator(lambda, 1/mu)

	g.SetCreator(&blocks.ColoredReqCreator{})

	// Create queues
	q1 := blocks.NewQueue()
	q2 := blocks.NewQueue()

	// Create processors
	p1 := blocks.NewBoundedProcessor(bufferSize)
	p2 := &blocks.BoundedProcessor2{}

	g.AddOutQueue(q1)
	p1.AddInQueue(q1)
	p1.AddOutQueue(q2)
	p1.SetReqDrain(droppedStats)
	engine.RegisterActor(p1)

	p2.AddInQueue(q2)
	p2.SetReqDrain(stats)
	engine.RegisterActor(p2)

	// Register the generator
	engine.RegisterActor(g)

	fmt.Printf("Cores:%v\tservice_rate:%v\tinterarrival_rate:%v\n", cores, mu, lambda)
	engine.Run(duration)
}
