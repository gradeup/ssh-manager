package apis

import (
	"html/template"
	"log"
	"net/http"
	"sshmanager/libraries"
	"time"
)

type HomePageVariables struct {
	Date      string
	Time      string
	PublicKey string
}

func HomePage(w http.ResponseWriter, r *http.Request, publicKeyFile string) {

	publicKey, err := libraries.ReadFile(publicKeyFile)
	if err != nil {
		ErrorPage(w)
		log.Print("template parsing error: ", err)
		return
	}

	now := time.Now()
	HomePageVars := HomePageVariables{
		Date:      now.Format("02-01-2006"),
		Time:      now.Format("15:04:05"),
		PublicKey: publicKey,
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
