package main

import "gopkg.in/mgo.v2/bson"

type Fields struct {
	Name string `json:"name"`
	Price string `json:"price"`
	Symbol string `json:"symbol"`
	Ts string `json:"ts"`
	Type string `json:"type"`
	Utctime string `json:"utctime"`
	Volume string `json:"volume"`
}

type SubResource struct {
	Classname string `json:"classname"`
	Fields Fields `json:"fields"`
}
type SubResources struct {
	Resource SubResource `json:"resource"`
}

type Meta struct {
	Type string `json:"type"`
	start int `json:"start"`
	count int `json:"count"`
}

type Quotes struct {
	Meta Meta `json:"meta"`
	Resources  []SubResources `json:"resources"`
}

type MainList struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	List Quotes `json:"list"`
}