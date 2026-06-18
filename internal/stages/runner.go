package stages

import (
	"context/internal/data"
	"runtime"
	"sync"
)

type Detector interface {
	Name() string
	Detect(filePath string) []data.Finding
}

type Runner struct {
	detectors []Detector
	workers   int
}

func NewRunner(detectors []Detector, workers int) *Runner {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &Runner{
		detectors: detectors,
		workers:   workers,
	}
}

func (r *Runner) Run(files []string) []data.Finding {
	var wg sync.WaitGroup

	fileCh := make(chan string, len(files))
	for _, f := range files {
		fileCh <- f
	}
	close(fileCh)

	resultCh := make(chan []data.Finding, r.workers)

	for i := 0; i < r.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var local []data.Finding

			for file := range fileCh {
				for _, det := range r.detectors {
					findings := det.Detect(file)
					if len(findings) > 0 {
						local = append(local, findings...)
					}
				}
			}

			if len(local) > 0 {
				resultCh <- local
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []data.Finding
	for batch := range resultCh {
		results = append(results, batch...)
	}

	return results
}