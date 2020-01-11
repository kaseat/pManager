package main

import (
	// "bufio"
	// "encoding/json"
	"fmt"
	"time"
	// "os"
	// "time"

	// "github.com/kaseat/tcssync"
	"github.com/kaseat/tcssync/portfolio"
)

func main() {
	cfg := portfolio.Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "tcs",
	}
	portfolio.Init(cfg)
	manageOperations()
}

func manageOperations() {
	var p portfolio.Portfolio

	// create portfolio if doesn't exists
	ps, err := portfolio.GetAllPortfolios()

	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	if len(ps) != 1 {
		err = portfolio.DeleteAllPortfolios()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

		p, err = portfolio.AddPortfolio("Awesome portfolio", "My awesome portfolio")
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
	} else {
		p = ps[0]
	}

	// add pay in operation if there's no funds
	op := portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         1,
		Quantity:      10000,
		FIGI:          "RUB",
		DateTime:      time.Now(),
		OperationType: portfolio.PayIn}

	if bal, err := p.GetBalanceByCurrency(portfolio.RUB); err != nil {
		fmt.Println("Something went wrong:", err)
		return
	} else if bal == 0 {
		opID, err := p.AddOperation(op)
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
		fmt.Println("Successfully added opeation with id:", opID)
	} else {
		fmt.Println("Current balance:", bal)
	}

	// add buy operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         679,
		Quantity:      3,
		FIGI:          "BBG005DXDPK9",
		DateTime:      time.Now().AddDate(0, 0, 1),
		OperationType: portfolio.Buy}

	var opID string
	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// add another buy operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         620,
		Quantity:      5,
		FIGI:          "BBG00NRFC2X2",
		DateTime:      time.Now().AddDate(0, 0, 1),
		OperationType: portfolio.Buy}

	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// add sell operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         625,
		Quantity:      2,
		FIGI:          "BBG00NRFC2X2",
		DateTime:      time.Now().AddDate(0, 0, 2),
		OperationType: portfolio.Sell}

	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// add another sell operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         690,
		Quantity:      1,
		FIGI:          "BBG005DXDPK9",
		DateTime:      time.Now().AddDate(0, 0, 2),
		OperationType: portfolio.Sell}

	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// now get all operations

	var ops []portfolio.Operation
	ops, err = p.GetAllOperations()
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	fmt.Println("Successfully fetched ", len(ops), "operations from", p.Name, "portfolio:")
	for _, op := range ops {
		fmt.Println(op)
	}
}
