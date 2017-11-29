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

const (
	CLAIM_PAIR    = 0
	CLAIM_FIFTEEN = 1
)

func (s *Sequence) WorthPoints(pointType int) bool {
	switch pointType {
	case CLAIM_PAIR:
		if s.Size() < 2 {
			return false
		}
		lastIdx := s.Size() - 1
		secLastIdx := lastIdx - 1
		_, lastCard := s.Get(lastIdx)
		_, secLastCard := s.Get(secLastIdx)
		return lastCard.Rank == secLastCard.Rank
	case CLAIM_FIFTEEN:
		if s.Size() < 2 {
			return false
		}
		total := 0
		for start := s.Size() - 1; total < 15; start++ {
			_, card := s.Get(start)
			total += card.Value()
		}
		return total == 15
	}
	return false
}
