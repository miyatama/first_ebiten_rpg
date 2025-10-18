package miyatama

func init() {

}

func NewGame() (*Game, error) {
	game := &Game{}
	game.Init()
	return game, nil
}
