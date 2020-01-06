package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/kaseat/tcssync"
)

func main() {
	tcssync.Init("config.json")
	fmt.Println("Starting main...")
	operationManagemrnt()
	sync()
}

func sync() {
	tcssync.SyncPrice()
	tcssync.SyncPriceLastDay()

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func operationManagemrnt() {
	fmt.Println("Starting operationManagemrnt...")
	op := tcssync.Operation{
		Currency:      tcssync.RUB,
		Price:         3574,
		Quantity:      1,
		FIGI:          "BBG005HLSZ23",
		DateTime:      time.Date(2020, time.January, 03, 11, 22, 33, 113, time.UTC),
		OperationType: tcssync.Buy}

	fmt.Println("Inserting operation...")
	opID := tcssync.AddOperation(op)

	fmt.Println("Getting operation by ID:", opID)
	op, _ = tcssync.GetOperation(opID)
	fmt.Println("GetOperation result:", op)

	fmt.Println("Getting all operations:")
	res, _ := tcssync.GetAllOperations()
	fmt.Println("Operations:", res)

	fmt.Println("Getting RUB balance:")
	b := tcssync.GetBalance(tcssync.RUB)
	fmt.Println("Balance:", b)
}
