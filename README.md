## migrent - A data migration tool for ent

migrent is a library to manage data migrations for applications using [ent](https://github.com/ent/ent).

Ent is one of the most popular ORM libraries for Golang, and [automatic schema migration](https://entgo.io/docs/migrate) is already supported.
However, with regard to data migration/incremental migration history management, there is some room for adding more functionality.
When an application using ent needs master data (=data needed to be registered in advance), the users need to have the code to manage the data by themselves because `Upsert` is not yet implemented, 
and it is not possible to "Up" and "Down" migrations like with
existing tools such as [goose](https://github.com/pressly/goose), [sql-migrate](https://github.com/rubenv/sql-migrate),
and [golang-migrate](https://github.com/golang-migrate/migrate).

migrent enables the applying and rolling back of data migrations by creating an internal "migration" entity in the DB, while
keeping the great part of ent, such as type-safety.

## Quick Installation

```console
go get github.com/dakimura/migrent
```

## Usage

```go
package main

import (
	"context"

	"github.com/dakimura/migrent/example"

	"github.com/dakimura/migrent/example/ent"

	"github.com/dakimura/migrent"

	_ "github.com/go-sql-driver/mysql"
)

// example master data definition 
var masterData = []ent.User{
	{Age: 12, Name: "Alice"},
	{Age: 24, Name: "Bob"},
	{Age: 36, Name: "Carol"},
	{Age: 48, Name: "David"},
	{Age: 60, Name: "Eve"},
}

func main() {
	ctx := context.Background()

	dialect := "mysql"
	dsn := "user:password@tcp(IPAddress:Port)/dbname?parse=True"
	
	// initialize ent client for your DB
	userDBCli, _ := ent.Open(dialect, dsn)
	// initialize migrent client (you can use the same dsn).
	// because migrent creates an internal entity to manage the migration history
	// in the DB, an ent client for it is necessary
	cli, _ := migrent.Open(dialect, dsn)

	// define your migration
	migrations := map[migrent.MigrationName]migrent.Migration{
		"20220712_user_data_v1": example.NewUserMasterDataMigration(masterData, userDBCli.User),
	}

	// ---  execute migration(Up)
	cli.Up(ctx, migrations)

	// --- execute migration(Down)
	cli.Down(ctx, migrations)
}

```

You need to define what operation you need when `Up` or `Down` is called.

```go
type Migration interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}
```

- Each migration doesn't have to be idempotent, but we recommend it to be.
  If your migration is idempotent, you can update the value of existent records
  by just changing the values in your data files and change the migration name like "data_v1" to "data_v2",
  and run the migration.
  Because migrent manages whether each migration has already been applied or not by the migration name,
  the "data_v2" migration will be executed and update your existent values.

- Migrations are executed by dictionary order (AtoZ) of migration names, so "{timestamp}_{name}" might be a good naming
  for migrations.

The following is the simple example of a migration that inserts master data to User entity.

```go
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
```

## License

migrent is licensed under MIT as found in the [LICENSE file](LICENSE).
