package example

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/dakimura/migrent/example/ent/user"

	"github.com/dakimura/migrent/example/ent"
)

type UserMasterDataMigration struct {
	Data   []ent.User
	Client *ent.UserClient
}

func NewUserMasterDataMigration(data []ent.User, client *ent.UserClient,
) *UserMasterDataMigration {
	return &UserMasterDataMigration{
		Data:   data,
		Client: client,
	}
}

// Up inserts data to User entity
func (m *UserMasterDataMigration) Up(ctx context.Context) error {
	for _, rec := range m.Data {
		err := m.Client.Create().SetID(rec.ID).SetAge(rec.Age).SetName(rec.Name).
			OnConflictColumns(user.FieldID).UpdateNewValues().Exec(ctx)
		if err != nil {
			return fmt.Errorf("create User entities: %w", err)
		}
	}

	return nil
}

// Down deletes data from User entity
func (m *UserMasterDataMigration) Down(ctx context.Context) error {
	users := make([]uuid.UUID, len(m.Data))
	for i, u := range m.Data {
		users[i] = u.ID
	}

	_, err := m.Client.
		Delete().
		Where(user.IDIn(users...)).Exec(ctx)

	if err != nil {
		return fmt.Errorf("delete User entities: %w", err)
	}

	return nil
}
