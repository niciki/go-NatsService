package order

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/e154/vydumschik"
	"github.com/essentialkaos/translit"
)

type Order struct {
	OrderUid          string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Items   `json:"items"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerId        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	Shardkey          string    `json:"shardkey" db:"shardkey"`
	SmId              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction"`
	RequestId    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Items struct {
	ChrtId      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	Rid         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmId        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
}

func generateNewOrder(test []byte) []byte {
	var order Order
	json.Unmarshal(test, &order)
	rand.Seed(time.Now().UnixNano())
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	// generate order_uid
	length := 19
	var order_uid string
	for i := 0; i < length; i++ {
		order_uid += string(chars[rand.Intn(len(chars))])
	}
	// generate name, surname, lastname
	order.OrderUid = order_uid
	order.Payment.Transaction = order_uid
	var nameGenerator vydumschik.Name
	order.Delivery.Name = translit.EncodeToPCGN(nameGenerator.Full_name(""))
	// generate adress and city
	var adressGenerator vydumschik.Address
	order.Delivery.Address = translit.EncodeToPCGN(adressGenerator.Street_address())
	cities := []string{"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg",
		"Kazan", "Nizhny Novgorod", "Chelyabinsk", "Samara", "Ufa"}
	order.Delivery.City = cities[rand.Intn(9)]
	// generate email
	email := ""
	lengthEmail := rand.Intn(10)
	for i := 0; i < lengthEmail; i++ {
		email += string(chars[rand.Intn(len(chars))])
	}
	order.Delivery.Email = email + "@gmail.com"
	// generate data of payment
	order.Payment.Amount = rand.Intn(20000)
	order.Payment.PaymentDt = rand.Intn(20000000)
	order.Payment.DeliveryCost = rand.Intn(order.Payment.Amount)

	order.DateCreated = time.Now()
	orderByte, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		log.Print(err)
		return []byte{}
	}
	return orderByte
}
