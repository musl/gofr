package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/musl/gofr/lib/gofr"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const Version = "0.2.2"

var render_jobs = make(map[string]RenderJob)
var render_jobs_mutex = &sync.Mutex{}

type RenderJob struct {
	Parameters gofr.Parameters
	Cancel     chan bool
	Threads    int
}

func (self *RenderJob) Render() (image.Image, error) {
	img := image.NewNRGBA64(image.Rect(0, 0, self.Parameters.ImageWidth, self.Parameters.ImageHeight))
	contexts := gofr.MakeContexts(img, self.Threads, &self.Parameters)

	err := gofr.Render(self.Threads, contexts, self.Cancel)
	if err != nil {
		return nil, err
	}

	image := resize.Resize(self.Parameters.Width, self.Parameters.Height, image.Image(img), resize.Lanczos3)
	return image, nil
}

type LogResponseWriter struct {
	http.ResponseWriter
	Status int
	Start  time.Time
	End    time.Time
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{w, 0, time.Now(), time.Now()}
}

func (self *LogResponseWriter) WriteHeader(code int) {
	self.Status = code
	self.ResponseWriter.WriteHeader(code)
}

func (self LogResponseWriter) Log(message string) {
	self.End = time.Now()
	log.Printf("%s %v\n", message, self.End.Sub(self.Start))
}

func wrapHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLogResponseWriter(w)

		lrw.Log(fmt.Sprintf("%s -> %s %s", r.Method, r.URL.Path, r.RemoteAddr))
		h.ServeHTTP(lrw, r)
		lrw.Log(fmt.Sprintf("%d <- %s %s %s", lrw.Status, r.Method, r.URL.Path, r.RemoteAddr))
	})
}

func wrapHandlerFunc(h http.HandlerFunc) http.Handler {
	return wrapHandler(http.Handler(h))
}

func finish(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	io.WriteString(w, message)
}

func make_spa_route(docroot, index string) http.HandlerFunc {

	// Rule: routes must end in '/', files must not. This means that
	// requests for directories will always be served with the spa index
	// file.
	pattern := regexp.MustCompile(`/$`)

	return func(w http.ResponseWriter, r *http.Request) {
		var req_path string

		if r.URL.Path == "/" || pattern.MatchString(r.URL.Path) {
			req_path = path.Join(docroot, index)
		} else {
			req_path = filepath.Clean(r.URL.Path)
			req_path = path.Join(docroot, req_path)
		}

		_, err := os.Stat(req_path)
		if err != nil {
			finish(w, 404, "File Not Found")
			return
		}

		f, err := os.Open(req_path)
		defer f.Close()
		if err != nil {
			finish(w, 500, fmt.Sprintf("Unable to open file: %s", req_path))
			return
		}

		i, err := os.Stat(req_path)
		if err != nil {
			finish(w, 500, fmt.Sprintf("Unable to stat file: %s", req_path))
			return
		}

		http.ServeContent(w, r, req_path, i.ModTime(), f)
	}
}

func route_png(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()

	// TODO either check a table for currently rendering IDs and return
	// HTTP 429 (one request per id in flight), or cancel the existing,
	// running request and continue this one.

	if r.Method != "GET" {
		finish(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	q := r.URL.Query()

	s, err := strconv.Atoi(q.Get("s"))
	if err != nil {
		s = 1
	}

	e, err := strconv.Atoi(q.Get("p"))
	if err != nil {
		e = 1
	}

	width, err := strconv.Atoi(q.Get("w"))
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid width")
		return
	}

	height, err := strconv.Atoi(q.Get("h"))
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid height")
		return
	}

	iterations, err := strconv.Atoi(q.Get("i"))
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid i")
		return
	}

	er, err := strconv.ParseFloat(q.Get("e"), 64)
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid e")
		return
	}

	rmin, err := strconv.ParseFloat(q.Get("rmin"), 64)
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid rmin")
		return
	}

	imin, err := strconv.ParseFloat(q.Get("imin"), 64)
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid rmin")
		return
	}

	rmax, err := strconv.ParseFloat(q.Get("rmax"), 64)
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid rmax")
		return
	}

	imax, err := strconv.ParseFloat(q.Get("imax"), 64)
	if err != nil {
		finish(w, http.StatusUnprocessableEntity, "Invalid rmin")
		return
	}

	j := RenderJob{
		Parameters: gofr.Parameters{
			Width:        uint(width),
			Height:       uint(height),
			ImageWidth:   width * s,
			ImageHeight:  height * s,
			MaxI:         iterations,
			EscapeRadius: er,
			Min:          complex(rmin, imin),
			Max:          complex(rmax, imax),
			RenderFunc:   q.Get("r"),
			ColorFunc:    q.Get("c"),
			MemberColor:  q.Get("m"),
			Power:        e,
		},
		Threads: runtime.NumCPU(),
		Cancel:  make(chan bool),
	}

	render_id := q.Get("render_id")
	if render_id == "" {
		finish(w, http.StatusUnprocessableEntity, "Missing render_id")
		return
	}

	defer func() {
		render_jobs_mutex.Lock()
		delete(render_jobs, render_id)
		render_jobs_mutex.Unlock()
	}()

	render_jobs_mutex.Lock()
	if _, exists := render_jobs[render_id]; exists {
		close(render_jobs[render_id].Cancel)
		delete(render_jobs, render_id)
	}
	render_jobs[render_id] = j
	render_jobs_mutex.Unlock()

	image, err := j.Render()
	if err != nil {
		finish(w, http.StatusTooManyRequests, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-Render-Job-ID", id.String())
	w.WriteHeader(http.StatusOK)
	png.Encode(w, image)
}

func route_status(w http.ResponseWriter, r *http.Request) {
	finish(w, http.StatusOK, "OK")
}

func main() {
	var value string

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	log.Printf("gofrd v%s", Version)
	log.Printf("libgofrd v%s", gofr.Version)

	static_dir := "./build"
	if value = os.Getenv("GOFR_STATIC_DIR"); value != "" {
		static_dir = value
	}
	static_dir, err := filepath.Abs(static_dir)
	if err != nil {
		panic(err)
	}
	log.Printf("Serving from: %s\n", static_dir)

	bind_addr := "0.0.0.0:8000"
	if value = os.Getenv("GOFR_BIND_ADDR"); value != "" {
		bind_addr = value
	}
	log.Printf("Listening on: %s\n", bind_addr)

	//http.Handle("/", wrapHandler(http.FileServer(http.Dir(static_dir))))
	http.Handle("/", wrapHandlerFunc(make_spa_route(static_dir, "index.html")))
	http.Handle("/png", wrapHandlerFunc(route_png))
	http.Handle("/status", wrapHandlerFunc(route_status))

	/* Run the thing. */
	log.Fatal(http.ListenAndServe(bind_addr, nil))
}
