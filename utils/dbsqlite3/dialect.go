package dbsqlite3

type sqlite3Dialect int

var dialect sqlite3Dialect

func (sqlite3Dialect) Bra() string {
	return "`"
}

func (sqlite3Dialect) Ket() string {
	return "`"
}
