package newsletter

import (
	"distudio.com/mage/model"
	"time"
)

var ZeroTime = time.Time{}

type Newsletter struct {
	model.Model
	Email string `json:"email"`
}
