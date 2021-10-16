## Running a single simulation

`go build`

`./schedsim [OPTION...]`

### Options
* topo: single queue (0), multi queue (1), bounded queue (2)
* mu: service rate [reqs/us]
* lambda: arrival rate [reqs/us]
* genType: MM (0), MD (1), MB[90-10] (2),  MB[99.9-0.1]
* procType: Fifo processing - number of cores from common.go (0), Processor sharing (1)

#### Examples
`./schedsim --topo=0 --mu=0.1 --lambda=0.005 --genType=2 --procType=0`

## genType Notation
[Kendall’s notation](https://en.wikipedia.org/wiki/Kendall%27s_notation#:~:text=In%20queueing%20theory%2C%20a%20discipline,and%20classify%20a%20queueing%20node.&text=When%20the%20final%20three%20parameters,%3D%20%E2%88%9E%20and%20D%20%3D%20FIFO.):

A/S/c
* A denotes the time between arrivals to the queue
    * M: Poisson
    * D: fixed inter-arrival time
* S the service time distribution
    * M: Exponential
    * D: Fixed
    * L: Lognormal
    * B: Bimodal
* c the number of service channels open at the node
