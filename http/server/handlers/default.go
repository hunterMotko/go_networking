package handlers

import (
	"html/template"
	"io"
	"net/http"
)

var t = template.Must(template.New("hello").Parse("Hello, {{.}}!"))

func DefaultHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
      // drain the request body
      // closing the body will not drain it
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var b []byte

		switch r.Method {
		case http.MethodGet:
			b = []byte("friend")
		case http.MethodPost:
			var err error
			b, err = io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		default:
			// not RFC-compliant due to lack of "Allow" header
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		_ = t.Execute(w, string(b))
	})
}
