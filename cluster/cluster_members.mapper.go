package cluster

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/canonical/lxd/lxd/db/query"
	"github.com/canonical/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var internalClusterMemberObjects = RegisterStmt(`
SELECT internal_cluster_members.id, internal_cluster_members.name, internal_cluster_members.address, internal_cluster_members.certificate, internal_cluster_members.schema_internal, internal_cluster_members.schema_external, internal_cluster_members.api_extensions, internal_cluster_members.heartbeat, internal_cluster_members.role
  FROM internal_cluster_members
  ORDER BY internal_cluster_members.name
`)

var internalClusterMemberObjectsByAddress = RegisterStmt(`
SELECT internal_cluster_members.id, internal_cluster_members.name, internal_cluster_members.address, internal_cluster_members.certificate, internal_cluster_members.schema_internal, internal_cluster_members.schema_external, internal_cluster_members.api_extensions, internal_cluster_members.heartbeat, internal_cluster_members.role
  FROM internal_cluster_members
  WHERE ( internal_cluster_members.address = ? )
  ORDER BY internal_cluster_members.name
`)

var internalClusterMemberObjectsByName = RegisterStmt(`
SELECT internal_cluster_members.id, internal_cluster_members.name, internal_cluster_members.address, internal_cluster_members.certificate, internal_cluster_members.schema_internal, internal_cluster_members.schema_external, internal_cluster_members.api_extensions, internal_cluster_members.heartbeat, internal_cluster_members.role
  FROM internal_cluster_members
  WHERE ( internal_cluster_members.name = ? )
  ORDER BY internal_cluster_members.name
`)

var internalClusterMemberID = RegisterStmt(`
SELECT internal_cluster_members.id FROM internal_cluster_members
  WHERE internal_cluster_members.name = ?
`)

var internalClusterMemberCreate = RegisterStmt(`
INSERT INTO internal_cluster_members (name, address, certificate, schema_internal, schema_external, api_extensions, heartbeat, role)
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`)

var internalClusterMemberDeleteByAddress = RegisterStmt(`
DELETE FROM internal_cluster_members WHERE address = ?
`)

var internalClusterMemberUpdate = RegisterStmt(`
UPDATE internal_cluster_members
  SET name = ?, address = ?, certificate = ?, schema_internal = ?, schema_external = ?, api_extensions = ?, heartbeat = ?, role = ?
 WHERE id = ?
