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

func (s *Sequence) CanPlay(card *Card) bool {
	if card != nil {
		return s.Total()+card.Value() < 32
	}
	return false
}

func (s *Sequence) Play(player int, card *Card) {
	if s.CanPlay(card) {
		*s = append(*s, Play{Card: card.Copy(), Player: player})
	}
}

func (s *Sequence) Get(index int) (int, *Card) {
	return (*s)[index].Player, (*s)[index].Card
}

func (s *Sequence) Size() int {
	return len(*s)
}

func (s *Sequence) Last() (int, *Card) {
	return s.Get(s.Size() - 1)
}

func (s *Sequence) Total() int {
	val := 0
	for _, play := range *s {
		val += play.Card.Value()
	}
	return val
}
