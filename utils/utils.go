package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"os"
	"strconv"
	"strings"

	"github.com/jihoon6372/hog/config"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v2"
)

// HashPassword Generator
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash Check
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPasswordPbkdf2Sha256 : hashing the password using PBKDF2_SHA256
func HashPasswordPbkdf2Sha256(password string) (string, error) {
	randByte := make([]byte, 8)

	_, err := rand.Read(randByte)
	if err != nil {
		return "", err
	}

	base64RandByte := base64.StdEncoding.EncodeToString(randByte)
	salt := []byte(base64RandByte)

	iter := 120000

	dk := pbkdf2.Key([]byte(password), salt, iter, 32, sha256.New)

	hashedPW := "pbkdf2_sha256$" + strconv.Itoa(iter) + "$" + string(salt) + "$" + base64.StdEncoding.EncodeToString(dk)

	return hashedPW, nil
}

// ComparePassword : compare the password
func ComparePassword(password string, hash string) bool {
	splitted := strings.Split(hash, "$")

	salt := []byte(splitted[2])

	// saved password iteration value should be converted to int
	iter, _ := strconv.Atoi(splitted[1])

	dk := pbkdf2.Key([]byte(password), salt, iter, 32, sha256.New)

	hashedPW := "pbkdf2_sha256$" + splitted[1] + "$" + splitted[2] + "$" + base64.StdEncoding.EncodeToString(dk)

	if subtle.ConstantTimeCompare([]byte(hash), []byte(hashedPW)) == 0 {
		return false
	}

	return true
}

// ReadConfig ...
func ReadConfig(cfg *config.Config) {
	f, err := os.Open("./config/config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
}