`)

// internalClusterMemberColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the InternalClusterMember entity.
func internalClusterMemberColumns() string {
	return "internal_cluster_members.id, internal_cluster_members.name, internal_cluster_members.address, internal_cluster_members.certificate, internal_cluster_members.schema_internal, internal_cluster_members.schema_external, internal_cluster_members.api_extensions, internal_cluster_members.heartbeat, internal_cluster_members.role"
}

// getInternalClusterMembers can be used to run handwritten sql.Stmts to return a slice of objects.
func getInternalClusterMembers(ctx context.Context, stmt *sql.Stmt, args ...any) ([]InternalClusterMember, error) {
	objects := make([]InternalClusterMember, 0)

	dest := func(scan func(dest ...any) error) error {
		i := InternalClusterMember{}
		err := scan(&i.ID, &i.Name, &i.Address, &i.Certificate, &i.SchemaInternal, &i.SchemaExternal, &i.APIExtensions, &i.Heartbeat, &i.Role)
		if err != nil {
			return err
		}

		objects = append(objects, i)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"internal_cluster_members\" table: %w", err)
	}

	return objects, nil
}

// getInternalClusterMembersRaw can be used to run handwritten query strings to return a slice of objects.
func getInternalClusterMembersRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]InternalClusterMember, error) {
	objects := make([]InternalClusterMember, 0)

	dest := func(scan func(dest ...any) error) error {
		i := InternalClusterMember{}
		err := scan(&i.ID, &i.Name, &i.Address, &i.Certificate, &i.SchemaInternal, &i.SchemaExternal, &i.APIExtensions, &i.Heartbeat, &i.Role)
		if err != nil {
			return err
		}

		objects = append(objects, i)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"internal_cluster_members\" table: %w", err)
	}

	return objects, nil
}

// GetInternalClusterMembers returns all available internal_cluster_members.
// generator: internal_cluster_member GetMany
func GetInternalClusterMembers(ctx context.Context, tx *sql.Tx, filters ...InternalClusterMemberFilter) ([]InternalClusterMember, error) {
	var err error

	// Result slice.
	objects := make([]InternalClusterMember, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = Stmt(tx, internalClusterMemberObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"internalClusterMemberObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Name != nil && filter.Address == nil {
			args = append(args, []any{filter.Name}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, internalClusterMemberObjectsByName)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"internalClusterMemberObjectsByName\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(internalClusterMemberObjectsByName)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"internalClusterMemberObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Address != nil && filter.Name == nil {
			args = append(args, []any{filter.Address}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, internalClusterMemberObjectsByAddress)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"internalClusterMemberObjectsByAddress\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(internalClusterMemberObjectsByAddress)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"internalClusterMemberObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Address == nil && filter.Name == nil {
			return nil, fmt.Errorf("Cannot filter on empty InternalClusterMemberFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getInternalClusterMembers(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getInternalClusterMembersRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"internal_cluster_members\" table: %w", err)
	}

	return objects, nil
}

// GetInternalClusterMember returns the internal_cluster_member with the given key.
// generator: internal_cluster_member GetOne
func GetInternalClusterMember(ctx context.Context, tx *sql.Tx, name string) (*InternalClusterMember, error) {
	filter := InternalClusterMemberFilter{}
	filter.Name = &name

	objects, err := GetInternalClusterMembers(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"internal_cluster_members\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "InternalClusterMember not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"internal_cluster_members\" entry matches")
	}
}

// GetInternalClusterMemberID return the ID of the internal_cluster_member with the given key.
// generator: internal_cluster_member ID
func GetInternalClusterMemberID(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	stmt, err := Stmt(tx, internalClusterMemberID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"internalClusterMemberID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, name)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "InternalClusterMember not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"internal_cluster_members\" ID: %w", err)
	}

	return id, nil
}

// InternalClusterMemberExists checks if a internal_cluster_member with the given key exists.
// generator: internal_cluster_member Exists
func InternalClusterMemberExists(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	_, err := GetInternalClusterMemberID(ctx, tx, name)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateInternalClusterMember adds a new internal_cluster_member to the database.
// generator: internal_cluster_member Create
func CreateInternalClusterMember(ctx context.Context, tx *sql.Tx, object InternalClusterMember) (int64, error) {
	// Check if a internal_cluster_member with the same key exists.
	exists, err := InternalClusterMemberExists(ctx, tx, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"internal_cluster_members\" entry already exists")
	}

	args := make([]any, 8)

	// Populate the statement arguments.
	args[0] = object.Name
	args[1] = object.Address
	args[2] = object.Certificate
	args[3] = object.SchemaInternal
	args[4] = object.SchemaExternal
	args[5] = object.APIExtensions
	args[6] = object.Heartbeat
	args[7] = object.Role

	// Prepared statement to use.
	stmt, err := Stmt(tx, internalClusterMemberCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"internalClusterMemberCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"internal_cluster_members\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"internal_cluster_members\" entry ID: %w", err)
	}

	return id, nil
}

// DeleteInternalClusterMember deletes the internal_cluster_member matching the given key parameters.
// generator: internal_cluster_member DeleteOne-by-Address
func DeleteInternalClusterMember(ctx context.Context, tx *sql.Tx, address string) error {
	stmt, err := Stmt(tx, internalClusterMemberDeleteByAddress)
	if err != nil {
		return fmt.Errorf("Failed to get \"internalClusterMemberDeleteByAddress\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(address)
	if err != nil {
		return fmt.Errorf("Delete \"internal_cluster_members\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "InternalClusterMember not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d InternalClusterMember rows instead of 1", n)
	}

	return nil
}

// UpdateInternalClusterMember updates the internal_cluster_member matching the given key parameters.
// generator: internal_cluster_member Update
func UpdateInternalClusterMember(ctx context.Context, tx *sql.Tx, name string, object InternalClusterMember) error {
	id, err := GetInternalClusterMemberID(ctx, tx, name)
	if err != nil {
		return err
	}

	stmt, err := Stmt(tx, internalClusterMemberUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"internalClusterMemberUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Name, object.Address, object.Certificate, object.SchemaInternal, object.SchemaExternal, object.APIExtensions, object.Heartbeat, object.Role, id)
	if err != nil {
		return fmt.Errorf("Update \"internal_cluster_members\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}
