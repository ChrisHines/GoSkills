package trueskill

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/numerics"
)

type message struct {
	source node
	value  interface{}
}

type node struct {
	name string
	in   chan message
}

func (n node) String() string {
	return n.name
}

type variable struct {
	node
	value     interface{}
	neighbors []*factor
}

type factor struct {
	node
	neighbors []*variable
}

type factorGraph struct {
	vars    []variable
	factors []factor
}

func (fg *factorGraph) NewVar(name string, value interface{}) *variable {
	fg.vars = append(fg.vars, variable{
		node: node{
			name: name,
			in:   make(chan message),
		},
		value: value,
	})
	return &fg.vars[len(fg.vars)-1]
}

//////////////////

func (v variable) Value() interface{} {
	return v.value
}

//////////////////

func NewGaussianPriorFactor(mean, variance float64, v *variable) *factor {
	return nil
}

func NewGaussianLikelihoodFactor(beta float64, skill, perf *variable) *factor {
	// beta2 := numerics.Sqr(beta)
	return nil
}

func NewGaussianWeightedSumFactor(in []*variable, weights []float64, sum *variable) *factor {
	return nil
}

func NewGaussianWithinFactor(epsilon float64, diff *variable) *factor {
	return nil
}

func NewGaussianGreaterThanFactor(epsilon float64, diff *variable) *factor {
	return nil
}

//////////////////

type TrueSkillFactorGraph struct {
	fg    factorGraph
	vp    map[*variable]interface{}
	gi    *skills.GameInfo
	teams []skills.Team
	ranks []int
}

func NewTrueSkillFactorGraph(gi *skills.GameInfo, teams []skills.Team, ranks ...int) *TrueSkillFactorGraph {
	fg := &TrueSkillFactorGraph{
		vp:    make(map[*variable]interface{}),
		gi:    gi,
		teams: teams,
		ranks: ranks,
	}
	return fg
}

func (g *TrueSkillFactorGraph) BuildGraph() {
	fg := &g.fg
	// PlayerPriorValuesToSkillsLayer
	outputVariableGroups := [][]*variable{}
	for _, t := range g.teams {
		tVars := []*variable{}
		tFactors := []*factor{}
		for p, s := range t.PlayerRatings {
			pVar := fg.NewVar(fmt.Sprintf("%v's skill", p), numerics.NewGaussDist(g.gi.InitialMean, g.gi.InitialStddev))
			g.vp[pVar] = p
			tVars = append(tVars, pVar)
			pFactor := NewGaussianPriorFactor(s.Mean(), s.Variance(), pVar)
			tFactors = append(tFactors, pFactor)
		}
		outputVariableGroups = append(outputVariableGroups, tVars)
	}

	inputVariableGroups := outputVariableGroups
	outputVariableGroups = [][]*variable{}

	// PlayerSkillsToPerformancesLayer
	for _, tInVars := range inputVariableGroups {
		tVars := []*variable{}
		tFactors := []*factor{}
		for _, pInVar := range tInVars {
			p := g.vp[pInVar]
			pVar := fg.NewVar(fmt.Sprintf("%v's performance", p), numerics.NewGaussDist(g.gi.InitialMean, g.gi.InitialStddev))
			g.vp[pVar] = p
			tVars = append(tVars, pVar)
			pFactor := NewGaussianLikelihoodFactor(g.gi.Beta, pInVar, pVar)
			tFactors = append(tFactors, pFactor)
		}
		outputVariableGroups = append(outputVariableGroups, tVars)
	}

	inputVariableGroups = outputVariableGroups
	outputVariables := []*variable{}

	// PlayerPerformancesToTeamPerformancesLayer
	for ti, tInVars := range inputVariableGroups {
		sumVars := []*variable{}
		weights := []float64{}
		for _, pInVar := range tInVars {
			sumVars = append(sumVars, pInVar)
			weights = append(weights, 1.0)
		}
		teamPerfomance := fg.NewVar(fmt.Sprintf("Team %v's performance", g.teams[ti]), numerics.NewGaussDist(g.gi.InitialMean, g.gi.InitialStddev))
		outputVariables = append(outputVariables, teamPerfomance)
		// tFactor :=
		NewGaussianWeightedSumFactor(sumVars, weights, teamPerfomance)
	}

	inputVariables := outputVariables
	outputVariables = []*variable{}

	// TeamPerformancesToTeamPerformanceDifferencesLayer
	for ti := range inputVariables[1:] {
		diffVar := fg.NewVar(fmt.Sprintf("Team %v - %v performance difference", g.teams[ti], g.teams[ti+1]),
			numerics.NewGaussDist(g.gi.InitialMean, g.gi.InitialStddev))
		outputVariables = append(outputVariables, diffVar)
		// diffFactor :=
		NewGaussianWeightedSumFactor(inputVariables[ti:ti+2], []float64{1.0, -1.0}, diffVar)
	}

	inputVariables = outputVariables
	outputVariables = []*variable{}

	// TeamDifferencesComparisonLayer
	for ti, diffVar := range inputVariables {
		epsilon := drawMarginFromDrawProbability(g.gi.DrawProbability, g.gi.Beta)
		if g.ranks[ti] == g.ranks[ti+1] {
			// factor =
			NewGaussianWithinFactor(epsilon, diffVar)
		} else {
			// factor =
			NewGaussianGreaterThanFactor(epsilon, diffVar)
		}
	}
}

func (g *TrueSkillFactorGraph) RunSchedule() {
	fmt.Printf("Running, %v vars\n", len(g.fg.vars))
}

func (g *TrueSkillFactorGraph) GetUpdatedRatings() skills.PlayerRatings {
	return make(skills.PlayerRatings)
}
