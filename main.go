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
)

var templates map[string]*template.Template
var yahooFinanceUrl string

func init() {
	yahooFinanceUrl = "http://query.yahooapis.com/v1/public/yql"
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("src/github.com/go_practice/index.html"))

}

func main() {


	fmt.Println("Starting...")
	r := mux.NewRouter().StrictSlash(false)
	//mux.HandleFunc("/welcome", index)
	r.HandleFunc("/", index)
	r.HandleFunc("/resource/{type}/{fileName}", ServeHTTP)
	r.HandleFunc("/search", getQuotesAndRender).Methods("Get")
	r.HandleFunc("/addToList", addToList).Methods("POST")
	r.HandleFunc("/getUserStockList", getUserStockList).Methods("GET")
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

	query, err := getQuotes(symbol)
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

func getUserStockList(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")

	if len(userName) == 0 {
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

	log.Printf("Search User: %s\n", userName)

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("testDB").C("userStockList")

/*	searchUser := User {
		Name: userName,
	}*/
	user := []CustomList{}

	err = c.Find(nil).All(&user)

	if err != nil {
		panic(err)
	}

	log.Printf("Find User: %s\n", user[0])

	if jsonResponse, err := json.Marshal(user); err == nil {
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

/**

 */
func getQuotes(symbol string) (QueryResult, error) {
	var queryResult QueryResult
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
		return queryResult, err
	}

	if err = json.Unmarshal(quoteResultRaw, &queryResult); err != nil {
		fmt.Printf("Parse JSON Error: %s\n", err.Error())
		return queryResult, err
	}

	return queryResult, nil
}

func addToList(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	var symbolRequest StoreSymbolRequest
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
	}

	stock := make(map[string]Stock)
	stock[symbolRequest.Symbol] = Stock{StockName: symbolRequest.StockName, Symbol: symbolRequest.Symbol}

	var customList CustomList
	if err = c.Find(bson.M{"user": user}).One(&customList); err != nil {
		//Not found
		log.Printf("The user is not exist %s", user)
		customList := CustomList{
			User: user,
			Stock: stock,
		}

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

var baseTemplateURL = "src/github.com/go_practice"

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fileType := vars["type"]
	fileName := vars["fileName"]
	path := fmt.Sprintf("%s/%s/%s", baseTemplateURL, fileType, fileName)

	log.Printf("Loading File: %s", path)

	fmt.Printf("Here is: %s", r.Method);

	data, err := ioutil.ReadFile(path)
	if err == nil {
		var contentType string
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
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

	//w.Write([]byte("First Name: " + p.fName))
}


