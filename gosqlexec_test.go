package gosqlexec

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gocraft/dbr/dialect"
	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/require"
)

func TestLineByLineReader(t *testing.T) {
	sqlSyntax := `DROP TABLE IF EXISTS Table1, Table2, Table3, Table4;` + "\n"
	loadedSyntax := LineByLineReader("dbmock/sql_tester.sql")
	if loadedSyntax != sqlSyntax {
		t.Error("LineByLineReader test FAIL")
	}
}

func TestSQLMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.MySQL,
	}
	sess := conn.NewSession(nil)

	mock.ExpectQuery("SELECT id FROM suggestions").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
	id, err := sess.Select("id").From("suggestions").ReturnInt64s()
	require.NoError(t, err)
	require.Equal(t, []int64{1, 2}, id)

	mock.ExpectClose()
	conn.Close()

	require.NoError(t, mock.ExpectationsWereMet())
}
