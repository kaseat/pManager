package postgres

import (
	"errors"
	"strconv"

	"github.com/jackc/pgconn"
	"github.com/kaseat/pManager/models"
)

// AddPortfolio adds new potrfolio
func (db Db) AddPortfolio(userID string, p models.Portfolio) (string, error) {
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return "", errors.New("Invalid user Id format. Expected positive number")
	}
	var id int
	query := "insert into portfolios (pid,name,title) values ($1,$2,$3) returning id;"
	err = db.connection.QueryRow(db.context, query, uid, p.Name, p.Description).Scan(&id)
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if !ok {
			return "", err
		}
		if pgerr.Code == "23503" {
			return "", errors.New("could not add portfolio to unknown user")
		}
		return "", err
	}
	return strconv.Itoa(id), nil
}

// GetPortfolio gets operation by id
func (db Db) GetPortfolio(userID string, portfolioID string) (models.Portfolio, error) {
	result := models.Portfolio{}
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return result, errors.New("Invalid user Id format. Expected positive number")
	}
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return result, errors.New("Invalid portfolio Id format. Expected positive number")
	}

	var name string
	var title string

	query := "select name,title from portfolios where pid = $1 and id = $2;"
	err = db.connection.QueryRow(db.context, query, uid, pid).Scan(&name, &title)
	if err != nil {
		return result, err
	}
	result.UserID = userID
	result.PortfolioID = portfolioID
	result.Name = name
	result.Description = title
	return result, nil
}

// GetPortfolios gets all portfolio fpvie user Id
func (db Db) GetPortfolios(userID string) ([]models.Portfolio, error) {
	result := []models.Portfolio{}
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return nil, errors.New("Invalid user Id format. Expected positive number")
	}
	query := "select id,name,title from portfolios where pid = $1;"
	rows, err := db.connection.Query(db.context, query, uid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		var name string
		var title string
		err = rows.Scan(&id, &name, &title)
		if err != nil {
			return nil, err
		}
		p := models.Portfolio{
			PortfolioID: strconv.Itoa(id),
			UserID:      userID,
			Name:        name,
			Description: title,
		}
		result = append(result, p)
	}
	return result, nil
}

// UpdatePortfolio updates portfolio with provided uid and pid
func (db Db) UpdatePortfolio(userID string, portfolioID string, p models.Portfolio) (bool, error) {
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return false, errors.New("Invalid user Id format. Expected positive number")
	}
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return false, errors.New("Invalid portfolio Id format. Expected positive number")
	}

	query := "update portfolios set name = $1, title = $2 where pid = $3 and id = $4;"
	r, err := db.connection.Exec(db.context, query, p.Name, p.Description, uid, pid)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// DeletePortfolio removes portfolio by Id
// Also removes all operations associated with this portfolio
func (db Db) DeletePortfolio(userID string, portfolioID string) (bool, error) {
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return false, errors.New("Invalid user Id format. Expected positive number")
	}
	pid, err := strconv.ParseInt(portfolioID, 10, 32)
	if err != nil {
		return false, errors.New("Invalid portfolio Id format. Expected positive number")
	}

	query := "delete from portfolios where pid = $1 and id = $2;"
	r, err := db.connection.Exec(db.context, query, uid, pid)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// DeletePortfolios removes all portfolios for provided user
// Also removes all operations associated with this portfolios
func (db Db) DeletePortfolios(userID string) (int64, error) {
	// todo: add operations deletion
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return 0, errors.New("Invalid user Id format. Expected positive number")
	}
	c, err := db.connection.Begin(db.context)
	query := "delete from operations o using portfolios p where o.pid = p.id and p.pid = $1;"
	_, err = c.Exec(db.context, query, uid)
	if err != nil {
		c.Rollback(db.context)
		return 0, err
	}
	query = "delete from portfolios where pid = $1;"
	r, err := c.Exec(db.context, query, uid)
	if err != nil {
		c.Rollback(db.context)
		return 0, err
	}
	err = c.Commit(db.context)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected(), nil
}
