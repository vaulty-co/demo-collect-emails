package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailgun/mailgun-go/v4"
)

type email struct {
	Address string
}

func main() {
	db := &DB{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// add subscriber to the list and send confirmation email
	r.Post("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.FormValue("email")

		db.Save(email{address})
		err := sendConfirmationEmail(address)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			w.Write([]byte("Error"))
			return
		}

		http.Redirect(w, r, "/success.html", 301)
	})

	// return list of subscribers
	r.Get("/subscribers", func(w http.ResponseWriter, r *http.Request) {
		type pageData struct {
			Emails []email
		}
		tmpl, err := template.ParseFiles("tpl/list.html")
		if err != nil {
			log.Fatal(err)
		}
		data := &pageData{
			Emails: db.List(),
		}
		tmpl.Execute(w, data)
	})

	// serve static content
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	r.Get("/*", http.FileServer(http.Dir(filesDir)).ServeHTTP)

	var port = flag.String("port", "3000", "web server port")
	flag.Parse()
	http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
}

func sendConfirmationEmail(to string) error {
	// create Mailgun client using env variables (MG_API_KEY, MG_DOMAIN)
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		return err
	}

	// create message with recipient "to"
	m := mg.NewMessage("Owner <pavel+demo@vaulty.co>", "Confirmation Email", "Hey! Your subscription is confirmed!", to)

	// set timeout to 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// send message
	_, _, err = mg.Send(ctx, m)
	return err
}
