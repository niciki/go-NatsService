package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
	ls "github.com/niciki/go-NatsService/structures/localStore"
	so "github.com/niciki/go-NatsService/structures/structOrder"
)

func main() {
	clusterID := "test-cluster"
	clientID := "client1"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("127.0.0.1:4222"))
	newRecord := *new(so.Order)
	CacheStore := ls.NewStore()
	if err != nil {
		log.Fatal(err)
	}
	subj := "sub"
	sub, err := sc.Subscribe(subj, func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &newRecord)
		// wrong model of data
		if err != nil {
			log.Print(err)
		} else {
			err := CacheStore.Add(newRecord)
			if err != nil {
				log.Print(err)
			}
			log.Print("+1\n")
		}
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Listening on [%s]", subj)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			tmpl, err := template.ParseFiles("../template/server.html")
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		case "POST":
			if val, err := CacheStore.Get(req.PostFormValue("order_uid")); err == nil {
				b, err := json.MarshalIndent(val, "", "\t")
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Отправлены данные с order_uid: %s\n", req.PostFormValue("order_uid"))
					fmt.Fprint(w, string(b))
				}
			} else {
				log.Println("Структура с таким order_uid отсутствует")
				fmt.Fprint(w, "Нет записей с таким order_uid ", err)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
	endChan := make(chan struct{})
	go func() {
		endChanSignal := make(chan os.Signal, 1)
		signal.Notify(endChanSignal, os.Interrupt)
		<-endChanSignal
		log.Print("end\n")
		sub.Unsubscribe()
		sc.Close()
		endChan <- struct{}{}
		close(endChanSignal)
		close(endChan)
	}()
	<-endChan
}
