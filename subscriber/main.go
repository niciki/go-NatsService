package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
	ls "github.com/niciki/go-NatsService/structures/localStore"
	so "github.com/niciki/go-NatsService/structures/structOrder"
)

func main() {
	clusterID := "test-cluster"
	clientID := "client1"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("127.0.0.1:4222"))
	example := new(so.Order)
	fmt.Print(example)
	var store ls.Store
	fmt.Print(store)
	if err != nil {
		log.Fatal(err)
	}
	subj := "sub"
	sub, err := sc.Subscribe(subj, func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Println(err.Error())
	}
	time.Sleep(10 * time.Minute)
	log.Printf("Listening on [%s]", subj)
	sub.Unsubscribe()
	sc.Close()
}
