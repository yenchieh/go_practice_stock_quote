package main

import (
	"net/http"
	"fmt"
	"log"
	"text/template"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"net/url"
	"github.com/rs/cors"
	"time"
)

var templates map[string]*template.Template
var yahooFinanceUrl string = "http://query.yahooapis.com/v1/public/yql"

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("src/github.com/go_practice/index.html"))

}

func main() {
	fmt.Println("Starting...")
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", index).Methods("Get")
	r.HandleFunc("/myStock", index).Methods("Get")
	r.HandleFunc("/resource/{type}/{fileName}", ServeHTTP)
	r.HandleFunc("/search", getQuotesAndRender).Methods("Get")
	r.HandleFunc("/addToList", addToList).Methods("POST")
	r.HandleFunc("/getUserStockList", getUserStockList).Methods("GET")
	r.HandleFunc("/checkUser", checkUser).Methods("GET")
	r.HandleFunc("/removeFromList", removeFromStockList).Methods("POST")
	log.Println("Listening...")

/*	server := &http.Server{
		Addr: ":8080",
		Handler: cors.Default().Handler(r)
	}*/

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":8080", handler)
}

func renderTemplate(w http.ResponseWriter, name string, template interface{}) {
	temp, ok := templates[name]
	if !ok {
		http.Error(w, "Can't find the page", http.StatusBadGateway)
	}

	if err := temp.Execute(w, template); err != nil {
		log.Printf("Error on render template: %s", err)
		http.Error(w, "Serve the wrong page", http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index")
	renderTemplate(w, "index", nil)
}

func getQuotesAndRender(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	fmt.Println(r.URL.RequestURI())
	log.Printf("Search Rquest: %s \n", symbol)

	if len(symbol) < 1 {
		http.Error(w, "The Symbol was not nice", http.StatusInternalServerError)
		return;
	}

	query, err := getQuote(symbol)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	if resp, err := json.Marshal(query); err != nil {
		panic(err)
	}else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func checkUser(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")
	userEmail := r.URL.Query().Get("userEmail")

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("testDB").C("userStockList")

	user := User{}

	if err := c.Find(bson.M{"user.name": userName, "user.email": userEmail}).One(&user); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)

	}else{
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	}

}

func getUserStockList(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")
	userEmail := r.URL.Query().Get("userEmail")

	if len(userName) == 0 || len(userEmail) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := Error{
			Status: string(http.StatusBadRequest),
			Message: "Fail to read user name",
		}
		if jsonResponse, err := json.Marshal(errorResponse); err == nil {
			w.Write(jsonResponse)
		}else {
			panic(err)
		}
		return;
	}

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	defer session.Close()

	log.Printf("Search User: %s %s\n", userName, userEmail)

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("testDB").C("userStockList")

	user := CustomList{}

	err = c.Find(bson.M{"user.name": userName, "user.email": userEmail}).One(&user)

	if err != nil {
		panic(err)
	}

	var stockList string
	for _, i := range user.Stock {
		stockList += i.Symbol + ","
	}

	fmt.Printf("List: %s\n", stockList)

	quotes, err := getQuotes(stockList)

	if(err != nil){
		panic(err)
	}

	var validQuote []StockQuotes
	for _, quote := range quotes.Query.Results.Quote {
		if len(quote.Name) != 0 {
			validQuote = append(validQuote, quote)
		}else{
			//Remove from DB
		}
	}

	quotes.Query.Results.Quote = validQuote

	//put db id into each quote
	for i, _ := range quotes.Query.Results.Quote {
		quotes.Query.Results.Quote[i].Id = user.Stock[i].Id
	}


	if jsonResponse, err := json.Marshal(quotes); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}else {
		panic(err)
	}

}

func requestServer(Url *url.URL) ([]byte, error) {
	log.Printf("Request URL: %s\n", Url.String())
	resp, err := http.Get(Url.String())

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Resposne: %s", body)

	return body, nil

}

