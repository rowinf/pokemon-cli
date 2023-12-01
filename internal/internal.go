package internal

type Context struct {
	MapUrl      string
	Prev        string
	Next        string
	Cache       *Cache
	CommandArgs []string
}
