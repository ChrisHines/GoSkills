package trueskill

import (
	"testing"
)

func TestTrueSkillCalc(t *testing.T) {
	AllTwoPlayerScenarios(t, &TrueSkillCalc{})
	AllTwoTeamScenarios(t, &TrueSkillCalc{})
	AllMultipleTeamScenarios(t, &TrueSkillCalc{})
	// AllPartialPlayScenarios(t, &TrueSkillCalc{})
}
