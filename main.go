package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailgun/mailgun-go/v3"
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

	var port = flag.String("port", "3001", "web server port")
	flag.Parse()
	http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
}

func sendConfirmationEmail(to string) error {
	domain := "mg.vaulty.co"
	apiKey := os.Getenv("MG_API_KEY")
	// configure Mailgun client
	mg := mailgun.NewMailgun(domain, apiKey)

	// create http client with using Vaulty as proxy
	httpClient, err := clientWithProxy()
	if err != nil {
		return err
	}

	// configure Mailgun to use client with proxy
	mg.SetClient(httpClient)

	// create message with recipient "to"
	m := mg.NewMessage("Owner <pavel+demo@vaulty.co>", "Confirmation Email", "Hey! Your subscription is confirmed!", to)

	// set timeout to 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// send message
	_, _, err = mg.Send(ctx, m)
	return err
}

func clientWithProxy() (*http.Client, error) {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile("ca.cert")
	if err != nil {
		log.Fatalf("Failed to append ca.cert to RootCAs: %v", err)
		return nil, err
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Fatalf("Failed to parse and add ca.cert to the pool of certificates")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		RootCAs: rootCAs,
	}

	proxyPass := os.Getenv("PROXY_PASS")

	proxyURL, _ := url.Parse(fmt.Sprintf("http://x:%s@vaulty:8080", proxyPass))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return proxyURL, nil
		},
		TLSClientConfig: config,
	}

	client := &http.Client{Transport: tr}
	return client, nil
}
