package example

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"entgo.io/ent/dialect"
	"github.com/dakimura/migrent/ent/enttest"
	enttest2 "github.com/dakimura/migrent/example/ent/enttest"

	"testing"

	"github.com/dakimura/migrent"

	ent2 "github.com/dakimura/migrent/example/ent"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	masterData = []ent2.User{
		{ID: uuid.New(), Age: 12, Name: "Alice"},
		{ID: uuid.New(), Age: 24, Name: "Bob"},
		{ID: uuid.New(), Age: 36, Name: "Carol"},
		{ID: uuid.New(), Age: 48, Name: "David"},
		{ID: uuid.New(), Age: 60, Name: "Eve"},
	}
	masterData2 = []ent2.User{
		{ID: uuid.New(), Age: 20, Name: "Frank"},
		{ID: uuid.New(), Age: 40, Name: "Grace"},
		{ID: uuid.New(), Age: 60, Name: "Heidi"},
	}
)

func TestUserDataMigration(t *testing.T) {
	// --- given ---
	ctx := context.Background()

	inMemoryMigrationDB := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	inMemoryUserDB := enttest2.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		inMemoryMigrationDB.Close()
		inMemoryUserDB.Close()
	})

	migrations := map[migrent.MigrationName]migrent.Migration{
		"user_data1": NewUserMasterDataMigration(masterData, inMemoryUserDB.User),
		"user_data2": NewUserMasterDataMigration(masterData2, inMemoryUserDB.User),
	}

	cli := migrent.NewClient(inMemoryMigrationDB)

	// --- when execute migration(Up)
	err := cli.Up(ctx, migrations)
	if err != nil {
		t.Fatal("failed to execute migration(Up):", err)
	}

	// --- then all the master data should be inserted
	u, _ := inMemoryUserDB.User.Query().All(ctx)
	assert.Len(t, u, len(masterData)+len(masterData2))

	// --- when execute the migration(Up) again
	err = cli.Up(ctx, migrations)
	if err != nil {
		t.Fatal("failed to execute migration(Up):", err)
	}

	// --- then it should not do anything (cause the migration is already applied)
	u, _ = inMemoryUserDB.User.Query().All(ctx)
	assert.Len(t, u, len(masterData)+len(masterData2))

	// --- when, rollback(migration Down) 1 migration,
	rollbackMigration := map[migrent.MigrationName]migrent.Migration{
		"user_data2": NewUserMasterDataMigration(masterData2, inMemoryUserDB.User),
	}

	err = cli.Down(ctx, rollbackMigration)
	if err != nil {
		t.Fatal("failed to execute migration(Down):", err)
	}

	// --- then, the rollback should be done
	u, _ = inMemoryUserDB.User.Query().All(ctx)
	assert.Len(t, u, len(masterData))

	// --- when, all the migrations are roll-backed,
	err = cli.Down(ctx, migrations)
	if err != nil {
		t.Fatal("failed to execute migration(Down):", err)
	}

	// --- then, all the record should be roll-backed
	u, err = inMemoryUserDB.User.Query().All(ctx)
	assert.Len(t, u, 0)
}

func TestDataMigrationOrder(t *testing.T) {
	// --- given ---
	ctx := context.Background()

	var (
		userID = uuid.New()
		data1  = []ent2.User{{ID: userID, Age: 10, Name: "Alice"}}
		data2  = []ent2.User{{ID: userID, Age: 20, Name: "Bob"}}
		data3  = []ent2.User{{ID: userID, Age: 30, Name: "Claire"}}
	)

	inMemoryMigrationDB := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	inMemoryUserDB := enttest2.Open(t, dialect.SQLite, "file:ent2?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		inMemoryMigrationDB.Close()
		inMemoryUserDB.Close()
	})

	migrations := map[migrent.MigrationName]migrent.Migration{
		"3_data1": NewUserMasterDataMigration(data1, inMemoryUserDB.User),
		"2_data2": NewUserMasterDataMigration(data2, inMemoryUserDB.User),
		"1_data3": NewUserMasterDataMigration(data3, inMemoryUserDB.User),
	}

	cli := migrent.NewClient(inMemoryMigrationDB)

	// --- when ---
	// execute migration(Up)
	err := cli.Up(ctx, migrations)
	if err != nil {
		t.Fatal("failed to execute migration(Up):", err)
	}

	// --- then ---
	// "3_data1" should be executed last
	// because it's the last one when migrations are sorted by dictionary order
	u := inMemoryUserDB.User.Query().OnlyX(ctx)
	require.Equal(t, data1[0].Age, u.Age)
	require.Equal(t, data1[0].Name, u.Name)
}
