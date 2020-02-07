package main

import (
	"fmt"
	"github.com/kaseat/pManager/portfolio"
	"time"
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

	if len(ps) != 0 {
		for _, p := range ps {
			n, err := p.DeleteAllOperations()
			if err != nil {
				fmt.Println("Something went wrong:", err)
				return
			}
			fmt.Println("Successfully removed", n, "operations from", p.Name)
		}
		_, err = portfolio.DeleteAllPortfolios()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
		fmt.Println("Successfully removed all portfolios")

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
		Volume:        10000,
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
		Volume:        3,
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

	// add buy operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         620,
		Volume:        1,
		FIGI:          "BBG00NRFC2X2",
		DateTime:      time.Now().AddDate(0, 0, 1),
		OperationType: portfolio.Buy}

	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// add buy operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         622,
		Volume:        5,
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
		Volume:        2,
		FIGI:          "BBG00NRFC2X2",
		DateTime:      time.Now().AddDate(0, 0, 2),
		OperationType: portfolio.Sell}

	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// add sell operation
	op = portfolio.Operation{
		Currency:      portfolio.RUB,
		Price:         690,
		Volume:        1,
		FIGI:          "BBG005DXDPK9",
		DateTime:      time.Now().AddDate(0, 0, 2),
		OperationType: portfolio.Sell}

	opID, err = p.AddOperation(op)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully added opeation with id:", opID)

	// get all operations
	var ops []portfolio.Operation
	ops, err = p.GetAllOperations()
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully fetched", len(ops), "operations from", p.Name)
	for _, op := range ops {
		fmt.Println(op)
	}

	// get all operations with figi BBG005DXDPK9
	ops, err = p.GetAllOperationsByFigi("BBG005DXDPK9")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully fetched", len(ops), "operations from", p.Name)
	for _, op := range ops {
		fmt.Println(op)
	}

	// get actual RUB balance
	bal, err := p.GetBalanceByCurrency(portfolio.RUB)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Current ballance is", bal, portfolio.RUB)

	// get RUB balance on date
	today := time.Now().UTC().AddDate(0, 0, 1)
	bal, err = p.GetBalanceByCurrencyTillDate(portfolio.RUB, today)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Ballance on", today, "is", bal, portfolio.RUB)

	// get actual BBG005DXDPK9 balance
	bal, err = p.GetBalanceByFigi("BBG005DXDPK9")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Current ballance for BBG005DXDPK9 is", bal, portfolio.RUB)

	// get actual BBG005DXDPK9 average price
	bal, err = p.GetAveragePriceByFigi("BBG00NRFC2X2")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Current average price for BBG005DXDPK9 is", bal, portfolio.RUB)

	// get BBG00NRFC2X2 balance on date
	bal, err = p.GetAveragePriceByFigiTillDate("BBG00NRFC2X2", today)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Average price for BBG00NRFC2X2 on", today, "is", bal, portfolio.RUB)

	// get BBG005DXDPK9 balance on date
	bal, err = p.GetBalanceByFigiTillDate("BBG005DXDPK9", today)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Ballance for BBG005DXDPK9 on", today, "is", bal, portfolio.RUB)

	// get operation by id
	op, err = p.GetOperationByID(ops[0].OperationID)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Successfully fetched", op.OperationID, "operation from", p.Name)
}
