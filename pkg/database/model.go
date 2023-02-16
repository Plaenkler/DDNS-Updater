package database

type User struct {
	ID       int64
	Username string
	Password string
}

type Updater struct {
	ID       int64
	Username string
	Password string
	Domain   string
	IP       string
	Interval int64
}
