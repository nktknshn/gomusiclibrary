package sqlutil

import (
	"time"

	"zombiezen.com/go/sqlite"
)

func SetDateTextOr(stmt *sqlite.Stmt, param string, value time.Time, orValue time.Time) {
	if value.IsZero() {
		stmt.SetText(param, orValue.Format(time.DateTime))
	} else {
		stmt.SetText(param, value.Format(time.DateTime))
	}
}

func SetDateTextIfNotZero(stmt *sqlite.Stmt, param string, value time.Time) {
	if !value.IsZero() {
		stmt.SetText(param, value.Format(time.DateTime))
	}
}

func DateFromColumnMaybeNull(stmt *sqlite.Stmt, column string) (time.Time, error) {

	idxCol := stmt.ColumnIndex(column)

	if stmt.ColumnType(idxCol) == sqlite.TypeNull {
		return time.Time{}, nil
	}

	return time.Parse(time.DateTime, stmt.ColumnText(idxCol))
}
