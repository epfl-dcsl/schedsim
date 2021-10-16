## Running a single simulation

`go build`

`./schedsim [OPTION...]`

### Options
* --topo: single queue (0), multi queue (1), bounded queue (2)
* --mu: service rate per core [reqs/us]
* --lambda: arrival rate [reqs/us]
* --genType: MM (0), MD (1), MB[90-10] (2),  MB[99.9-0.1] (3)
* --procType: FIFO processing - number of cores from common.go (0), Processor sharing (1)

#### Examples
`./schedsim --topo=0 --mu=0.1 --lambda=0.005 --genType=2 --procType=0`

## genType Notation
[Kendall’s notation](https://en.wikipedia.org/wiki/Kendall%27s_notation):

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

## Running for multiple arrival rates and configs

Add schedsim to path:

`export PATH="$PATH:${PWD}"` (from where schedsim is)

Example: 

`./scripts/run_new.py "single_queue"`

### Running for multiple arrival rates (partial implementation)
`python3 ./scripts/run_many.py run --topo=0 --mu=0.1 --gen_type=1 --proc_type=0 --num_cores=10`

`python3 ./scripts/run_many.py csv`
