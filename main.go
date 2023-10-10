package main

import (
	"log"
)

var LatestOekaki Oekaki

func main() {
	LatestOekaki = Oekaki{
		Answer: "",
		Image:  "",
	}

	InitDB()
	log.Fatal(Handler().Listen(":3000"))
}
