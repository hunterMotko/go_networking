package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

type handler struct {
  db *sql.DB
  log *log.Logger
}

func (h *handler) handler1() http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      err := h.db.Ping()
      if err != nil {
        h.log.Printf("db ping: %v", err)
      }
      // do something with the db here
    },
  )
}

func (h *handler) handler2() http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // ...
    },
  )
}

func some() {
  // example
  // h := &handler{
  //   db: &sql.DB{},
  //   log: &log.Logger{},
  // }
  // http.Handle("/one", h.handler1())
  // http.Handle("/two", h.handler2())
}
