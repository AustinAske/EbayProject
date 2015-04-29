package main

import (
	"fmt"
// 	"io"
// 	"web"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"


	"net/http"
//	"html/template"
)

type AuctionItem struct{
	Id				int
	Name			string
	StartingBid		float32
	Description 	string
}

var Db sql.DB
var connectionString = "Austin:@tcp(localhost:3306)/ebay_store"


var Templates = template.Must(template.ParseFiles("shop.html", "auctionAdded.html", "history.html", "templates/shopItem.html", "templates/table.html"))


// servers static pages in file structure
func home(writer http.ResponseWriter, request *http.Request) {
        http.ServeFile(writer, request, request.URL.Path[1:])    
}

func updateBid( writer http.ResponseWriter, request *http.Request){
	request.ParseForm()
	
		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer Db.Close()
	
	// 	create prepare statement
	addItemStmt, err := Db.Prepare("INSERT INTO Bids(price, customer_id, auction_id, customer_name, timestamp) VALUES( ?, ?, ?, ?, ?)")
	if err!= nil {
		panic(err.Error())
	}
	
	price := request.PostFormValue("bid_amount")
	customer_id := 1
	auction_id := request.PostFormValue("auction_id")
// 	endTime := request.PostFormValue("end_time")
	customer_name := request.PostFormValue("bid_uname")
	currentTime := request.PostFormValue("currentTime")
	
	_, err = addItemStmt.Exec(price, customer_id, auction_id, customer_name, currentTime)
	if err != nil{
		panic(err.Error())
	}
		fmt.Printf("%q", currentTime)
		http.Redirect(writer, request, "/shop", 302)

	

}

func shop(writer http.ResponseWriter, request *http.Request) {
// 	http.ServeFile(writer, request, request.URL.Path[1:] + ".html")	

		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer Db.Close()
	
	results, err := Db.Query("SELECT id, name,format(starting_bid,2),description FROM Auctions;") // where close time is < current time
	if err != nil{
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer results.Close()
	
	auctions := []AuctionItem{}

	for results.Next(){
		var newResult AuctionItem
		results.Scan(&newResult.Id, &newResult.Name, &newResult.StartingBid, &newResult.Description) // get rows from query		
		auctions = append(auctions, newResult)
	}


	Templates.ExecuteTemplate(writer, "shop" ,auctions)

}

func post(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, request.URL.Path[1:] + ".html")
}

func addAuction(writer http.ResponseWriter, request *http.Request) {
	// parse form data into json element
	request.ParseForm()
	
		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer Db.Close()
	
	// 	create prepare statement
	addItemStmt, err := Db.Prepare("INSERT INTO Auctions(name, category, description, starting_bid, open_time) VALUES( ?, ?, ?, ?, ?)")
	if err!= nil {
		panic(err.Error())
	}
	
	itemName := request.PostFormValue("item_name")
	category := request.PostFormValue("item_category")
	description := request.PostFormValue("item_description")
// 	endTime := request.PostFormValue("end_time")
	startingBid := request.PostFormValue("starting_bid")
	openTime := request.PostFormValue("currentTime")
	
	auctions := AuctionItem{Name: itemName}

	_, err = addItemStmt.Exec(itemName, category, description, startingBid, openTime)
	if err != nil{
		panic(err.Error())
	}
	
	Templates.ExecuteTemplate(writer, "auctionAdded", auctions) 
// 	http.Redirect(writer, request, "/shop", 302)
		
// 		fmt.Fprintf(writer, "Post request sent to server\n%q",request)	
	
}

// testing state
func history(writer http.ResponseWriter, request *http.Request) {
		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer Db.Close()
	
	results, err := Db.Query("SELECT name,format(starting_bid,2),description FROM Auctions Limit 3;") // where close time is < current time
	if err != nil{
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer results.Close()
	
	auctions := []AuctionItem{}

// 	i := 0
	for results.Next(){
		var newResult AuctionItem
		results.Scan(&newResult.Name, &newResult.StartingBid, &newResult.Description) // get rows from query		
		auctions = append(auctions, newResult)

// 		fmt.Printf("%s\n", auctions[i].Name)
// 		i++
	}

	Templates.ExecuteTemplate(writer, "history", auctions)
}


func main() {
	
	
	http.HandleFunc("/", home) // respond to any file path
	http.HandleFunc("/post", post) // respond to any file path
	http.HandleFunc("/shop", shop)
	http.HandleFunc("/updateBid", updateBid)
	http.HandleFunc("/history", history)
	http.HandleFunc("/submit", addAuction)
	http.ListenAndServe(":8000", nil)
}