package main

import (
	"fmt"
	"github.com/kaseat/pManager/portfolio"
)

func main() {
	cfg := portfolio.Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "tcs",
	}
	portfolio.Init(cfg)
	managePortfolio()
}

func managePortfolio() {
	// create portfolio
	p, err := portfolio.AddPortfolio("Awesome portfolio", "Portfolio for my awesome investments")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	fmt.Println(p)

	// edit portfolio
	p.Name = "Good one"

	if ok, err := p.UpdatePortfolio(); ok {
		fmt.Println("Updated successfully!")
	} else if err == nil {
		fmt.Println("Nothing updated!")
		return
	} else {
		fmt.Println("Something went wrong:", err)
		return
	}

	// get all portfolios
	var ps []portfolio.Portfolio
	ps, err = portfolio.GetAllPortfolios()

	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	for _, p := range ps {
		println(p.String())
	}

	// delete all portfolios
	_, err = portfolio.DeletePortfolio(ps[0].PortfolioID)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	// delete all portfolios
	_, err = portfolio.DeleteAllPortfolios()
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	fmt.Println("Delete all done!")
}
