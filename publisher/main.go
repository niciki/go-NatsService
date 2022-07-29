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

func InitConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	clusterID := "test-cluster" // nats cluster id
	clientID := "client0"
	example := new(so.Order)
	fmt.Print(example)
	if err := InitConfig(); err != nil {
		log.Fatal(err)
	}
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("127.0.0.1:"+viper.GetString("port_nats")))
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
	data := [][]byte{[]byte(""), []byte("bad input")}
Loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			data[0] = so.GenerateNewOrder(test)
			err = sc.Publish(subj, data[0])
			if err == nil {
				log.Printf("%s sends successfully\n", string(data[0]))
			} else {
				log.Print(err)
			}
		case <-time.After(17 * time.Second):
			err = sc.Publish(subj, data[1])
			if err == nil {
				log.Printf("%s sends successfully\n", string(data[1]))
			} else {
				log.Print(err)
			}
		case <-endChan:
			break Loop
		}
	}
	sc.Close()
}
