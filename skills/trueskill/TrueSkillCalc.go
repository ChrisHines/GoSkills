package trueskill

import (
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/numerics"
	"sort"
)

// Calculates TrueSkill using a full factor graph.
type TrueSkillCalc struct{}

// Calculates new ratings based on the prior ratings and team ranks use 1 for first place, repeat the number for a tie (e.g. 1, 2, 2).
func (calc *TrueSkillCalc) CalcNewRatings(gi *skills.GameInfo, teams []skills.Team, ranks ...int) skills.PlayerRatings {
	// Basic argument checking
	validateTeamCount(teams, trueSkillTeamRange)
	validatePlayersPerTeam(teams, trueSkillPlayerRange)

	// Copy slices so we don't confuse the client code
	steams := append([]skills.Team{}, teams...)
	sranks := append([]int{}, ranks...)

	// Make sure things are in order
	sort.Sort(skills.NewRankedTeams(steams, sranks))

	factorGraph := NewTrueSkillFactorGraph(gi, teams, sranks...)
	factorGraph.BuildGraph()
	factorGraph.RunSchedule()

	return factorGraph.GetUpdatedRatings()
}

func (calc *TrueSkillCalc) CalcMatchQual(gi *skills.GameInfo, teams []skills.Team) float64 {
	return 0
}

var (
	trueSkillTeamRange   = numerics.AtLeast(2)
	trueSkillPlayerRange = numerics.AtLeast(1)
)
