package cribbage

import (
	"fmt"
	"github.com/fatih/color"
)

const (
	SPADE   = "♠"
	HEART   = "♥"
	DIAMOND = "♦"
	CLUB    = "♣"
)

func RenderCard(card *Card) string {
	var printer func(...interface{}) string
	if isUnknown(card) {
		printer = color.New(color.BgBlue, color.FgWhite).SprintFunc()
	} else if isRed(card) {
		printer = color.New(color.BgWhite, color.FgRed).SprintFunc()
	} else {
		printer = color.New(color.BgWhite, color.FgBlack).SprintFunc()
	}
	return printer(renderRank(card), renderSuit(card))
}

func RenderHand(hand *Hand) string {
	out := ""
	for i, card := range hand.cards {
		if card == nil {
			card = &Card{}
		}
		out += fmt.Sprintf("%d:", i) + RenderCard(card)
		if i < len(hand.cards) {
			out += " "
		}
	}
	return out
}

func RenderSeq(seq *Sequence) string {
	out := ""
	for i := 0; i < seq.Size(); i++ {
		out += fmt.Sprintf("%d:", i) + RenderCard(seq.Get(i)) + " "
	}
	return out
}

func isRed(card *Card) bool {
	return card.Suit == "Hearts" || card.Suit == "Diamonds"
}

func isUnknown(card *Card) bool {
	return card.Suit == "" || card.Suit == ""
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
		return "▹"
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
		return "◃"
	}
}
