package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string `json:"name"`
	Password string
	Email string `json:"emial"`

}

type CustomList struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	User User `json:"user"`
	Stock []Stock `json:"stock"`
	DateCreated time.Time `json:"dateCreated"`
}

type Stock struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	StockName string `json:"stockName"`
	Symbol string `json:"symbol"`
}

func (list *CustomList) AddStock(stock Stock) []Stock {
	var found bool = false
	for _, v := range list.Stock {
		if v.Symbol == stock.Symbol {
			found = true
			break;
		}
	}

	if found == false {
		stock.Id = bson.NewObjectId()
		list.Stock = append(list.Stock, stock)
	}

	return list.Stock
}