package parallel

import (
	"fmt"
	"sync"
)

func feedInputs[T any](done <-chan int, inputs []T) (<-chan T, <-chan error) {
	inputsChan := make(chan T)
	errChan := make(chan error, 1)

	go func() {
		defer close(inputsChan)

		errChan <- func() error {
			for _, input := range inputs {
				select {
				case inputsChan <- input:
				case <-done:
					return fmt.Errorf("loop canceled")
				}
			}

			return nil
		}()

	}()

	return inputsChan, errChan
}

type resultWithError[T any] struct {
	result T
	err    error
}

type Worker[TInput any, TResult any] func(input TInput) (TResult, error)

func work[TInput any, TResult any](done <-chan int, inputs <-chan TInput, c chan<- resultWithError[TResult], w Worker[TInput, TResult]) {
	for input := range inputs {
		re := resultWithError[TResult]{}
		re.result, re.err = w(input)
		select {
		case c <- re:
		case <-done:
			return
		}
	}
}

func Run[TInput any, TResult any](workerNum int, w Worker[TInput, TResult], inputs []TInput) ([]TResult, error) {
	done := make(chan int)
	defer close(done)

	inputsChan, errChan := feedInputs(done, inputs)

	c := make(chan resultWithError[TResult])

	var wg sync.WaitGroup

	wg.Add(workerNum)

	for i := 0; i < workerNum; i++ {
		go func() {
			work(done, inputsChan, c, w)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	results := make([]TResult, 0, len(c))

	for r := range c {
		if r.err != nil {
			return nil, r.err
		}

		results = append(results, r.result)
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	return results, nil
}
