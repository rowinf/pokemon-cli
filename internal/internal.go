package internal

type Context struct {
	LocationAreaUrl string
	CatchPokemonUrl string
	MapUrl          string
	Prev            string
	Next            string
	Pokedex         map[string]Pokemon
	Cache           *Cache
	CommandArgs     []string
}
