package contracts

import "SPBHistoryBot/clients"

type IClient interface {
	Updates(offset int, limit int) ([]clients.Update, error)
	SendMessage()
}
