package trueskill

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/numerics"
)

func ValidateTeamCountAndPlayersCountPerTeam(teams []skills.Team, teamsAllowed, playersAllowed numerics.Range) {
	if n := len(teams); !teamsAllowed.In(n) {
		panic(fmt.Errorf("len(teams) [%v] outside of expected range [%v]", n, teamsAllowed))
	}
	for _, t := range teams {
		if n := t.PlayerCount(); !playersAllowed.In(n) {
			panic(fmt.Errorf("PlayerCount [%v] outside of expected range [%v]", n, playersAllowed))
		}
	}
}

func cond(c bool, t, f int) int {
	if c {
		return t
	}
	return f
}
