package routine

import (
	"context"
	"sync"
	"time"
)

// Routine represents an operation which should be run on a cycle.
type Routine struct {
	// Period is the length between runs of Op.
	Period time.Duration
	// Op must be safe for parallel execution.
	Op func(context.Context, time.Time) error
}

// RunRoutines will start go routines to cyclically run the defined operations.
// RunRoutines also starts a
// Canceling ctx will end all running routines and clean up all resources.
// errorHandler must be safe for parallel execution.
func RunRoutines(ctx context.Context, errorHandler func(error), routines []Routine) {
	rr := routineRunner{
		wg:          &sync.WaitGroup{},
		errorStream: make(chan error),
	}

	for _, routine := range routines {
		rr.wg.Add(1)
		go rr.runRoutine(ctx, routine)
	}

	go rr.handleErrors(errorHandler)

	go func() {
		rr.wg.Wait()
		close(rr.errorStream)
	}()
}

type routineRunner struct {
	wg          *sync.WaitGroup
	errorStream chan error
}

func (rr routineRunner) runRoutine(ctx context.Context, routine Routine) {
	ticker := time.NewTicker(routine.Period)
	for {
		select {
		case t := <-ticker.C:
			if err := routine.Op(ctx, t); err != nil {
				rr.errorStream <- err
			}
		case <-ctx.Done():
			rr.wg.Done()
			return
		}
	}

}

// handleErrors handles all errors produced by running routines.
// handleErrors accepts no context because
// the a context canclation of RunRoutines will eventually result
// in the error channel being closed, thus ending handleErrors.
func (rr routineRunner) handleErrors(handler func(error)) {
	for {
		err, ok := <-rr.errorStream
		if !ok {
			return
		}
		handler(err)
	}
}
