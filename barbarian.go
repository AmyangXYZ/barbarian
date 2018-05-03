package main

import (
	"sync/atomic"
	"time"
)

// Statistics record runtime data.
type Statistics struct {
	Start         time.Time
	Finish        time.Time
	UpTime        time.Duration
	Inputs        int
	Execps        float64
	Execs         uint64
	PositiveExecs uint64
}

// RunHandler handles main work.
type RunHandler func(line string) interface{}

// ResultHandler handles result processing.
type ResultHandler func(result interface{})

// Barbarian is a crazy brute machine.
type Barbarian struct {
	Stats      Statistics
	runHandler RunHandler
	resHandler ResultHandler
	workList   []string
	blockers   chan struct{}
	input      chan string
	output     chan interface{}
	stop       chan bool
}

// New returns a new Barbarin.
func New(runHandler RunHandler, resHandler ResultHandler, workList []string, concurrency int) *Barbarian {
	return &Barbarian{
		Stats:      Statistics{},
		runHandler: runHandler,
		resHandler: resHandler,
		workList:   workList,
		blockers:   make(chan struct{}, concurrency),
		input:      make(chan string),
		output:     make(chan interface{}),
		stop:       make(chan bool),
	}
}

// Run .
func (bb *Barbarian) Run() {
	bb.Stats.Start = time.Now()

	// wait input
	go func() {
		for in := range bb.input {
			go func(in string) {
				atomic.AddUint64(&bb.Stats.Execs, 1)
				res := bb.runHandler(in)
				if res != nil {
					bb.output <- res
					atomic.AddUint64(&bb.Stats.PositiveExecs, 1)
				}
				<-bb.blockers
			}(in)
		}
	}()

	// wait output
	go func() {
		for res := range bb.output {
			bb.resHandler(res)
		}
	}()

	bb.Stats.Inputs = len(bb.workList)
	for _, line := range bb.workList {
		select {
		case <-bb.stop:
			close(bb.input)
			return
		default:
			bb.blockers <- struct{}{}
			bb.input <- line

		}
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(bb.blockers); i++ {
		bb.blockers <- struct{}{}
	}
	bb.Report()
}

// Stop brute.
func (bb *Barbarian) Stop() {
	bb.stop <- true
}

// Report update stats.
func (bb *Barbarian) Report() {
	bb.Stats.Finish = time.Now()
	bb.Stats.UpTime = bb.Stats.Finish.Sub(bb.Stats.Start)
	bb.Stats.Execps = float64(bb.Stats.Execs) / bb.Stats.UpTime.Seconds()
}
