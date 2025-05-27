package entity

type Location struct {
	ID           uint   `db:"id"`
	State        string `db:"state"`
	County       string `db:"county"`
	Municipality string `db:"municipality"`
}
