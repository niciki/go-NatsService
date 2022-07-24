package main

import (
	"encoding/json"
	"log"
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
		json.Unmarshal(m.Data, &newRecord)
		CacheStore.Add(newRecord)
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Listening on [%s]", subj)
	doneCh := make(chan bool)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		sub.Unsubscribe()
		doneCh <- true
	}()
	<-doneCh
	sub.Unsubscribe()
	sc.Close()
}
