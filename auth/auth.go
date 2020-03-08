package auth

import "errors"

// CheckСredentials checks credentials
func CheckСredentials(user, password string) (bool, error) {
	if user == "admin" && password == "password" {
		return true, nil
	}
	return false, errors.New("Invalid credentials")
}
