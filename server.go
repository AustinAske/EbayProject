package main

import (
// 	"fmt"
// 	"io"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"


	"net/http"
//	"html/template"
)

// servers static pages in file structure
func home(writer http.ResponseWriter, request *http.Request) {
        http.ServeFile(writer, request, request.URL.Path[1:])    
}

func shop(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, request.URL.Path[1:] + ".html")
}

func post(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, request.URL.Path[1:] + ".html")
}

func addAuction(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm() // parse form data into json element
// 	fmt.Fprintf(writer, "Post request sent to server\n%q",request)
	err = db.Ping()
	if err != nil {
		    panic(err.Error()) // proper error handling instead of panic in your app
		}	
	
	http.Redirect(writer, request, "/shop", 302)
}


/*
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
	    panic(err.Error()) // proper error handling instead of panic in your app
	}
	
// 	Test the database conneciton by sending data to it!! WORKED!!!
	result, err := db.Exec("INSERT INTO `Customers` (`id`, `name`, `email`) VALUES(NULL, 'Jazzy John', 'Jazzy@google.com');")	
	if err != nil {
		panic(err.Error())
	}
*/
		

func main() {

 	db, err := sql.Open("mysql", "Austin:@tcp(localhost:3306)/ebay_store")
    	if err != nil {
        	panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    	}
		defer db.Close()
	
	http.HandleFunc("/", home) // respond to any file path
	http.HandleFunc("/post", post) // respond to any file path
	http.HandleFunc("/shop", shop)
	http.HandleFunc("/submit", addAuction)
	http.ListenAndServe(":8000", nil)
}