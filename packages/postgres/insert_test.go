package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType struct {
	FieldA       int  `db:"field_a"`
	FieldB       bool `db:"field_b"`
	IgnoredField int
}

func TestInsertQuery(t *testing.T) {
	query, values := insertQuery("table", testType{1, false, 10}, testType{2, true, 20})

	expectedQuery := "INSERT INTO table (field_a,field_b) VALUES ($1,$2),($3,$4)"
	expectedValues := []interface{}{1, false, 2, true}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

func BenchmarkInsertQuery1(b *testing.B) {
	b.StopTimer()
	rows := generateInserts(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		insertQuery("table", rows...)
	}
}

func BenchmarkInsertQuery10000(b *testing.B) {
	b.StopTimer()
	rows := generateInserts(10000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		insertQuery("table", rows...)
	}
}

func generateInserts(n int) []interface{} {
	rows := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rows[i] = testType{FieldA: 0, FieldB: false}
	}
	return rows
}
