package example

import (
	"context"
	"fmt"

	"github.com/dakimura/migrent/example/ent/user"

	ent "github.com/dakimura/migrent/example/ent"
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
	data := make([]*ent.UserCreate, len(m.Data))
	for i, rec := range m.Data {
		data[i] = m.Client.Create().SetAge(rec.Age).SetName(rec.Name)
	}
	_, err := m.Client.CreateBulk(data...).Save(ctx)

	if err != nil {
		return fmt.Errorf("create User entities: %w", err)
	}
	return nil
}

// Down deletes data from User entity
func (m *UserMasterDataMigration) Down(ctx context.Context) error {
	users := make([]string, len(m.Data))
	for i, u := range m.Data {
		users[i] = u.Name
	}

	_, err := m.Client.
		Delete().
		Where(user.NameIn(users...)).Exec(ctx)

	if err != nil {
		return fmt.Errorf("delete User entities: %w", err)
	}

	return nil
}
