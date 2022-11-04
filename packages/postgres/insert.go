package postgres

import (
	"context"
	"fmt"
	"strings"
)

// Insert inserts the given rows into the provided table. Each row must be a struct whose database
// fields are exported and tagged with the `db` tag. The struct must not be nested and anonymous
// structs are ignored. A single insert statement is generated in order to maximize throughput.
// Passing zero rows to this function results in a noop.
func (c *Connection) Insert(ctx context.Context, table string, rows ...interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	query, values := insertQuery(table, rows...)
	_, err := c.ExecContext(ctx, query, values...)
	return err
}

// InsertIfNotExists performs an upsert which ignores rows that already exist according to the
// given conflict statement. The conflict statement must be valid SQL syntax for the `target` of
// the ON CONFLICT clause. Passing zero rows to this function results in a noop.
func (c *Connection) InsertIfNotExists(
	ctx context.Context, table, conflict string, rows ...interface{},
) error {
	if len(rows) == 0 {
		return nil
	}
	query, values := insertQuery(table, rows...)
	query = fmt.Sprintf("%s ON CONFLICT %s DO NOTHING", query, conflict)
	_, err := c.ExecContext(ctx, query, values...)
	return err
}

//-------------------------------------------------------------------------------------------------

func insertQuery(table string, rows ...interface{}) (string, []interface{}) {
	// First, we get all the columns that ought to be inserted
	columns := taggedFieldsFromStruct(rows[0], "db")
	columnSet := strings.Join(columns, ",")

	// Then, we generate the corresponding placeholders
	placeholders := valuePlaceholders(len(rows), len(columns))

	// And also get the values from the rows
	values := make([]interface{}, 0, len(rows)*len(columns))
	for _, row := range rows {
		values = append(values, valuesFromTaggedFields(row, "db")...)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, columnSet, placeholders)
	return query, values
}
