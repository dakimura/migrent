package migrent

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"

	"github.com/dakimura/migrent/ent/migration"

	"github.com/dakimura/migrent/ent"
)

type MigrationName string
type Migration interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

type Client struct {
	entclient *ent.Client
}

func NewClient(entclient *ent.Client) *Client {
	return &Client{entclient: entclient}
}

func Open(driverName, dataSourceName string, options ...ent.Option) (*Client, error) {
	cli, err := ent.Open(driverName, dataSourceName, options...)
	if err != nil {
		return nil, err
	}
	return &Client{entclient: cli}, nil
}

func OpenByMySQLDB(db *sql.DB) *Client {
	return OpenByDB(db, "mysql")
}

func OpenByDB(db *sql.DB, driver string) *Client {
	drv := entsql.OpenDB(driver, db)
	client := ent.NewClient(ent.Driver(drv))

	return &Client{entclient: client}
}

func (c *Client) Up(ctx context.Context, m map[MigrationName]Migration) error {
	// create internal table if not exists
	err := c.createMigrationTable(ctx)
	if err != nil {
		return err
	}

	for name, mi := range m {
		// --- check if the migration is already applied or not
		m, err := c.entclient.Migration.Query().
			Where(migration.NameEQ(string(name))).
			All(ctx)
		if err != nil {
			return fmt.Errorf("querying migration history for %s: %w", name, err)
		}

		// this migration is already applied to DB, skip
		if len(m) > 0 {
			continue
		}

		// --- apply the migration
		err = mi.Up(ctx)
		if err != nil {
			return fmt.Errorf("migration(Up) for %s: %w", name, err)
		}

		// --- record the migration to the internal table
		_, err = c.entclient.Migration.Create().SetName(string(name)).SetAppliedAt(time.Now()).Save(ctx)
		if err != nil {
			return fmt.Errorf("record migration(Up) of %s: %w", name, err)
		}
	}

	return nil
}

func (c *Client) Down(ctx context.Context, m map[MigrationName]Migration) error {
	// implement me. the below is currently just a copypaste of Up()
	err := c.createMigrationTable(ctx)
	if err != nil {
		return err
	}

	for name, mi := range m {
		// --- check if the migration is already applied or not
		m, err := c.entclient.Migration.Query().
			Where(migration.NameEQ(string(name))).
			All(ctx)
		if err != nil {
			return fmt.Errorf("querying migration history for %s: %w", name, err)
		}

		// this migration is not applied to DB, skip
		if len(m) == 0 {
			continue
		}

		// --- apply the migration
		err = mi.Down(ctx)
		if err != nil {
			return fmt.Errorf("migration(Down) for %s: %w", name, err)
		}

		// --- record the migration to the internal table
		_, err = c.entclient.Migration.Delete().Where(migration.NameEQ(string(name))).Exec(ctx)
		if err != nil {
			return fmt.Errorf("record migration(Down) of %s: %w", name, err)
		}
	}

	return nil
}

// TODO: where to run "defer client.Close()" ?

// createMigrationTable creates an internal migration table if not exists
func (c *Client) createMigrationTable(ctx context.Context) error {
	err := c.entclient.Schema.Create(ctx)
	if err != nil {
		return fmt.Errorf("create the internal migration table: %w", err)
	}

	return nil
}
