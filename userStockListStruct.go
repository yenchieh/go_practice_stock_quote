package main

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
)

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string `bson:"omitempty"`
	Password string
	Email string
	DateCreated time.Time `json:"dateCreated"`

}

type CustomList struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	User User
	Stock map[string]Stock
}

type Stock struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	StockName string
	Symbol string
}

func (list *CustomList) AddStock(stock map[string]Stock) map[string]Stock {
	for _, v := range stock {
		if lan, ok := list.Stock[v.Symbol]; ok == false {
			fmt.Println(lan)
			list.Stock[v.Symbol] = v
		}
	}



	return list.Stock
}