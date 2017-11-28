package cribbage

type PlayerInfo struct {
	NumPlayers, LocalPlayerNum, CurrentDealerNum int
}

func NewPlayerInfo(numPlayers, localPlayerNum, currDealerNum int) *PlayerInfo {
	return &PlayerInfo{
		NumPlayers:       numPlayers,
		LocalPlayerNum:   localPlayerNum,
		CurrentDealerNum: currDealerNum,
	}
}

func (p *PlayerInfo) LocalPlayerIsDealer() bool {
	return p.LocalPlayerNum == p.CurrentDealerNum
}

func (p *PlayerInfo) OpponentNum() int {
	if p.LocalPlayerNum == 1 {
		return 2
	}
	return 1
}
