package main

type SymbolRequest struct {
	Symbol string `json:"symbol"`
	StockName string `json:"StockName"`
	UserName string `json:"userName"`
}

type RemoveFromListRequest struct {
	UserName string `json:"userName"`
	ListId string `json:"listId"`
}