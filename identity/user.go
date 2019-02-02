package identity

import (
	"crypto/sha1"
	"crypto/sha256"
	"distudio.com/mage/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	tokenSeparator = "|"
	hashLen = 28
)

type User struct {
	model.Model
	Name        string
	Surname     string
	Email       string
	Password    string
	Token       string
	Locale      string
	Permission Permission
	LastLogin time.Time
}

func (user *User) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name    string    `json:"name"`
		Surname string    `json:"surname"`
		Email   string    `json:"email"`
		Token   string    `json:"token"`
		Permission   Permission `json:"level"`
	}

	return json.Marshal(&struct {
		Id string `json:"id"`
		Alias
	}{
		user.Key.Encode(),
		Alias{
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
			Token:   user.Token,
			Permission:  user.Permission,
		},
	})
}

func (user User) hash() string {
	now := time.Now().UTC().Unix()
	s := fmt.Sprintf("%s%s%s%s%d",user.StringID(), tokenSeparator, user.Password, tokenSeparator, now)
	hasher := sha1.New()
	hasher.Write([]byte(s))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

func HashPassword(password string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	if salt != "" {
		hasher.Write([]byte(salt))
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (user User) GenerateToken() (string, error) {
	if user.Key == nil {
		return "", errors.New("can't generate token. User does not exists")
	}
	hash := user.hash()
	return fmt.Sprintf("%s%s",hash, user.EncodedKey()), nil
}
