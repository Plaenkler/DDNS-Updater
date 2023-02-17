package model

type User struct {
	ID       int64
	Username string
	Password string
}

type Updater struct {
	ID         int64
	Provider   string
	Domain     string
	User       string
	Password   string
	LastIPAddr string
}
