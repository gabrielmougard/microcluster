package update

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/canonical/lxd/lxd/db/schema"
)

// CreateSchema is the default schema applied when bootstrapping the database.
const CreateSchema = `
CREATE TABLE schemas (
  id          INTEGER    PRIMARY  KEY    AUTOINCREMENT  NOT  NULL,
  version     INTEGER    NOT      NULL,
  updated_at  DATETIME   NOT      NULL,
  UNIQUE      (version)
);
`

// Template for schema files (can't use backticks since we need to use backticks
// inside the template itself).
const dotGoTemplate = "package %s\n\n" +
	"// DO NOT EDIT BY HAND\n" +
	"//\n" +
	"// This code was generated by the schema.DotGo function. If you need to\n" +
	"// modify the database schema, please add a new schema update to update.go\n" +
	"// and the run 'make update-schema'.\n" +
	"const freshSchema = `\n" +
	"%s`\n"

type SchemaUpdateManager struct {
	overrides map[int]schema.Update
	updates   map[int]schema.Update
}

func NewSchema() *SchemaUpdateManager {
	return &SchemaUpdateManager{
		overrides: map[int]schema.Update{
			2: overrideUpdateFromV1, // If this override is executed before `updateFromV1`, we should not apply `updateFromV1`.
		},
		updates: map[int]schema.Update{
			1: updateFromV0,
			2: updateFromV1,
		},
	}
}

func (m *SchemaUpdateManager) Schema() *SchemaUpdate {
	schema := NewFromMap(m.overrides, m.updates)
	schema.Fresh("")
	return schema
}

func (m *SchemaUpdateManager) AppendSchema(extensions map[int]schema.Update) {
	currentVersion := len(m.updates)
	schema := NewFromMap(nil, extensions)
	for _, extension := range schema.updates {
		m.updates[currentVersion+1] = extension.update
		currentVersion = len(m.updates)
	}
}

func (m *SchemaUpdateManager) SchemaDotGo() error {
	// Apply all the updates that we have on a pristine database and dump
	// the resulting schema.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return fmt.Errorf("failed to open schema.go for writing: %w", err)
	}

	schema := NewFromMap(m.overrides, m.updates)

	_, err = schema.Ensure(db)
	if err != nil {
		return err
	}

	dump, err := schema.Dump(db)
	if err != nil {
		return err
	}

	// Passing 1 to runtime.Caller identifies our caller.
	_, filename, _, _ := runtime.Caller(1)

	file, err := os.Create(path.Join(path.Dir(filename), "schema_update.go"))
	if err != nil {
		return fmt.Errorf("failed to open Go file for writing: %w", err)
	}

	pkg := path.Base(path.Dir(filename))
	_, err = file.Write([]byte(fmt.Sprintf(dotGoTemplate, pkg, dump)))
	if err != nil {
		return fmt.Errorf("failed to write to Go file: %w", err)
	}

	return nil
}

func updateFromV0(ctx context.Context, tx *sql.Tx) error {
	stmt := fmt.Sprintf(`
%s

CREATE TABLE internal_token_records (
  id           INTEGER         PRIMARY  KEY    AUTOINCREMENT  NOT  NULL,
  name         TEXT            NOT      NULL,
  secret       TEXT            NOT      NULL,
  UNIQUE       (name),
  UNIQUE       (secret)
);

CREATE TABLE internal_cluster_members (
  id                   INTEGER   PRIMARY  KEY    AUTOINCREMENT  NOT  NULL,
  name                 TEXT      NOT      NULL,
  address              TEXT      NOT      NULL,
  certificate          TEXT      NOT      NULL,
  schema               INTEGER   NOT      NULL,
  heartbeat            DATETIME  NOT      NULL,
  role                 TEXT      NOT      NULL,
  UNIQUE(name),
  UNIQUE(certificate)
);
`, CreateSchema)

	_, err := tx.ExecContext(ctx, stmt)
	return err
}

func updateFromV1(ctx context.Context, tx *sql.Tx) error {
	stmt := `
ALTER TABLE internal_cluster_members ADD COLUMN internal_api_extensions TEXT;
ALTER TABLE internal_cluster_members ADD COLUMN external_api_extensions TEXT;
`
	_, err := tx.ExecContext(ctx, stmt)
	return err
}

// This update is an override for the `updateFromV1` update.
// If this override is executed before `updateFromV1`, we should not apply `updateFromV1`.
//
// This update is necessary when an already clustered node with schema version 1 is upgraded.
func overrideUpdateFromV1(ctx context.Context, tx *sql.Tx) error {
	stmt := `
CREATE TABLE internal_cluster_members_new (
	id                      INTEGER   PRIMARY  KEY    AUTOINCREMENT  NOT  NULL,
	name                    TEXT      NOT      NULL,
	address                 TEXT      NOT      NULL,
	certificate             TEXT      NOT      NULL,
	schema                  INTEGER   NOT      NULL,
	heartbeat               DATETIME  NOT      NULL,
	role                    TEXT      NOT      NULL,
	internal_api_extensions TEXT,
	external_api_extensions TEXT,
	UNIQUE(name),
	UNIQUE(certificate)
);

INSERT INTO internal_cluster_members_new (id, name, address, certificate, schema, heartbeat, role, internal_api_extensions, external_api_extensions)
SELECT id, name, address, certificate, schema, heartbeat, role, NULL, NULL FROM internal_cluster_members;

DROP TABLE internal_cluster_members;
ALTER TABLE internal_cluster_members_new RENAME TO internal_cluster_members;
`
	_, err := tx.ExecContext(ctx, stmt)
	return err
}
