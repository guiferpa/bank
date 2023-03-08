package account

import "time"

type Account struct {
	ID             uint
	DocumentNumber string
	CreatedAt      time.Time
}
