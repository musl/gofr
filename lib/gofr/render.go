package gofr

import (
	"fmt"
)

type RenderFunc func(*Context, chan bool) int

func RenderFuncFromString(name string) (RenderFunc, error) {
	switch name {
	case "mandelbrot":
		return Mandelbrot, nil
	case "ebrot":
		return Ebrot, nil
	case "experimental":
		return Experimental, nil
	default:
		return nil, fmt.Errorf("Invalid RenderFunc name: %#v", name)
	}
}

func Render(threads int, contexts []*Context, cancel chan bool) error {
	count := len(contexts)
	jobs := make(chan *Context, count)
	results := make(chan int, count)

	for i := 0; i < threads; i++ {
		go func() {
			for job := range jobs {
				results <- job.RenderFunc(job, cancel)
			}
		}()
	}

	for _, context := range contexts {
		jobs <- context
	}
	close(jobs)

	r := 0
	for {
		select {
		case <-results:
			r++
			if r == count {
				return nil
			}
		case <-cancel:
			return fmt.Errorf("Render job cancelled.")
		default:
		}
	}
}
