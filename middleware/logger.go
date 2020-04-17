package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	chiMiddleware "github.com/go-chi/chi/middleware"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white. Logger prints a
// request ID if one is provided.

// RequestLogger returns a logger handler using a custom LogFormatter.
// It will logs the aplication on the default Stdout and on the logWritter given
func RequestLogger(logWriter io.Writer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			DefaultStdout := &chiMiddleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: false}
			DefaultLogger := &chiMiddleware.DefaultLogFormatter{Logger: log.New(logWriter, "", log.LstdFlags), NoColor: true}
			entryStdout := DefaultStdout.NewLogEntry(r)
			entryLog := DefaultLogger.NewLogEntry(r)
			ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Escrevendo na entrada
			entryStdout.Write(ww.Status(), ww.BytesWritten(), ww.Header(), 0, nil)
			entryLog.Write(ww.Status(), ww.BytesWritten(), ww.Header(), 0, nil)

			// Escrevendo na saida com a duração total da requisição
			t1 := time.Now()
			defer func() {
				entryStdout.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
				entryLog.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
