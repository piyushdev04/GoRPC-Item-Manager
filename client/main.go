package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// Item represents a simple data structure with a title and body
type Item struct {
	Title string
	Body  string
}

func main() {
	var reply Item
	var db []Item

	// Connect to the RPC server
	client, err := rpc.DialHTTP("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	// Create items to add to the database
	a := Item{"First", "A first item"}
	b := Item{"Second", "A second item"}
	c := Item{"Third", "A third item"}

	// Add items to the database
	client.Call("API.CreateItem", a, &reply)
	client.Call("API.CreateItem", b, &reply)
	client.Call("API.CreateItem", c, &reply)

	// List all items in the database
	client.Call("API.ListItems", "", &db)
	fmt.Println("Database:", db)

	// Update an item in the database
	client.Call("API.UpdateItem", Item{"Second", "A new second item"}, &reply)

	// Remove an item from the database
	client.Call("API.RemoveItem", c, &reply)

	// List all items in the database after removal
	client.Call("API.ListItems", "", &db)
	fmt.Println("Database: ", db)

	// Get a specific item by title
	client.Call("API.GetItemByTitle", "First", &reply)
	fmt.Println("First item: ", reply)
}
