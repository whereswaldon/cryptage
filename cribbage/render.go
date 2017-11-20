package cribbage

import "github.com/fatih/color"

const (
	SPADE   = "♠"
	HEART   = "♥"
	DIAMOND = "♦"
	CLUB    = "♣"
)

func RenderCard(card *Card) string {
	var printer func(...interface{}) string
	if isRed(card) {
		printer = color.New(color.BgWhite, color.FgRed).SprintFunc()
	} else {
		printer = color.New(color.BgWhite, color.FgBlack).SprintFunc()
	}
	return printer(renderRank(card), renderSuit(card))
}

func RenderCards(cards []*Card) string {
	out := ""
	for i, card := range cards {
		out += RenderCard(card)
		if i < len(cards) {
			out += " "
		}
	}
	return out
}

func isRed(card *Card) bool {
	return card.Suit == "Hearts" || card.Suit == "Diamonds"
}

func renderSuit(card *Card) string {
	switch card.Suit {
	case "Hearts":
		return HEART
	case "Spades":
		return SPADE
	case "Clubs":
		return CLUB
	case "Diamonds":
		return DIAMOND
	default:
		return "?"
	}
}

func renderRank(card *Card) string {
	switch card.Rank {
	case "Two":
		return "2"
	case "Three":
		return "3"
	case "Four":
		return "4"
	case "Five":
		return "5"
	case "Six":
		return "6"
	case "Seven":
		return "7"
	case "Eight":
		return "8"
	case "Nine":
		return "9"
	case "Ten":
		return "10"
	case "Jack":
		return "J"
	case "Queen":
		return "Q"
	case "King":
		return "K"
	case "Ace":
		return "A"
	default:
		return "?"
	}
}
