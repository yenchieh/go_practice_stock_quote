package main

type StoreSymbolRequest struct {
	Symbol string `json:"symbol"`
	StockName string `json:"StockName"`
	UserName string `json:"userName"`
}