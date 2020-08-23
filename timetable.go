package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Page struct {
	Title string
	Count int
}

func main() {
	http.HandleFunc("/test/", testHandle)
	http.HandleFunc("/save/", saveHandle)
	http.ListenAndServe(":8080", nil)
}

func testHandle(w http.ResponseWriter, r *http.Request) {
	page := Page{"Hello World.", 1}
	tmpl, err := template.ParseFiles("test.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}

}

func saveHandle(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	event := r.FormValue("event")
	title := date + " " + event + ".txt"
	open := []byte(r.FormValue("open"))
	start := []byte(r.FormValue("start"))

	//時間修正

	stage := []byte(r.FormValue("stage"))
	interval := r.FormValue("interval")
	var a, b string = stage[0:2], stage[3:5]
	var sh, sm string
	var time, artname []byte
	var h, m, st, lt int
	h, _ = strconv.Atoi(a)
	m, _ = strconv.Atoi(b)
	st, _ = strconv.Atoi(stage)
	lt, _ = strconv.Atoi(interval)

	artistname := r.FormValue("art1")
	timetable := []byte(artistname)
	//--------------------------
	for i := 2; i >= 10; i++ {
		sh = strconv.Itoa(h)
		sm = strconv.Itoa(m)
		time = []byte("/n" + sh + ":" + sm)
		timetable = append(timetable, time...)

		m += st
		if m >= 60 {
			h += m / 60
			m -= 60
		}

		sh = strconv.Itoa(h)
		sm = strconv.Itoa(m)

		time = []byte("～" + sh + ":" + sm + "  ")

		timetable = append(timetable, time...)

		fmt.Println(i)
		var s string
		s = strconv.Itoa(i)
		artistname = r.FormValue("art" + s)
		if len(artistname) == 0 {
			break
		}
		artname = []byte(artistname)
		timetable = append(timetable, artname...)

	}
	//-------------------------
	head := []byte(event, +"\n", +"OPEN ", +open, +"/", +"START ")
	(start, +"\nステージ :", +stage, +"分 /", +" 転換 : ", +interval, +"分\n")
	body := []byte(head, +timetable)
	ioutil.WriteFile(title, head, os.ModePerm)
	ioutil.WriteFile("タイムテーブル.txt", timetable, os.ModePerm)
}
