package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/kaseat/pManager/models/provider"
)

// AddUserLastUpdateTime saves last date when specified provider made sync
func (db Db) AddUserLastUpdateTime(login string, provider provider.Type, date time.Time) error {
	pvds := getSyncPviderTypeByName()
	query := "insert into user_sync (uid,provider_id,last_sync) select u.id, $2,$3 from users u where login = $1 on conflict on constraint pk_user_sync do update set last_sync = $3;"
	r, err := db.connection.Exec(db.context, query, login, pvds[provider], date)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23503" {
				return fmt.Errorf("last uptate time with login: %s and provider: %s already exists", login, provider)
			}
		}
		return err
	}
	if r.RowsAffected() == 0 {
		return fmt.Errorf("could not fint user with login: %s", login)
	}
	return nil
}

// GetUserLastUpdateTime receives last date when specified provider made sync
func (db Db) GetUserLastUpdateTime(login string, provider provider.Type) (time.Time, error) {
	var result *time.Time
	pvds := getSyncPviderTypeByName()
	query := "select last_sync from user_sync us inner join users u on u.id = us.uid where u.login = $1 and us.provider_id = $2;"
	err := db.connection.QueryRow(db.context, query, login, pvds[provider]).Scan(&result)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	if result != nil {
		return *result, nil
	}
	return time.Time{}, nil
}

// DeleteUserLastUpdateTime removes last date when specified provider made sync
func (db Db) DeleteUserLastUpdateTime(login string, provider provider.Type) error {
	pvds := getSyncPviderTypeByName()
	query := "delete from user_sync us using users u where u.id = us.uid and u.login = $1 and us.provider_id = $2;"
	_, err := db.connection.Exec(db.context, query, login, pvds[provider])
	return err
}

func getSyncPviderTypeByName() map[provider.Type]int {
	return map[provider.Type]int{
		provider.Sber: 1,
		provider.Tcs:  2,
		provider.Vtb:  3,
	}
}
