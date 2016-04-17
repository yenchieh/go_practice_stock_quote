package main

import "gopkg.in/mgo.v2/bson"

type QueryResult struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Query Query `json:"query"`
}

type Query struct {
	Count int `json:"count"`
	Created string `json:"created"`
	Lang string `json:"lang"`
	Results Quote `json:"results"`
}

type Quote struct {
	Quote StockQuotes `json:"quote"`
}

type StockQuotes struct {
	Change string `json:"change"`
	PercentChange string `json:"percentChange"`
	DaysLow string `json:"daysLow"`
	DaysHigh string `json:"daysHigh"`
	Open string `json:"open"`
	PreviousClose string `json:"previousClose"`
	Symbol string `json:"symbol"`
	Name string `json:"name"`
	Volume string `json:"volume"`
}