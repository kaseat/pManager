package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/kaseat/pManager/storage"
	"golang.org/x/crypto/argon2"
)

// CheckСredentials checks credentials
func CheckСredentials(user, password string) (bool, error) {
	s := storage.GetStorage()
	hash, err := s.GetPassword(user)
	if err != nil {
		return false, err
	}
	if hash == "" {
		return false, errors.New("could not find user " + user)
	}
	return comparePassword(password, hash)
}

// SaveСredentials saves user/password
func SaveСredentials(user, password string) (bool, error) {
	config := &passwordConfig{
		time:    1,
		memory:  1024,
		threads: 2,
		keyLen:  32,
	}
	hash, err := generatePassword(config, password)
	if err != nil {
		return false, err
	}
	s := storage.GetStorage()
	err = s.SavePassword(user, hash)
	if err != nil {
		return false, err
	}
	return true, nil
}

type passwordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// generatePassword is used to generate a new password hash for storing and
// comparing at a later date.
func generatePassword(c *passwordConfig, password string) (string, error) {

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return full, nil
}

// comparePassword is used to compare a user-inputted password to a hash to see
// if the password matches or not.
func comparePassword(password, hash string) (bool, error) {

	parts := strings.Split(hash, "$")

	c := &passwordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.keyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}
