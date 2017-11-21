package cribbage

const (
	DRAW_STATE         = 0
	DISCARD_STATE      = 1
	DISCARD_WAIT_STATE = 2
	CIRCULAR_STATE     = 3
	INTERNAL_STATE     = 4
	CRIB_STATE         = 5
	END_STATE          = 6
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
