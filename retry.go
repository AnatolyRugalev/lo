package lo

import (
	"sync"
	"time"
)

type debounce struct {
	after     time.Duration
	mu        *sync.Mutex
	timer     *time.Timer
	done      bool
	callbacks []func()
}

func (d *debounce) reset() *debounce {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.done {
		return d
	}

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.after, func() {
		for _, f := range d.callbacks {
			f()
		}
	})
	return d
}

func (d *debounce) cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	d.done = true
}

// NewDebounce creates a debounced instance that delays invoking functions given until after wait milliseconds have elapsed.
// Play: https://go.dev/play/p/mz32VMK2nqe
func NewDebounce(duration time.Duration, f ...func()) (func(), func()) {
	d := &debounce{
		after:     duration,
		mu:        new(sync.Mutex),
		timer:     nil,
		done:      false,
		callbacks: f,
	}

	return func() {
		d.reset()
	}, d.cancel
}

// Attempt invokes a function N times until it returns valid output. Returning either the caught error or nil. When first argument is less than `1`, the function runs until a successful response is returned.
// Play: https://go.dev/play/p/3ggJZ2ZKcMj
func Attempt(maxIteration int, f func(index int) error) (int, error) {
	var err error

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		// for retries >= 0 {
		err = f(i)
		if err == nil {
			return i + 1, nil
		}
	}

	return maxIteration, err
}

// AttemptWithDelay invokes a function N times until it returns valid output,
// with a pause between each call. Returning either the caught error or nil.
// When first argument is less than `1`, the function runs until a successful
// response is returned.
// Play: https://go.dev/play/p/tVs6CygC7m1
func AttemptWithDelay(maxIteration int, delay time.Duration, f func(index int, duration time.Duration) error) (int, time.Duration, error) {
	var err error

	start := time.Now()

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		err = f(i, time.Since(start))
		if err == nil {
			return i + 1, time.Since(start), nil
		}

		if maxIteration <= 0 || i+1 < maxIteration {
			time.Sleep(delay)
		}
	}

	return maxIteration, time.Since(start), err
}

// AttemptWhile invokes a function N times until it returns valid output.
// Returning either the caught error or nil, and along with a bool value to identify
// whether it needs invoke function continuously. It will terminate the invoke
// immediately if second bool value is returned with falsy value. When first
// argument is less than `1`, the function runs until a successful response is
// returned.
func AttemptWhile(maxIteration int, f func(int) (error, bool)) (int, error) {
	var err error
	var shouldContinueInvoke bool

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		// for retries >= 0 {
		err, shouldContinueInvoke = f(i)
		if !shouldContinueInvoke { // if shouldContinueInvoke is false, then return immediately
			return i + 1, err
		}
		if err == nil {
			return i + 1, nil
		}
	}

	return maxIteration, err
}

// AttemptWhileWithDelay invokes a function N times until it returns valid output,
// with a pause between each call. Returning either the caught error or nil, and along
// with a bool value to identify whether it needs to invoke function continuously.
// It will terminate the invoke immediately if second bool value is returned with falsy
// value. When first argument is less than `1`, the function runs until a successful
// response is returned.
func AttemptWhileWithDelay(maxIteration int, delay time.Duration, f func(int, time.Duration) (error, bool)) (int, time.Duration, error) {
	var err error
	var shouldContinueInvoke bool

	start := time.Now()

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		err, shouldContinueInvoke = f(i, time.Since(start))
		if !shouldContinueInvoke { // if shouldContinueInvoke is false, then return immediately
			return i + 1, time.Since(start), err
		}
		if err == nil {
			return i + 1, time.Since(start), nil
		}

		if maxIteration <= 0 || i+1 < maxIteration {
			time.Sleep(delay)
		}
	}

	return maxIteration, time.Since(start), err
}

// throttle ?
