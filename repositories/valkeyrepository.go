package repositories

import (
	"context"

	"github.com/valkey-io/valkey-go"
)

type ValkeyRepository struct {
	client   valkey.Client
	address  string
	username string
	password string
}

func NewValkeyRepository(username, password, addr string) *ValkeyRepository {
	return &ValkeyRepository{address: addr, username: username, password: password}
}

func (r *ValkeyRepository) Open() error {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{r.address}, Username: r.username, Password: r.password})
	r.client = client
	return err
}

func (r *ValkeyRepository) Set(key string, value []byte) error {
	ctx := context.Background()
	return r.client.Do(ctx, r.client.B().Set().Key(key).Value(valkey.BinaryString(value)).Build()).Error()

}

func (r *ValkeyRepository) Get(key string) ([]byte, error) {
	ctx := context.Background()
	return r.client.Do(ctx, r.client.B().Get().Key(key).Build()).AsBytes()
}

func (r *ValkeyRepository) Close() error {
	r.client.Close()
	return nil
}
