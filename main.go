package main

import (
	"net/http"
	"fmt"
	"log"
	"text/template"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

var templates map[string]*template.Template



func init(){
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
	r.HandleFunc("/addToList", addToList).Methods("Get")
	log.Println("Listening...")

	server := &http.Server{
		Addr: ":8080",
		Handler: handlers.CORS()(r),
	}

	server.ListenAndServe()
}

func renderTemplate(w http.ResponseWriter, name string, template interface{}){
	temp, ok := templates[name]
	if !ok {
		http.Error(w, "Can't find the page", http.StatusBadGateway)
	}

	if err := temp.Execute(w, template); err != nil {
		log.Printf("Error on render template: %s", err)
		http.Error(w, "Serve the wrong page", http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request){
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

	mainList, err := getQuotes(symbol)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	if resp, err := json.Marshal(mainList.List.Resources); err != nil {
		panic(err)
	}else{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func requestServer(url string, symbol string) (MainList, error){
	var quote MainList
	requestUrl := fmt.Sprintf(url, symbol)
	log.Printf("Request URL: %s\n", requestUrl)
	resp, err := http.Get(requestUrl)

	if err != nil {
		return quote, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return quote, err
	}

	log.Printf("Resposne: %s", body)


	if err = json.Unmarshal(body, &quote); err != nil {
		fmt.Printf("Parse JSON Error: %s\n", err.Error())
		return quote, err
	}

	return quote, nil

}

/**

 */
func getQuotes(symbol string) (MainList, error){
	mainList, err := requestServer("http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json", symbol)

	if err != nil {
		return mainList, err
	}

	return mainList, nil
}

func addToList(w http.ResponseWriter, r *http.Request){
	symbol := r.URL.Query().Get("symbolForAdd")
	stockName := r.URL.Query().Get("stockName")
	username := r.URL.Query().Get("username")
	if symbol == "" || stockName == "" || username == "" {
		http.Error(w, "Error on getting parameter", http.StatusInternalServerError)
		return;
	}

	//mainList, err := getQuotes(symbol)

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("testDB").C("userStockList")

	user := User {
		Name: username,
	}


	stock := make(map[string]Stock)
	stock[symbol] = Stock{ StockName: stockName, Symbol: symbol}

	var existingUser CustomList
	if err = c.Find(bson.M{"user": user}).One(&existingUser); err != nil {
		//Not found
		log.Println("The user is not exist")
		customList := CustomList {
			User: user,
			Stock: stock,
		}

		if err = c.Insert(&customList); err != nil {
			log.Fatal(err)
		}
	}else{
		log.Println("Existing User")
		existingUser.AddStock(stock)
		if err := c.Update(bson.M{"_id": existingUser.Id}, existingUser); err != nil {
			log.Panic(err)
		}
	}
}

var baseTemplateURL = "src/github.com/go_practice"
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fileType := vars["type"]
	fileName := vars["fileName"]
	path := fmt.Sprintf("%s/%s/%s", baseTemplateURL, fileType, fileName)

	log.Printf("Loading File: %s", path)

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


