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
)

var templates map[string]*template.Template



func init(){
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index"] = template.Must(template.ParseFiles("src/github.com/go_practice/template/index.html"))

}


func main() {
	fmt.Println("Starting...")
	r := mux.NewRouter().StrictSlash(false)
	//mux.HandleFunc("/welcome", index)
	r.HandleFunc("/", index)
	r.HandleFunc("/search", getQuotesAndRender).Methods("Get")
	r.HandleFunc("/dbTest", dbTest).Methods("Get")
	log.Println("Listening...")

	server := &http.Server{
		Addr: ":8080",
		Handler: r,
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

	renderTemplate(w, "index", mainList.List.Resources)
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

func dbTest(w http.ResponseWriter, r *http.Request){
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("taskdb").C("quotes")

	symbol := r.URL.Query().Get("symbol")
	mainList, err := getQuotes(symbol)

	if err = c.Insert(&mainList); err != nil {
		log.Fatal(err)
	}


}
