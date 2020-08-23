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
	http.HandleFunc("/", Handle)
	http.HandleFunc("/timetable/", timetableHandle)
	http.HandleFunc("/save/", saveHandle)
	http.ListenAndServe(":8080", nil)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle")
	page := Page{"index", 1}
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func timetableHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("timetableHandle")
	page := Page{"timetable", 1}
	tmpl, err := template.ParseFiles("timetable.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func saveHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("saveHandle")

	title := r.FormValue("date") + r.FormValue("eventTitle") + ".txt"

	var info string
	info = (r.FormValue("eventTitle") +
		"\nOPEN:" + r.FormValue("openTime") +
		"/START:" + r.FormValue("startTime") +
		"\n前売￥" +
		r.FormValue("beforePrice") +
		"-(DRINK別途￥500-)  " +
		"\n当日￥" +
		r.FormValue("todayPrice") +
		"-(DRINK別途￥500-)\n")

	var artistList, artLen, s, next, timetable, stageStart, stageEnd, h, m string
	i := 1
	startTime := r.FormValue("startTime")

	h = startTime[0:2]
	m = startTime[3:5]
	stageStart = h + ":" + m

	hour, _ := strconv.Atoi(h)
	min, _ := strconv.Atoi(m)

	stageTime, _ := strconv.Atoi(r.FormValue("stageTime"))
	interTime, _ := strconv.Atoi(r.FormValue("interTime"))
	for {
		if i != 1 {
			min += interTime
			if min >= 60 {
				hour += min / 60
				min -= 60
			}

			h = fmt.Sprintf("%02d", hour)
			m = fmt.Sprintf("%02d", min)
			stageStart = (h + ":" + m)
		}
		//終演時間
		min += stageTime
		if min >= 60 {
			hour += min / 60
			min -= 60
		}
		//h = strconv.Itoa(hour)
		//m = strconv.Itoa(min)

		h = fmt.Sprintf("%02d", hour)
		m = fmt.Sprintf("%02d", min)
		stageEnd = (h + ":" + m)
		//転換時間

		//一行目、一番手スタート～終演

		s = strconv.Itoa(i)

		artistList += (r.FormValue("artName" + s))

		//ループの処理

		i++
		next = strconv.Itoa(i)

		artLen = (r.FormValue("artName" + s))
		if len(artLen) == 0 {
			break
		}

		artLen = (r.FormValue("artName" + next))
		if len(artLen) != 0 {
			artistList += " / "
		}
		timetable += stageStart + "～" + stageEnd + " " + (r.FormValue("artName" + s)) + "\n"

	}

	fmt.Println(timetable)

	//インフォメーション文の作成
	infomation := info + artistList
	infomationBody := []byte(infomation + "\n" + timetable)
	ioutil.WriteFile(title, infomationBody, os.ModePerm)

}
