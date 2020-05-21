package pipeline

import (
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/hashicorp/go-multierror"
	"sync"
)

type asyncTaskRunner struct {
	tasks []AsyncTask
}

func NewAsyncTaskRunner(tasks []AsyncTask) *asyncTaskRunner {
	return &asyncTaskRunner{tasks: tasks}
}

func (ar *asyncTaskRunner) Run(p Payload) (Payload, error) {
	var wg sync.WaitGroup
	errCh := make(chan error, len(ar.tasks))
	defer close(errCh)
	for _, task := range ar.tasks {
		wg.Add(1)
		go func(p Payload, errCh chan <- error) {
			task.Run(errCh, p)
			wg.Done()
		}(p, errCh)
	}

	wg.Wait()

	var err error
	for i := 0; i < len(ar.tasks); i++ {
		pErr := <- errCh
		if pErr != nil {
			err = multierror.Append(err, pErr)
		}
	}
	if err != nil {
		return p, errors.NewError("error occurred in the calculator stage", errors.CalculatorError, err)
	}
	return p, nil
}