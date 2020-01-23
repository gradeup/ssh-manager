package apis

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorPage(w http.ResponseWriter) {
	t, err := template.ParseFiles("ui/errorpage.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
