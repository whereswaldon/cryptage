package cribbage

const (
	DRAW_STATE         = 0
	DISCARD_STATE      = 1
	DISCARD_WAIT_STATE = 2
	CUT_STATE          = 3
	CUT_WAIT_STATE     = 4
	CIRCULAR_STATE     = 5
	INTERNAL_STATE     = 6
	CRIB_STATE         = 7
	END_STATE          = 8
)

type State uint

func instructionsForState(s State) string {
	switch s {
	case DRAW_STATE:
		return STR_DRAW_INSTRUCTIONS
	case DISCARD_STATE:
		return STR_DISCARD_INSTRUCTIONS
	case DISCARD_WAIT_STATE:
		return STR_DISCARD_WAIT_INSTRUCTIONS
	case CUT_STATE:
		return STR_CUT_INSTRUCTIONS
	case CUT_WAIT_STATE:
		return STR_CUT_WAIT_INSTRUCTIONS
	case CIRCULAR_STATE:
		return STR_CIRCULAR_INSTRUCTIONS
	case INTERNAL_STATE:
		return STR_INTERNAL_INSTRUCTIONS
	case CRIB_STATE:
		return STR_CRIB_INSTRUCTIONS
	case END_STATE:
		return STR_END_INSTRUCTIONS
	default:
		return "Unknown game state, cannot fetch instructions"
	}
}
