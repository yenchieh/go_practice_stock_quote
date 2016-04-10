package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string
	Password string
	Email string
}

type CustomList struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	User User
	Stock []Stock
}

type Stock struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	StockName string
	Symbol string
}