package identity

import (
	"crypto/sha1"
	"crypto/sha256"
	"decodica.com/flamel/model"
	"decodica.com/spellbook"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	guser "google.golang.org/appengine/user"
	"strings"
	"time"
	"unicode"
)

const (
	tokenSeparator = "|"
	hashLen        = 28
	UsernameMaxLen = 32
	UsernameMinLen = 4
)

type User struct {
	model.Model `json:"-"`
	//Resource
	Name    string
	Surname string
	//username    string `model:"-"`
	Email      string
	Password   string
	Token      string
	Locale     string
	Permission spellbook.Permission
	LastLogin  time.Time
	gUser      *guser.User `model:"-",json:"-"`
}

func (user *User) UnmarshalJSON(data []byte) error {
	// username (alias StringID) must be handled by the consumer of the model
	alias := struct {
		Name        string   `json:"name"`
		Surname     string   `json:"surname"`
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	user.Name = alias.Name
	user.Surname = alias.Surname
	user.Email = alias.Email
	//user.username = alias.Username
	user.GrantNamedPermissions(alias.Permissions)
	return nil
}

func (user *User) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name        string   `json:"name"`
		Surname     string   `json:"surname"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
	}

	return json.Marshal(&struct {
		Username string `json:"username"`
		Alias
	}{
		user.Username(),
		Alias{
			Name:        user.Name,
			Surname:     user.Surname,
			Email:       user.Email,
			Permissions: user.Permissions(),
		},
	})
}

func (user User) Permissions() []string {
	var perms []string
	for permission, description := range spellbook.Permissions {
		if user.HasPermission(permission) {
			perms = append(perms, description)
		}
	}
	return perms
}

func (user User) HasPermission(permission spellbook.Permission) bool {
	return user.Permission&permission != 0
}

// sanitizes a string to be used a username
// if there is an error or the username is invalid an empty string is returned
func SanitizeUserName(username string) string {
	u := strings.TrimSpace(username)

	if len(username) > UsernameMaxLen {
		return ""
	}

	if len(username) < UsernameMinLen {
		return ""
	}

	for _, c := range u {
		if unicode.IsLetter(c) || unicode.IsNumber(c) || c == '.' || c == '_' {
			continue
		}
		return ""
	}

	u = strings.ToLower(u)

	return u
}

func (user User) hash() string {
	now := time.Now().UTC().Unix()
	s := fmt.Sprintf("%s%s%s%s%d", user.StringID(), tokenSeparator, user.Password, tokenSeparator, now)
	hasher := sha1.New()
	hasher.Write([]byte(s))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

func (user User) IsGUser() bool {
	return user.gUser != nil
}

func (user User) Username() string {
	if user.IsGUser() {
		return user.gUser.String()
	}
	return user.StringID()
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
	return fmt.Sprintf("%s%s", hash, user.EncodedKey()), nil
}

/**
-- Resource implementation
*/

func (user *User) Id() string {
	return user.Username()
}

func (user *User) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, user)
	}
	return spellbook.NewUnsupportedError()
}

func (user *User) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(user)
	}
	return nil, spellbook.NewUnsupportedError()
}
