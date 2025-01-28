package database

import (
	"database/sql"
	"testing"
)

type mockDB struct {
	*sql.DB
	execCalls []string
}

func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.execCalls = append(m.execCalls, query)
	return nil, nil
}
func TestSetAndGetDB(t *testing.T) {
	db := &sql.DB{}
	SetDB(db)

	if got := GetDB(); got != db {
		t.Errorf("GetDB() = %v, want %v", got, db)
	}
}
