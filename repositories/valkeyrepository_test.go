package repositories

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/valkey"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestValkeyRepository(t *testing.T) {
	valkeyContainer, err := valkey.Run(context.Background(), "docker.io/valkey/valkey:latest",
		valkey.WithConfigFile("./../.docker/valkey.conf"),

		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections tcp").
				WithOccurrence(1).WithStartupTimeout(20*time.Second)))

	defer func() {
		if err := testcontainers.TerminateContainer(valkeyContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
			t.FailNow()
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		t.FailNow()
	}
	uri, err := valkeyContainer.ConnectionString(context.Background())
	valkeyAddress := strings.TrimPrefix(uri, "redis://")
	t.Run("Test init",
		func(t *testing.T) {
			repo := NewValkeyRepository("user", "pass", valkeyAddress)
			assert.NotNil(t, repo)
			assert.Equal(t, "user", repo.username)
			assert.Equal(t, "pass", repo.password)
			assert.Equal(t, valkeyAddress, repo.address)
		})

	t.Run("Test open",
		func(t *testing.T) {
			repo := NewValkeyRepository("default", "", valkeyAddress)
			err := repo.Open()
			require.NoError(t, err)
			defer repo.Close()
		})
	t.Run("Test set and get",
		func(t *testing.T) {
			repo := NewValkeyRepository("default", "", valkeyAddress)
			err := repo.Open()
			require.NoError(t, err)
			defer repo.Close()

			key := "testkey"
			value := []byte("testvalue")

			err = repo.Set(key, value)
			require.NoError(t, err)

			result, err := repo.Get(key)
			require.NoError(t, err)
			assert.Equal(t, value, result)
		})
	t.Run("Test get nonexistent key",
		func(t *testing.T) {
			repo := NewValkeyRepository("default", "", valkeyAddress)
			err := repo.Open()
			require.NoError(t, err)
			defer repo.Close()

			result, err := repo.Get("nonexistent")
			assert.Nil(t, result)
			assert.Error(t, err)
		})
	t.Run("Test close",
		func(t *testing.T) {
			repo := NewValkeyRepository("default", "", valkeyAddress)
			err := repo.Open()
			require.NoError(t, err)
			repo.Close()
		})

}
