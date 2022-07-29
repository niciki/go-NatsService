package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
	db "github.com/niciki/go-NatsService/structures/database"
	ls "github.com/niciki/go-NatsService/structures/localStore"
	so "github.com/niciki/go-NatsService/structures/structOrder"
	"github.com/spf13/viper"
)

func waitExit(endChan chan struct{}) {
	endChanSignal := make(chan os.Signal, 1)
	signal.Notify(endChanSignal, os.Interrupt)
	<-endChanSignal
	log.Print("end\n")
	endChan <- struct{}{}
	close(endChanSignal)
	close(endChan)
}

func main() {
	clusterID := "test-cluster"
	clientID := "client1"
	if err := InitConfig(); err != nil {
		log.Fatal(err)
	}
	// initialise database and restore data
	database, err := db.InitDb(viper.GetString("port_postgresql"))
	if err != nil {
		log.Fatal(err)
	}
	// restore data to cache
	CacheStore := ls.NewStore()
	err = database.UploadCache(CacheStore)
	if err != nil {
		log.Fatal(err)
	}
	// connect to nats-stream
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("127.0.0.1:"+viper.GetString("port_nats")))
	newRecord := *new(so.Order)
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
			// add to cache
			err := CacheStore.Add(newRecord)
			if err != nil {
				log.Print(err)
			} else {
				// add to database if record doesn't exist
				err = database.AddRecord(newRecord)
				if err != nil {
					log.Print(err)
				}
				log.Print("+1\n")
			}
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
	server := &http.Server{Addr: ("127.0.0.1:" + viper.GetString("port_http"))}
	server.ListenAndServe()
	endChan := make(chan struct{})
	go waitExit(endChan)
	<-endChan
	sub.Unsubscribe()
	sc.Close()
	// disconnect from port of http server
	server.Shutdown(context.Background())
	// save data in database
	database.AddRecord(CacheStore.GetAll()...)
}

func InitConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
