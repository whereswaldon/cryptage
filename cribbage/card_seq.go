package cribbage

type Sequence []Play

type Play struct {
	Card   *Card
	Player int
}

func NewSeq() *Sequence {
	plays := make([]Play, 0)
	seq := Sequence(plays)
	return &seq
}

func (s *Sequence) Play(player int, card *Card) {
	*s = append(*s, Play{Card: card, Player: player})
}

func (s *Sequence) Get(index int) (int, *Card) {
	return (*s)[index].Player, (*s)[index].Card
}

func (s *Sequence) Size() int {
	return len(*s)
}
