package identity

import (
	"appengine/datastore"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"distudio.com/mage/model"
	"distudio.com/page"
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
	salt           = "AnticmS"
)

type User struct {
	model.Model `json:"-"`
	//Resource
	Name       string
	Surname    string
	Email      string
	Password   string
	Token      string
	Locale     string
	Permission Permission
	LastLogin  time.Time
	gUser      *guser.User `model:"-",json:"-"`
}

func (user *User) UnmarshalJSON(data []byte) error {
	// username (alias StringID) must be handled by the consumer of the model
	alias := struct {
		Name        string   `json:"name"`
		Surname     string   `json:"surname"`
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
	for permission, description := range Permissions {
		if user.HasPermission(permission) {
			perms = append(perms, description)
		}
	}
	return perms
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

// populates the user struct.
func (user *User) Create(ctx context.Context) error {
	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionCreateUser) {
		return page.NewPermissionError(PermissionName(PermissionCreateUser))
	}

	username := SanitizeUserName(user.Username())

	uf := page.NewRawField("username", true, username)
	uf.AddValidator(page.DatastoreKeyNameValidator{})

	// validate the username. Accepted values for the username are implementation dependent
	if err := uf.Validate(); err != nil {
		msg := fmt.Sprintf("invalid username %s", user.Username())
		return page.NewFieldError("username", errors.New(msg))
	}

	pf := page.NewRawField("password", true, user.Password)
	pf.AddValidator(page.LenValidator{MinLen: 8})

	if err := pf.Validate(); err != nil {
		msg := fmt.Sprintf("invalid password %s for username %s", user.Password, username)
		return page.NewFieldError("password", errors.New(msg))
	}

	if !current.HasPermission(PermissionEditPermissions) {
		// user without the EditPermission perm can only enable or disable a user
		if !((len(user.Permissions()) == 1 && user.IsEnabled()) || (len(user.Permissions()) == 0 && !user.IsEnabled())) {
			return page.NewPermissionError(PermissionName(PermissionEditPermissions))
		}
	}

	// check for user existence
	err := model.FromStringID(ctx, &User{}, username, nil)

	if err == nil {
		// user already exists
		msg := fmt.Sprintf("user %s already exists.", username)
		return page.NewFieldError("user", errors.New(msg))
	}

	if err != datastore.ErrNoSuchEntity {
		// generic datastore error
		return fmt.Errorf("error retrieving user with username %s: %s", username, err.Error())
	}

	user.Password = HashPassword(user.Password, salt)

	return nil
}


func (user *User) Update(ctx context.Context, res page.Resource) error {

	current, _ := ctx.Value(KeyUser).(User)
	if !current.HasPermission(PermissionEditUser) {
		return page.NewPermissionError(PermissionName(PermissionEditUser))
	}

	other := res.(*User)
	user.Name = other.Name

	if other.Password != "" {
		pf := page.NewRawField("password", true, other.Password)
		pf.AddValidator(page.LenValidator{MinLen: 8})

		if err := pf.Validate(); err != nil {
			msg := fmt.Sprintf("invalid password %s for username %s", other.Password, other.Username())
			return page.NewFieldError("user", errors.New(msg))
		}
		user.Password = other.Password
	}

	if other.Email != "" {
		ef := page.NewRawField("email", true, other.Email)
		if err := ef.Validate(); err != nil {
			msg := fmt.Sprintf("invalid email address: %s", other.Email)
			return page.NewFieldError("user", errors.New(msg))
		}
		user.Email = other.Email
	}

	if !current.HasPermission(PermissionEditPermissions) && other.ChangedPermission(*user) {
		return page.NewPermissionError(PermissionName(PermissionEditPermissions))
	}

	user.Name = other.Name
	user.Surname = other.Surname
	user.Permission = other.Permission

	return nil
}
