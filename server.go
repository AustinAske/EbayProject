package main

import (
// 	"fmt"
// 	"io"
// 	"os"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"


	"net/http"
//	"html/template"
)

type AuctionItem struct{
	Name			string
	StartingBid		float32
	Description 	string
}

var Db sql.DB
var connectionString = "Austin:@tcp(localhost:3306)/ebay_store"


var Templates = template.Must(template.ParseFiles("shop.html", "auctionAdded.html", "templates/shopItem.html"))


// servers static pages in file structure
func home(writer http.ResponseWriter, request *http.Request) {
        http.ServeFile(writer, request, request.URL.Path[1:])    
}

func shop(writer http.ResponseWriter, request *http.Request) {
// 	http.ServeFile(writer, request, request.URL.Path[1:] + ".html")	

		// 	Open Database connection
	Db, err := sql.Open("mysql", connectionString)
	if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer Db.Close()
	
	results, err := Db.Query("SELECT name,format(starting_bid,2),description FROM Auctions Limit 1;") // where close time is < current time
	if err != nil{
        http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	defer results.Close()
	
	var auctions AuctionItem

	for results.Next(){
		results.Scan(&auctions.Name, &auctions.StartingBid, &auctions.Description) // get rows from query
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

func main() {
	
	
	http.HandleFunc("/", home) // respond to any file path
	http.HandleFunc("/post", post) // respond to any file path
	http.HandleFunc("/shop", shop)
	http.HandleFunc("/submit", addAuction)
	http.ListenAndServe(":8000", nil)
}