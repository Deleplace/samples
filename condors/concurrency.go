package condors

import "sync"

// RunConcurrent launches funcs,
// and waits for their completion.
func RunConcurrent(funcs ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		f := f
		go func() {
			f()
			wg.Done()
		}()
	}
	wg.Wait()
}
