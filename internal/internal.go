package internal

type Context struct {
	MapUrl          string
	LocationAreaUrl string
	Prev            string
	Next            string
	Cache           *Cache
	CommandArgs     []string
}
