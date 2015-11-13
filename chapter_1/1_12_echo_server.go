package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"io"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"strconv"
)

var mu sync.Mutex
var count int

var palette = []color.Color{color.Black, color.RGBA{0x00,0xff,0x00,0xff}}

const (
	whiteIndex = 0 //first color in palette
	blackIndex = 1 //next color in palette
)
func main() {
	/*handler := func(w http.ResponseWriter, r *http.Request) {
		//do stuff
	}*/


	http.HandleFunc("/", handler)//each request calls handler
	http.HandleFunc("/count", counter)
	http.HandleFunc("/lissajous/",handlerLiss)

	//this is an anonymous function 
	// http.HandleFunc("/lissajous", func (w http.ResponseWriter, r *http.Request) {
	// 	lissajous(w)
	// })
	

	log.Fatal(http.ListenAndServe("localhost:8000",nil))
}

func handlerLiss(w http.ResponseWriter, r *http.Request) {
	
	//if used the other way parseForm must first be called
	// if err := r.ParseForm(); err != nil {
	// 	log.Print(err)
	// }

	//this will print to the server
	// fmt.Println(r.FormValue("cycles"))
	//alternate form access with c[0]
	// c := r.Form["cycles"]

	cyclesParam, err := strconv.Atoi(r.FormValue("cycles")); if err !=  nil {
		log.Print(err)
	}

	lissajous(w, float64(cyclesParam))
}

//handler echoes the path componet of the request url r.
func handler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)

	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}

	//forms field submitted like /?cycles=30
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}

	// lissajous(w);

}

//echos the number of requests so far
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n",count)
	mu.Unlock()
}

func lissajous(out io.Writer, cyclesParam float64) {
	const(
		
		res = 0.001 //angular resolution
		size = 100  //image canvas covers [-size...+size]
		nframes = 64 //number of animation in frames
		delay = 8 	//delay between frames in seconds
	)
	cycles :=  cyclesParam	//number of complete x oscillator revolutions
	freq := rand.Float64() * 3.0 //relative freq of y
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 //phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect,palette)
		for t := 0.0; t < cycles*2*math.Pi; t+= res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		} 
		phase += 0.1
		anim.Delay = append(anim.Delay,delay)
		anim.Image = append(anim.Image,img)
	}
	gif.EncodeAll(out, &anim)//ignoring encoding errors
}