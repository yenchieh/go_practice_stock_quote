package main

type SymbolRequest struct {
	Symbol string `json:"symbol"`
	StockName string `json:"StockName"`
	UserName string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

type RemoveFromListRequest struct {
	UserName string `json:"userName"`
	ListId string `json:"listId"`
}