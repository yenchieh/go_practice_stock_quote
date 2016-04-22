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
	Quote StockQuote `json:"quote"`
}

type StockQuote struct {
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


type QueryResults struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Query Queries `json:"query"`
}

type Queries struct {
	Count int `json:"count"`
	Created string `json:"created"`
	Lang string `json:"lang"`
	Results Quotes `json:"results"`
}

type Quotes struct {
	Quote []StockQuotes `json:"quote"`
}

type StockQuotes struct {
	Id bson.ObjectId `json:"id"`
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