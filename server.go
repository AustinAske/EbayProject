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
	BuyPrice		float32
	Description 	string
	CloseTime		string
	CurrentTime 	string
	BidTime			string
	CurrentPrice	float32
	State			int
	CustomerId		int
}

var currenttime = "1"
var Db sql.DB
var connectionString = "Austin:@tcp(localhost:3306)/ebay_store"


var Templates = template.Must(template.ParseFiles("shop.html", "post.html", "auctionAdded.html", "history.html", "templates/shopItem.html", "templates/table.html"))


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
	if currenttime < request.PostFormValue("currenttime"){
		addItemStmt, err := Db.Prepare("INSERT INTO Bids(price, customer_id, auction_id, customer_name, timestamp) VALUES( ?, ?, ?, ?, ?)")
		if err!= nil {
			panic(err.Error())
		}
		
		
		price := request.PostFormValue("bid_amount")
		customer_id := 1
		auction_id := request.PostFormValue("auction_id")
		customer_name := request.PostFormValue("bid_uname")
		
		_, err = addItemStmt.Exec(price, customer_id, auction_id, customer_name, currenttime)
		if err != nil{
			panic(err.Error())
		}
	}
		http.Redirect(writer, request, "/shop", 302)
}



func setTime(writer http.ResponseWriter, request *http.Request){
	request.ParseForm()

	currenttime = request.PostFormValue("currentTime")
	returnAddress := request.Referer()
	fmt.Printf("Current Time changed to %s\n", currenttime)
	http.Redirect(writer, request, returnAddress, 302)
}

func shop(writer http.ResponseWriter, request *http.Request) {	
		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer Db.Close()
	
	results, err := Db.Query("Select Auctions.id, name, price, description, buy_price, close_time from (select * from Bids where timestamp <= ? order by price DESC)as max join Auctions on Auctions.id = max.`auction_id` where state = 0 group by `auction_id`;", currenttime) // where close time is < current time
	if err != nil{
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer results.Close()
	
	auctions := []AuctionItem{}

	for results.Next(){
		var newResult AuctionItem
		newResult.CurrentTime = currenttime
		results.Scan(&newResult.Id, &newResult.Name, &newResult.CurrentPrice, &newResult.Description, &newResult.BuyPrice, &newResult.CloseTime) // get rows from query		
		auctions = append(auctions, newResult)
	}


	Templates.ExecuteTemplate(writer, "shop" ,auctions)

}

func post(writer http.ResponseWriter, request *http.Request) {
	Templates.ExecuteTemplate(writer, "post", currenttime)
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
	addItemStmt, err := Db.Prepare("INSERT INTO Auctions(name, category, description, buy_price, close_time) VALUES( ?, ?, ?, ?, ?)")
	if err!= nil {
		panic(err.Error())
	}
	
	itemName := request.PostFormValue("item_name")
	category := request.PostFormValue("item_category")
	description := request.PostFormValue("item_description")
	endTime := request.PostFormValue("close_time")
	buyPrice := request.PostFormValue("buy_price")
	
// 	auctions := AuctionItem{Name: itemName}

	_, err = addItemStmt.Exec(itemName, category, description, buyPrice, endTime)
	if err != nil{
		panic(err.Error())
	}
	http.Redirect(writer, request, "/shop", 302)

// 	Templates.ExecuteTemplate(writer, "", auctions) 
		
	
}

// testing state
func history(writer http.ResponseWriter, request *http.Request) {
		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer Db.Close()
	
	results, err := Db.Query("Select Auctions.id, name, price, description, buy_price, timestamp, close_time, Auctions.state, customer_id from (select * from Bids where timestamp <= ?) as max join Auctions on Auctions.id = max.`auction_id` where customer_id is not NULL order by name, price DESC, timestamp DESC ;", currenttime) // where close time is < current time
	if err != nil{
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer results.Close()
	
	auctions := []AuctionItem{}

// 	i := 0
	for results.Next(){
		var newResult AuctionItem
		newResult.CurrentTime = currenttime
		results.Scan(&newResult.Id, &newResult.Name, &newResult.CurrentPrice, &newResult.Description, &newResult.BuyPrice, &newResult.BidTime, &newResult.CloseTime, &newResult.State, &newResult.CustomerId) // get rows from query		
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
	go http.HandleFunc("/setTime", setTime)
	http.HandleFunc("/history", history)
	http.HandleFunc("/submit", addAuction)
	http.ListenAndServe(":8000", nil)
}