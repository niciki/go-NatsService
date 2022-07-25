package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/stan.go"
	so "github.com/niciki/go-NatsService/structures/structOrder"
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
	clusterID := "test-cluster" // nats cluster id
	clientID := "client0"
	example := new(so.Order)
	fmt.Print(example)
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("127.0.0.1:4222"))
	if err != nil {
		log.Fatal(err)
	}
	subj := "sub"
	f, err := os.Open("../model.json")
	if err != nil {
		fmt.Println("read file fail", err)
		return
	}
	defer f.Close()
	test, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return
	}
	endChan := make(chan struct{})
	go waitExit(endChan)
	// sends test order
	err = sc.Publish(subj, test)
	if err == nil {
		log.Printf("%s sends successfully\n", string(test))
	}
Loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			rec := so.GenerateNewOrder(test)
			err = sc.Publish(subj, rec)
			if err == nil {
				log.Printf("%s sends successfully\n", string(rec))
			} else {
				log.Print(err)
			}
		case <-endChan:
			break Loop
		}
	}
	sc.Close()
}
