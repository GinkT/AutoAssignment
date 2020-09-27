package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/AutoAssignment/db"
	"github.com/gorilla/mux"
	"hash/crc32"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Env struct {
	Database *sql.DB
}

func main() {
	Database, err := db.NewDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer Database.Close()
	env := &Env{Database: Database}

	router := mux.NewRouter()
	router.HandleFunc("/{id}", env.RedirectHandler)
	router.Handle("/", ValidLink(http.HandlerFunc(env.ShortenerHandler)))

	srv := &http.Server{Handler: router}
	ln, err := net.Listen("tcp", ":8181")
	if err != nil {
		log.Fatalln(err)
	}
	go func () {
		log.Fatalln(srv.Serve(ln))
	} ()

	// Graceful shutdown
	log.Println("Started to listen and serve at :8181")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

// Хэндлер который выполняет перенаправление на сохраненную ссылку
func (env *Env)RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	log.Println("[Redirect] Got request with path:", path)
	link, err := db.GetLinkFromDB(env.Database, path)
	if err != nil {
		log.Println(err)
		Json404Response(w, "Internal error")
		return
	}
	log.Println("[Redirect] Successfully redirected user to:", link)
	http.Redirect(w, r, link, http.StatusMovedPermanently)
}

// Хэндлер который выполняет укорачивание ссылки
func (env *Env)ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	linkToShort := r.URL.Query().Get("link")
	customLink := r.URL.Query().Get("custom")

	var shortenedLink string
	switch {
	case customLink != "":
		shortenedLink = customLink
		log.Printf("[Shortener] Got a request with link: %s Shortened it to(custom): %s\n", linkToShort, shortenedLink)

		err := db.AddLinkToDB(env.Database, shortenedLink, linkToShort)
		switch  {
		case err == db.ErrorAlreadyInDB:
			Json404Response(w, "Custom link is already taken!")
			log.Println("[DB]", err)
			return
		case err != nil:
			Json404Response(w, "Internal error")
			log.Println("[DB]", err)
			return
		}
	default:
		shortenedLink = ShrinkLink(linkToShort)
		log.Printf("[Shortener] Got a request with link: %s Shortened it to: %s\n", linkToShort, shortenedLink)
		err := db.AddLinkToDB(env.Database, shortenedLink, linkToShort)
		// ErrorAlreadyInDB не критична.
		if err != nil && err != db.ErrorAlreadyInDB {
			Json404Response(w, "Internal error")
			return
		}
	}

	JsonOKResponse(w, linkToShort, shortenedLink)
}

// Middleware для валидации ссылки
func ValidLink(shortener http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("link")
		_, err := url.ParseRequestURI(link)
		if err != nil {
			log.Printf("[Shortener] User link(%s) does not match URL\n", link)
			Json404Response(w, "Invalid link!")
			return
		}
		if customLink := r.URL.Query().Get("custom"); customLink != "" {
			log.Println("[DEBUG!!] CustomLink", customLink)
			pattern := `^\w+$`
			re, err := regexp.Match(pattern, []byte(customLink))
			if err != nil {
				log.Println("RegExp error!")
				Json404Response(w, "Internal server error")
				return
			}
			if !re {
				Json404Response(w, "Invalid custom link!")
				return
			}
		}
		shortener.ServeHTTP(w, r)
	})
}

// Подсчитывает crc32 сумму предложенной ссылки. Переводит в строку
func ShrinkLink(data string) string {
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 36)
	return dataHash
}

// ----------------------------------------------------- JSON Encoders

type JsonOKStruct struct {
	Status 			int			`json:"status"`
	Link			string		`json:"link"`
	ConvertedTo 	string		`json:"convertedTo"`
}

func JsonOKResponse(w http.ResponseWriter, link, convertedTo string) {
	response := &JsonOKStruct{
		Status:        	200,
		Link:			link,
		ConvertedTo: 	convertedTo,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

type Json404Struct struct {
	Status 			int 		`json:"status"`
	Message			string		`json:"message"`
}

func Json404Response(w http.ResponseWriter, message string) {
	response := &Json404Struct{
		Status:        	404,
		Message: 		message,
	}
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(response)
}

