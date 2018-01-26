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

// Version is a semantic version for the package.
const Version = "0.2.2"

var renderJobs = make(map[string]RenderJob)
var renderJobsMutex = &sync.Mutex{}

// RenderJob contains all information necessary to complete a Render.
type RenderJob struct {
	Parameters gofr.Parameters
	Cancel     chan bool
	Threads    int
}

// Render executes a RenderJob's unit of work.
func (rj *RenderJob) Render() (image.Image, error) {
	img := image.NewNRGBA64(image.Rect(0, 0, rj.Parameters.ImageWidth, rj.Parameters.ImageHeight))
	contexts := gofr.MakeContexts(img, rj.Threads, &rj.Parameters)

	err := gofr.Render(rj.Threads, contexts, rj.Cancel)
	if err != nil {
		return nil, err
	}

	image := resize.Resize(rj.Parameters.Width, rj.Parameters.Height, image.Image(img), resize.Lanczos3)
	return image, nil
}

// LogResponseWriter logs how long a response took and what it's
// resulting status code was.
type LogResponseWriter struct {
	http.ResponseWriter
	Status int
	Start  time.Time
	End    time.Time
}

// NewLogResponseWriter returns a new LogResponseWriter that wraps a
// given http.ResponseWriter.
func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{w, 0, time.Now(), time.Now()}
}

// WriteHeader implements http.ResponseWriter
func (lrw *LogResponseWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Log writes out the time difference between being initialized and when
// it is called.
func (lrw LogResponseWriter) Log(message string) {
	lrw.End = time.Now()
	log.Printf("%s %v\n", message, lrw.End.Sub(lrw.Start))
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

func makeSPARoute(docroot, index string) http.HandlerFunc {

	// Rule: routes must end in '/', files must not. This means that
	// requests for directories will always be served with the spa index
	// file.
	pattern := regexp.MustCompile(`/$`)

	return func(w http.ResponseWriter, r *http.Request) {
		var reqPath string

		if r.URL.Path == "/" || pattern.MatchString(r.URL.Path) {
			reqPath = path.Join(docroot, index)
		} else {
			reqPath = filepath.Clean(r.URL.Path)
			reqPath = path.Join(docroot, reqPath)
		}

		_, err := os.Stat(reqPath)
		if err != nil {
			finish(w, 404, "File Not Found")
			return
		}

		f, err := os.Open(reqPath)
		defer f.Close()
		if err != nil {
			finish(w, 500, fmt.Sprintf("Unable to open file: %s", reqPath))
			return
		}

		i, err := os.Stat(reqPath)
		if err != nil {
			finish(w, 500, fmt.Sprintf("Unable to stat file: %s", reqPath))
			return
		}

		http.ServeContent(w, r, reqPath, i.ModTime(), f)
	}
}

func routePNG(w http.ResponseWriter, r *http.Request) {
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

	renderID := q.Get("render-id")
	if renderID == "" {
		finish(w, http.StatusUnprocessableEntity, "Missing render-id")
		return
	}

	defer func() {
		renderJobsMutex.Lock()
		delete(renderJobs, renderID)
		renderJobsMutex.Unlock()
	}()

	renderJobsMutex.Lock()
	if _, exists := renderJobs[renderID]; exists {
		close(renderJobs[renderID].Cancel)
		delete(renderJobs, renderID)
	}
	renderJobs[renderID] = j
	renderJobsMutex.Unlock()

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

func routeStatus(w http.ResponseWriter, r *http.Request) {
	finish(w, http.StatusOK, "OK")
}

func main() {
	var value string

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	log.Printf("gofrd v%s", Version)
	log.Printf("libgofrd v%s", gofr.Version)

	staticDir := "./build"
	if value = os.Getenv("GOFR_STATIC_DIR"); value != "" {
		staticDir = value
	}
	staticDir, err := filepath.Abs(staticDir)
	if err != nil {
		panic(err)
	}
	log.Printf("Serving from: %s\n", staticDir)

	bindAddr := "0.0.0.0:8000"
	if value = os.Getenv("GOFR_BIND_ADDR"); value != "" {
		bindAddr = value
	}
	log.Printf("Listening on: %s\n", bindAddr)

	//http.Handle("/", wrapHandler(http.FileServer(http.Dir(staticDir))))
	http.Handle("/", wrapHandlerFunc(makeSPARoute(staticDir, "index.html")))
	http.Handle("/png", wrapHandlerFunc(routePNG))
	http.Handle("/status", wrapHandlerFunc(routeStatus))

	/* Run the thing. */
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
