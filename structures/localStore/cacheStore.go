package store

import (
	"errors"

	so "github.com/niciki/go-NatsService/structures/structOrder"
)

type Store struct {
	data map[string]so.Order
}

func (s *Store) Add(order so.Order) error {
	_, ok := s.data[order.OrderUid]
	if ok {
		return errors.New("record with this OrderUid exists in database")
	}
	s.data[order.OrderUid] = order
	return nil
}

func (s *Store) Get(orderUid string) (so.Order, error) {
	val, ok := s.data[orderUid]
	if ok {
		return val, nil
	}
	return so.Order{}, errors.New("there isn't record with such OrderUid")
}
