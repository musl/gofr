package gofr

type RenderFunc func(*Context) int

func Render(ncpus int, contexts []*Context, render RenderFunc) {
	jobs := make(chan *Context, len(contexts))
	results := make(chan int, len(contexts))
	i := 0

	for i = 0; i < ncpus; i++ {
		go func() {
			for job := range jobs {
				results <- render(job)
			}
		}()
	}

	for _, context := range contexts {
		jobs <- context
	}
	close(jobs)

	for i = 0; i < len(contexts); i++ {
		<-results
	}
	close(results)
}
