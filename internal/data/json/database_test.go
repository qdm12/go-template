package json

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/qdm12/go-template/internal/data/memory"
	"github.com/qdm12/go-template/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Database(t *testing.T) {
	t.Parallel()

	memory, err := memory.NewDatabase()
	require.NoError(t, err)

	filePath := filepath.Join(t.TempDir(), "database.json")

	// Initialize database file
	database := NewDatabase(memory, filePath)

	runError, err := database.Start(context.Background())
	require.NoError(t, err)
	assert.Nil(t, runError)

	userOne := models.User{
		ID: 1,
	}
	err = database.CreateUser(context.Background(), userOne)
	require.NoError(t, err)

	err = database.Stop()
	require.NoError(t, err)

	runError, err = database.Start(context.Background())
	require.NoError(t, err)
	assert.Nil(t, runError)

	// Check we still get the user previously created and stored on file
	userRetrieved, err := database.GetUserByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, userOne, userRetrieved)

	userTwo := models.User{
		ID: 2,
	}
	err = database.CreateUser(context.Background(), userTwo)
	require.NoError(t, err)

	// Check we still have the user previously created
	userRetrieved, err = database.GetUserByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, userOne, userRetrieved)

	// Check we have the new user
	userRetrieved, err = database.GetUserByID(context.Background(), 2)
	require.NoError(t, err)
	assert.Equal(t, userTwo, userRetrieved)

	err = database.Stop()
	require.NoError(t, err)
}
