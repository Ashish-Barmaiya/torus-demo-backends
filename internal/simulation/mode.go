package simulation

import "fmt"

type Mode string

const (
	ModeNormal Mode = "normal"
	ModeSlow   Mode = "slow"
	ModeError  Mode = "error"
)

func (m Mode) Valid() bool {
	switch m {
	case ModeNormal, ModeSlow, ModeError:
		return true
	default:
		return false
	}
}

func ParseMode(value string) (Mode, error) {
	mode := Mode(value)

	if !mode.Valid() {
		return "", fmt.Errorf("unsupported simulation mode %q", value)
	}

	return mode, nil
}
