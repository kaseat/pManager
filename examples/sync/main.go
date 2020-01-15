package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kaseat/tcssync/sync"
)

func main() {
	cfg := sync.Config{
		MongoURL:   "mongodb://localhost:27017",
		DbName:     "tcs",
		TcsToken:   "SANDBOX TOKEN",
		TcsTimeout: 5,
	}
	sync.Init(cfg)
	manageSync()
}

func manageSync() {
	sync.Price()

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
