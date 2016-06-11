package main

import (
  "speedserver"
  "net/http"
)

func main () {
  http.HandleFunc("/generate", speedserver.GenerateHandler)
  http.ListenAndServe(":3000", nil)
}
