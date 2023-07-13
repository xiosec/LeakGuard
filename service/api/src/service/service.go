package service

import (
	"LeakGuard/config"
	"LeakGuard/databases"
	"LeakGuard/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func index(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Leaked passwords verification service (%s).", utils.VERSION)
	w.Write([]byte(message))
}

func check(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("value") == "" || r.PostFormValue("token") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password := r.Form.Get("value")
	token := r.Form.Get("token")

	if token != config.Conf.Service.Token {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	startTime := time.Now()
	exist, err := databases.ExistPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	endTime := time.Now()

	if exist {
		host := strings.Split(r.RemoteAddr, ":")[0]
		fmt.Printf("[%s] Client %s used a leaked password!\n", time.Now().Format(time.RFC822), host)
		fmt.Printf("\t\\__search in : %s\n", endTime.Sub(startTime))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Run(config config.Service) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/check", check)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