func getQuoteFromServer(symbol string)([]byte, error){
	var Url *url.URL

	Url, err := url.Parse(yahooFinanceUrl)

	if err != nil {
		log.Panic(err)
	}
	query := fmt.Sprintf("select finance, Name, Symbol, Change, PercentChange, DaysLow, DaysHigh, Open,PreviousClose, Volume from yahoo.finance.quotes where symbol in (\"%s\")", symbol)
	parameters := url.Values{}
	parameters.Add("q", query)
	parameters.Add("format", "json")
	parameters.Add("env", "http://datatables.org/alltables.env")

	Url.RawQuery = parameters.Encode()

	quoteResultRaw, err := requestServer(Url)

	if err != nil {
		return nil, err
	}

	return quoteResultRaw, nil
}

/**

 */
func getQuote(symbol string) (QueryResult, error) {
	var queryResult QueryResult

	quoteResultRaw, err := getQuoteFromServer(symbol)

	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(quoteResultRaw, &queryResult); err != nil {
		fmt.Printf("Parse JSON Error: %s\n", err.Error())
		return queryResult, err
	}

	return queryResult, nil
}

func getQuotes(symbol string)(QueryResults, error){
	var queryResult QueryResults

	quoteResultRaw, err := getQuoteFromServer(symbol)

	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(quoteResultRaw, &queryResult); err != nil {
		fmt.Printf("Parse JSON Error: %s\n", err.Error())
		return queryResult, err
	}

	return queryResult, nil
}

func addToList(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	var symbolRequest SymbolRequest
	if err := json.Unmarshal(body, &symbolRequest); err != nil {
		panic(err)
	}

	if symbolRequest.Symbol == "" || symbolRequest.UserName == "" {
		http.Error(w, "Error on getting parameter", http.StatusInternalServerError)
		return;
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("testDB").C("userStockList")

	user := User{
		Name: symbolRequest.UserName,
		Email: symbolRequest.UserEmail,
	}

	stock := Stock{StockName: symbolRequest.StockName, Symbol: symbolRequest.Symbol}

	var customList CustomList
	if err = c.Find(bson.M{"user": user}).One(&customList); err != nil {
		//Not found
		log.Printf("The user is not exist %s", user)
		customList := CustomList{
			User: user,
			DateCreated: time.Now(),
		}
		customList.AddStock(stock)

		if err = c.Insert(&customList); err != nil {
			log.Fatal(err)
		}

	}else {
		log.Printf("Existing User %s", user)
		customList.AddStock(stock)
		if err := c.Update(bson.M{"_id": customList.Id}, customList); err != nil {
			log.Panic(err)
		}
	}

	if jsonResponse, err := json.Marshal(customList); err == nil {
		w.Header().Set("Context-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}else {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func removeFromStockList(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	var removeRequest RemoveFromListRequest
	if err := json.Unmarshal(body, &removeRequest); err != nil {
		panic(err)
	}

	if removeRequest.ListId == "" || removeRequest.UserName == "" {
		http.Error(w, "Error on getting parameter", http.StatusInternalServerError)
		return;
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("testDB").C("userStockList")

	user := User{Name: removeRequest.UserName}


	if err := c.Update(bson.M{"user.name": user.Name}, bson.M{"$pull": bson.M{"stock": bson.M{"_id": bson.ObjectIdHex(removeRequest.ListId)}}}); err != nil {
		panic(err)
	}else{
		w.Header().Set("Context-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	}

}

var baseTemplateURL = "src/github.com/go_practice/build"

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileType := vars["type"]
	fileName := vars["fileName"]
	path := fmt.Sprintf("%s/%s", baseTemplateURL, fileName)

	log.Printf("Loading File: %s", path)

	data, err := ioutil.ReadFile(path)
	if err == nil {
		var contentType string
		if strings.HasSuffix(fileType, "css") {
			contentType = "text/css"
		} else if strings.HasSuffix(fileType, "js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(fileType, "png") {
			contentType = "image/png"
		} else {
			contentType = "text/html"
		}
		w.Header().Add("Content-Type", contentType)
		w.Write(data)

	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 Not Found - " + http.StatusText(404)))
	}
}


