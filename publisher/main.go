package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	so "github.com/niciki/go-NatsService/structures/structOrder"

	"github.com/nats-io/stan.go"
)

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
	rec, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return
	}
	err = sc.Publish(subj, rec)
	if err == nil {
		log.Printf("%s sends successfully\n", string(rec))
	}
	sc.Close()
}
