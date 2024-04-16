package dbrepo

func (m *sqliteDBRepo) Close() error {
	return m.DB.Close()
}

func (m *sqliteDBRepo) AllUsers() bool {
	return true
}
