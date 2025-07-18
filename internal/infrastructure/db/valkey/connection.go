package valkey

import "github.com/valkey-io/valkey-go"

func NewConnection() (client valkey.Client, err error) {
	return valkey.NewClient(valkey.ClientOption{InitAddress: []string{"localhost:6379"}, Password: "valkey-password"})
}
