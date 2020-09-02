package postgres

import "fmt"

// AddTcsToken adds token to access tcs API
func (db Db) AddTcsToken(token string) error {
	query := fmt.Sprintf("update settings set settings = settings::jsonb - 'tcs_token' || jsonb_build_object('tcs_token', '%s') where settings->>'ver' = '1';", token)
	_, err := db.connection.Exec(db.context, query)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTcsToken deletes token to access tcs API
func (db Db) DeleteTcsToken() error {
	query := "update settings set settings = settings::jsonb - 'tcs_token' where settings->>'ver' = '1';"
	_, err := db.connection.Exec(db.context, query)
	if err != nil {
		return err
	}
	return nil
}

// GetTcsToken finds token to access tcs API
func (db Db) GetTcsToken() (string, error) {
	var out *string
	query := "select settings->>'tcs_token' as token from settings where settings->>'ver' = '1';"
	err := db.connection.QueryRow(db.context, query).Scan(&out)
	if err != nil {
		return "", err
	}
	if out != nil {
		return *out, nil
	}
	return "", nil
}
