package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

// Item represent a simple data structure with a title and body
type Item struct {
	Title string
	Body  string
}

// API defines the RPC methods for item management
type API int

var (
	database []Item
	mu       sync.Mutex // Mutex to protect the database slice
)

// ListItems fetches the entire database
func (a *API) ListItems(empty string, reply *[]Item) error {
	mu.Lock()
	defer mu.Unlock()

	*reply = database
	return nil
}

// GetItemByTitle fetches an item by its title
func (a *API) GetItemByTitle(title string, reply *Item) error {
	mu.Lock()
	defer mu.Unlock()

	for _, val := range database {
		if val.Title == title {
			*reply = val
		}
	}

	return errors.New("item not found")
}

// CreateItem adds a new item to the database
func (a *API) CreateItem(item Item, reply *Item) error {
	mu.Lock()
	defer mu.Unlock()

	// Check for duplicate titles
	for _, val := range database {
		if val.Title == item.Title {
			return errors.New("item with this title already exists")
		}
	}

	database = append(database, item)
	*reply = item
	return nil
}

// UpdateItem edits an existing item based on its title
func (a *API) EditItem(item Item, reply *Item) error {
	mu.Lock()
	defer mu.Unlock()

	for idx, val := range database {
		if val.Title == item.Title {
			database[idx] = item
			*reply = item
			return nil
		}
	}

	return errors.New("item not found")
}

// RemoveItem deletes an item from the database
func (a *API) DeleteItem(item Item, reply *Item) error {
	mu.Lock()
	defer mu.Unlock()

	for idx, val := range database {
		if val.Title == item.Title && val.Body == item.Body {
			database = append(database[:idx], database[idx+1:]...)
			*reply = item
			return nil
		}
	}

	return errors.New("item not found")
}

func main() {
	api := new(API)
	err := rpc.Register(api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")

	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("serving rpc on port %d", 4040)
	http.Serve(listener, nil)
}
