package postgres

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/kaseat/pManager/models"
	"golang.org/x/oauth2"
)

// AddUser saves user and password hash to storage
func (db Db) AddUser(login, email, hash string) (string, error) {
	var id int
	var err error
	if email == "" {
		query := "insert into users (login,hash,role_id) values ($1,$2,2) returning id;"
		err = db.connection.QueryRow(db.context, query, login, hash).Scan(&id)
	} else {
		query := "insert into users (login,hash,role_id,email) values ($1,$2,2,$3) returning id;"
		err = db.connection.QueryRow(db.context, query, login, hash, email).Scan(&id)
	}
	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "pk_users_login" (SQLSTATE 23505)` {
			return "", errors.New("User with this login already exists")
		}
		return "", err
	}

	return strconv.Itoa(id), nil
}

// GetUserByLogin gets user by login
func (db Db) GetUserByLogin(login string) (models.User, error) {
	result := models.User{}
	var id int
	var role int
	var email *string

	query := "select id,role_id,email from users where login = $1;"
	err := db.connection.QueryRow(db.context, query, login).Scan(&id, &role, &email)
	if err != nil {
		return result, err
	}

	if email != nil {
		result.Email = *email
	}

	result.Login = login
	result.UserID = strconv.Itoa(id)
	if role == 1 {
		result.IsAdmin = true
	} else {
		result.IsAdmin = false
	}

	return result, nil
}

// AddUserState adds state to user
func (db Db) AddUserState(login string, state string) error {
	query := "update users set g_sync_state = $1 where login = $2;"
	_, err := db.connection.Exec(db.context, query, state, login)
	if err != nil {
		return err
	}
	return nil
}

// GetUserState gets user's state
func (db Db) GetUserState(login string) (string, error) {
	state := ""
	query := "select g_sync_state from users where login = $1;"
	err := db.connection.QueryRow(db.context, query, login).Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}

// AddUserToken adds oauth2 token to user
func (db Db) AddUserToken(state string, token *oauth2.Token) error {
	bytes, err := json.Marshal(*token)
	if err != nil {
		return err
	}
	query := "update users set g_sync_token = $1 where g_sync_state = $2;"
	_, err = db.connection.Exec(db.context, query, bytes, state)
	if err != nil {
		return err
	}
	return nil
}

// GetUserToken gets user's oauth2 token
func (db Db) GetUserToken(login string) (oauth2.Token, error) {
	token := oauth2.Token{}
	query := "select g_sync_token from users where login = $1;"
	err := db.connection.QueryRow(db.context, query, login).Scan(&token)
	if err != nil {
		return token, err
	}
	return token, nil
}

// GetUserPassword gets password hash from storage
func (db Db) GetUserPassword(login string) (string, error) {
	hash := ""
	query := "select hash from users where login = $1;"
	err := db.connection.QueryRow(db.context, query, login).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// UpdateUserPassword updates user password
func (db Db) UpdateUserPassword(login, hash string) (bool, error) {
	query := "update users set hash = $1 where login = $2;"
	r, err := db.connection.Exec(db.context, query, hash, login)
	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

// DeleteUser removes password hash from storage
// Also removes all portfolios associated with this user
// Also removes all operations associated with portfolios of this user
func (db Db) DeleteUser(login string) (bool, error) {
	c, err := db.connection.Begin(db.context)
	if err != nil {
		return false, err
	}
	query := "delete from operations o using portfolios p, users u where u.id = p.uid and o.pid =p.id and u.login = $1;"
	_, err = c.Exec(db.context, query, login)
	if err != nil {
		c.Rollback(db.context)
		return false, err
	}
	query = "delete from portfolios p using users u where u.id = p.uid and u.login = $1;"
	_, err = c.Exec(db.context, query, login)
	if err != nil {
		c.Rollback(db.context)
		return false, err
	}
	query = "delete from user_sync s using users u where u.id = s.uid and u.login = $1;"
	_, err = c.Exec(db.context, query, login)
	if err != nil {
		c.Rollback(db.context)
		return false, err
	}
	query = "delete from users where login = $1;"
	r, err := c.Exec(db.context, query, login)
	if err != nil {
		c.Rollback(db.context)
		return false, err
	}
	err = c.Commit(db.context)

	if err != nil {
		return false, err
	}
	if r.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}
