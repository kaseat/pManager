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
	instrumentBalanceInfo()
	operationManagemrnt()
	sync()
}

func portfolioBalanceInfo() {
	currency := tcssync.RUB
	tillDate, _ := time.Parse("20060102", "20200104")
	balance, _ := tcssync.GetBalanceByCurrency(currency)
	fmt.Println("Current balance is", balance, currency)
	balance, _ = tcssync.GetBalanceByCurrencyTillDate(currency, tillDate)
	fmt.Println("Balance on", tillDate.Format("2006-01-02"), "is", balance, currency)
}

func instrumentBalanceInfo() {
	figi := "BBG005HLSZ23"
	tillDate, _ := time.Parse("20060102", "20200104")
	balance, _ := tcssync.GetBalanceByFigi(figi)
	fmt.Println("Current balance of", figi, "is", balance)
	balance, _ = tcssync.GetBalanceByFigiTillDate(figi, tillDate)
	fmt.Println("Balance on", tillDate.Format("2006-01-02"), "of", figi, "is", balance)
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
		DateTime:      time.Now(),
		OperationType: tcssync.Sell}

	fmt.Println("Inserting operation...")

	tcssync.AddOperation(op)

	opID, _ := tcssync.AddOperation(op)

	fmt.Println("Getting operation by ID:", opID)
	op, _ = tcssync.GetOperationByID(opID)
	fmt.Println("GetOperation result:", op)

	fmt.Println("Getting all operations:")
	res, _ := tcssync.GetAllOperations()
	fmt.Println("Operations:", res)
}
