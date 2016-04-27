package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "github.com/gtcno/gostudy/ch3/surface/draw"
)
const (
    defaultWidth, defaultHeight = 600, 320
)

func main() {
    http.HandleFunc("/surface", handler) // each request calls handler
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    height, err := strconv.Atoi(r.URL.Query().Get("height"))
    if err != nil {
        height = defaultHeight
    }

    width, err := strconv.Atoi(r.URL.Query().Get("width"))
    if err != nil {
        width = defaultWidth
    }

    w.Header().Set("Content-Type","image/svg+xml")
    fmt.Fprintf(w, "%s", draw.Surface(width,height))
}

