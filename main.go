package main

import (
	"log"
)

func main() {
	log.Println("SpiceDB Go Library")
	log.Println("================")
	log.Println("")
	log.Println("This is the SpiceDB Go library package.")
	log.Println("")
	log.Println("For examples, see:")
	log.Println("- examples/basic/main.go - Basic usage examples")
	log.Println("- examples/http-service/main.go - HTTP service example")
	log.Println("")
	log.Println("To use in your project:")
	log.Println("  go get github.com/spicedb/spicedb-go")
	log.Println("")
	log.Println("Quick start:")
	log.Println(`
	import "github.com/spicedb/spicedb-go/spicedb"
	
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
`)
	log.Println("")
	log.Println("See spicedb/README.md for full documentation")
}
