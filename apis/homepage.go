package apis

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageVariables struct {
	Date string
	Time string
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	HomePageVars := PageVariables{
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("ui/homepage.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, HomePageVars)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
