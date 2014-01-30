package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"text/template"
	"time"
)

type Tuple struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func wsHandler(f func(float64) float64) websocket.Handler {
	return func(ws *websocket.Conn) {
		log.Printf("Client connected: %s", ws.Request().RemoteAddr)
		tuple := Tuple{}
		for x := 0.0; ; x += (math.Pi / 30) {
			tuple.X, tuple.Y = x, f(x)
			err := websocket.JSON.Send(ws, tuple)
			if err != nil {
				log.Printf("Client disconnected: %s, %v", ws.Request().RemoteAddr, err)
				return
			}
			time.Sleep(30 * time.Millisecond)
		}
	}
}

func main() {

	http.Handle("/ws", websocket.Handler(wsHandler(math.Sin)))

	// lol, who needs file servers anyway
	http.HandleFunc("/", HTMLPageHandler("index.html", []string{
		"js/lib/d3.js",
		"js/lib/rickshaw.js",
		"js/graphs.js",
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Silly stuff is very silly

func HTMLPageHandler(pageTmplname string, jsFilename []string) http.HandlerFunc {
	page := NewHTMLPage(pageTmplname, jsFilename)
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		io.Copy(rw, bytes.NewReader(page))
		log.Printf("Done in %v", time.Since(start))
	}
}

func NewHTMLPage(pageTmplname string, jsFilename []string) []byte {
	jsBuf := bytes.NewBuffer(nil)
	for _, filename := range jsFilename {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		jsBuf.WriteString("// Starting '" + filename + "'\n")
		jsBuf.Write(data)
		jsBuf.WriteRune('\n')
	}
	pageBuf := bytes.NewBuffer(nil)

	templ := template.Must(template.ParseFiles(pageTmplname))
	err := templ.Execute(pageBuf, struct {
		Libraries string
	}{
		jsBuf.String(),
	})
	if err != nil {
		panic(err)
	}
	return pageBuf.Bytes()
}
