package trueskill

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"runtime"
	"testing"
)

const (
	errorTolerance = 0.085
)

func AllTwoPlayerScenarios(t *testing.T, calc skills.Calc) {
	TwoPlayerTestNotDrawn(t, calc)
	TwoPlayerTestDrawn(t, calc)
	OneOnOneMassiveUpsetDrawTest(t, calc)

	TwoPlayerChessTestNotDrawn(t, calc)
}

func AllTwoTeamScenarios(t *testing.T, calc skills.Calc) {
	OneOnTwoSimpleTest(t, calc)
	OneOnTwoSomewhatBalanced(t, calc)
	OneOnThreeSimpleTest(t, calc)
	OneOnTwoDrawTest(t, calc)
	OneOnThreeDrawTest(t, calc)
	OneOnSevenSimpleTest(t, calc)

	TwoOnTwoSimpleTest(t, calc)
	TwoOnTwoDrawTest(t, calc)
	TwoOnTwoUnbalancedDrawTest(t, calc)
	TwoOnTwoUpsetTest(t, calc)

	ThreeOnTwoTests(t, calc)

	FourOnFourSimpleTest(t, calc)
}

func AllMultipleTeamScenarios(t *testing.T, calc skills.Calc) {
	ThreeTeamsOfOneNotDrawn(t, calc)
	// ThreeTeamsOfOneDrawn(t, calc)
	// FourTeamsOfOneNotDrawn(t, calc)
	// FiveTeamsOfOneNotDrawn(t, calc)
	// EightTeamsOfOneDrawn(t, calc)
	// EightTeamsOfOneUpset(t, calc)
	// SixteenTeamsOfOneNotDrawn(t, calc)

	// TwoOnFourOnTwoWinDraw(t, calc)
}

//------------------- Actual Tests ---------------------------
// If you see more than 3 digits of precision in the decimal point, then the expected values calculated from
// F# RalfH's implementation with the same input. It didn't support teams, so team values all came from the
// online calculator at http://atom.research.microsoft.com/trueskill/rankcalculator.aspx
//
// All match quality expected values came from the online calculator

// In both cases, there may be some discrepancy after the first decimal point. I think this is due to my implementation
// using slightly higher precision in GaussianDistribution.

//------------------------------------------------------------------------------
// Two Player Tests
//------------------------------------------------------------------------------

func TwoPlayerTestNotDrawn(t *testing.T, calc skills.Calc) {
	Convey("Using default game info", t, func() {
		gameInfo := skills.DefaultGameInfo
		Convey("Given two rookie teams of one", func() {
			teams := []skills.Team{skills.NewTeam(), skills.NewTeam()}

			teams[0].AddPlayer(2, gameInfo.DefaultRating())
			teams[1].AddPlayer(1, gameInfo.DefaultRating())

			Convey("The match quality should be 0.447", func() {
				So(calc.CalcMatchQual(gameInfo, teams), ShouldAlmostEqual, 0.447, 0.001)
			})

			Convey("Given team one wins", func() {
				ranks := []int{2, 1}

				Convey("Calculate the new ratings", func() {
					newRatings := calc.CalcNewRatings(gameInfo, teams, ranks...)

					Convey("The passed slices should not be reordered", func() {
						So(teams[0].Players()[0], ShouldEqual, 2)
						So(teams[1].Players()[0], ShouldEqual, 1)
						So(ranks, ShouldResemble, []int{2, 1})
					})

					Convey("Team one's mean should be higher and the standard deviations should have gone down.", func() {
						So(newRatings[1].Mean(), ShouldAlmostEqual, 29.39583201999924, errorTolerance)
						So(newRatings[1].Stddev(), ShouldAlmostEqual, 7.171475587326186, errorTolerance)

						So(newRatings[2].Mean(), ShouldAlmostEqual, 20.60416798000076, errorTolerance)
						So(newRatings[2].Stddev(), ShouldAlmostEqual, 7.171475587326186, errorTolerance)
					})
				})
			})
		})
	})
}

func TwoPlayerTestDrawn(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	player1NewRating := newRatings[1]
	AssertRating(t, 25.0, 6.4575196623173081, player1NewRating)

	player2NewRating := newRatings[2]
	AssertRating(t, 25.0, 6.4575196623173081, player2NewRating)
	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoPlayerChessTestNotDrawn(t *testing.T, calc skills.Calc) {
	// Inspired by a real bug :-)
	gameInfo := &skills.GameInfo{
		InitialMean:     1200,
		InitialStddev:   1200 / 3,
		Beta:            200,
		DynamicsFactor:  1200 / 300,
		DrawProbability: 0.03,
	}

	team1 := skills.NewTeam()
	team1.AddPlayer(1, skills.NewRating(1301.0007, 42.9232))

	team2 := skills.NewTeam()
	team2.AddPlayer(2, skills.NewRating(1188.7560, 42.5570))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	player1NewRating := newRatings[1]
	AssertRating(t, 1304.7820836053318, 42.843513887848658, player1NewRating)

	player2NewRating := newRatings[2]
	AssertRating(t, 1185.0383099003536, 42.485604606897752, player2NewRating)
}

func OneOnOneMassiveUpsetDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, skills.NewRating(50, 12.5))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	player1NewRating := newRatings[1]
	AssertRating(t, 31.662, 7.137, player1NewRating)

	player2NewRating := newRatings[2]
	AssertRating(t, 35.010, 7.910, player2NewRating)
	AssertMatchQuality(t, 0.110, calc.CalcMatchQual(gameInfo, teams))
}

//------------------------------------------------------------------------------
// Two Team Tests
//------------------------------------------------------------------------------

func TwoOnTwoSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())
	team1.AddPlayer(2, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(3, gameInfo.DefaultRating())
	team2.AddPlayer(4, gameInfo.DefaultRating())

	teams := []skills.Team{team2, team1}

	ranks := []int{2, 1}
	newRatings := calc.CalcNewRatings(gameInfo, teams, ranks...)
	if ranks[0] != 2 {
		t.Errorf("client ranks slice reordered")
	}

	// Winners
	AssertRating(t, 28.108, 7.774, newRatings[1])
	AssertRating(t, 28.108, 7.774, newRatings[2])

	// Losers
	AssertRating(t, 21.892, 7.774, newRatings[3])
	AssertRating(t, 21.892, 7.774, newRatings[4])

	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoOnTwoDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())
	team1.AddPlayer(2, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(3, gameInfo.DefaultRating())
	team2.AddPlayer(4, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 25, 7.455, newRatings[1])
	AssertRating(t, 25, 7.455, newRatings[2])

	// Losers
	AssertRating(t, 25, 7.445, newRatings[3])
	AssertRating(t, 25, 7.445, newRatings[4])

	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoOnTwoUnbalancedDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, skills.NewRating(15, 8))
	team1.AddPlayer(2, skills.NewRating(20, 6))

	team2 := skills.NewTeam()
	team2.AddPlayer(3, skills.NewRating(25, 4))
	team2.AddPlayer(4, skills.NewRating(30, 3))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 21.570, 6.556, newRatings[1])
	AssertRating(t, 23.696, 5.418, newRatings[2])

	// Losers
	AssertRating(t, 23.357, 3.833, newRatings[3])
	AssertRating(t, 29.075, 2.931, newRatings[4])

	AssertMatchQuality(t, 0.214, calc.CalcMatchQual(gameInfo, teams))
}

func TwoOnTwoUpsetTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, skills.NewRating(20, 8))
	team1.AddPlayer(2, skills.NewRating(25, 6))

	team2 := skills.NewTeam()
	team2.AddPlayer(3, skills.NewRating(35, 7))
	team2.AddPlayer(4, skills.NewRating(40, 5))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 29.698, 7.008, newRatings[1])
	AssertRating(t, 30.455, 5.594, newRatings[2])

	// Losers
	AssertRating(t, 27.575, 6.346, newRatings[3])
	AssertRating(t, 36.211, 4.768, newRatings[4])

	AssertMatchQuality(t, 0.084, calc.CalcMatchQual(gameInfo, teams))
}

func FourOnFourSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())
	team1.AddPlayer(2, gameInfo.DefaultRating())
	team1.AddPlayer(3, gameInfo.DefaultRating())
	team1.AddPlayer(4, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(5, gameInfo.DefaultRating())
	team2.AddPlayer(6, gameInfo.DefaultRating())
	team2.AddPlayer(7, gameInfo.DefaultRating())
	team2.AddPlayer(8, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 27.198, 8.059, newRatings[1])
	AssertRating(t, 27.198, 8.059, newRatings[2])
	AssertRating(t, 27.198, 8.059, newRatings[3])
	AssertRating(t, 27.198, 8.059, newRatings[4])

	// Losers
	AssertRating(t, 22.802, 8.059, newRatings[5])
	AssertRating(t, 22.802, 8.059, newRatings[6])
	AssertRating(t, 22.802, 8.059, newRatings[7])
	AssertRating(t, 22.802, 8.059, newRatings[8])

	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnTwoSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, gameInfo.DefaultRating())
	team2.AddPlayer(3, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 33.730, 7.317, newRatings[1])

	// Losers
	AssertRating(t, 16.270, 7.317, newRatings[2])
	AssertRating(t, 16.270, 7.317, newRatings[3])

	AssertMatchQuality(t, 0.135, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnTwoSomewhatBalanced(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, skills.NewRating(40, 6))

	team2 := skills.NewTeam()
	team2.AddPlayer(2, skills.NewRating(20, 7))
	team2.AddPlayer(3, skills.NewRating(25, 8))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 42.744, 5.602, newRatings[1])

	// Losers
	AssertRating(t, 16.266, 6.359, newRatings[2])
	AssertRating(t, 20.123, 7.028, newRatings[3])

	AssertMatchQuality(t, 0.478, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnThreeSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, gameInfo.DefaultRating())
	team2.AddPlayer(3, gameInfo.DefaultRating())
	team2.AddPlayer(4, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 36.337, 7.527, newRatings[1])

	// Losers
	AssertRating(t, 13.663, 7.527, newRatings[2])
	AssertRating(t, 13.663, 7.527, newRatings[3])
	AssertRating(t, 13.663, 7.527, newRatings[4])

	AssertMatchQuality(t, 0.012, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnTwoDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, gameInfo.DefaultRating())
	team2.AddPlayer(3, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 31.660, 7.138, newRatings[1])

	// Losers
	AssertRating(t, 18.340, 7.138, newRatings[2])
	AssertRating(t, 18.340, 7.138, newRatings[3])

	AssertMatchQuality(t, 0.135, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnThreeDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, gameInfo.DefaultRating())
	team2.AddPlayer(3, gameInfo.DefaultRating())
	team2.AddPlayer(4, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 34.990, 7.455, newRatings[1])

	// Losers
	AssertRating(t, 15.010, 7.455, newRatings[2])
	AssertRating(t, 15.010, 7.455, newRatings[3])
	AssertRating(t, 15.010, 7.455, newRatings[4])

	AssertMatchQuality(t, 0.012, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnSevenSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(2, gameInfo.DefaultRating())
	team2.AddPlayer(3, gameInfo.DefaultRating())
	team2.AddPlayer(4, gameInfo.DefaultRating())
	team2.AddPlayer(5, gameInfo.DefaultRating())
	team2.AddPlayer(6, gameInfo.DefaultRating())
	team2.AddPlayer(7, gameInfo.DefaultRating())
	team2.AddPlayer(8, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 40.582, 7.917, newRatings[1])

	// Losers
	AssertRating(t, 9.418, 7.917, newRatings[2])
	AssertRating(t, 9.418, 7.917, newRatings[3])
	AssertRating(t, 9.418, 7.917, newRatings[4])
	AssertRating(t, 9.418, 7.917, newRatings[5])
	AssertRating(t, 9.418, 7.917, newRatings[6])
	AssertRating(t, 9.418, 7.917, newRatings[7])
	AssertRating(t, 9.418, 7.917, newRatings[8])

	AssertMatchQuality(t, 0.000, calc.CalcMatchQual(gameInfo, teams))
}

func ThreeOnTwoTests(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(1, skills.NewRating(28, 7))
	team1.AddPlayer(2, skills.NewRating(27, 6))
	team1.AddPlayer(3, skills.NewRating(26, 5))

	team2 := skills.NewTeam()
	team2.AddPlayer(4, skills.NewRating(30, 4))
	team2.AddPlayer(5, skills.NewRating(31, 3))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 28.658, 6.770, newRatings[1])
	AssertRating(t, 27.484, 5.856, newRatings[2])
	AssertRating(t, 26.336, 4.917, newRatings[3])

	// Losers
	AssertRating(t, 29.785, 3.958, newRatings[4])
	AssertRating(t, 30.879, 2.983, newRatings[5])

	newRatings = calc.CalcNewRatings(gameInfo, teams, 2, 1)

	// Losers
	AssertRating(t, 21.840, 6.314, newRatings[1])
	AssertRating(t, 22.474, 5.575, newRatings[2])
	AssertRating(t, 22.857, 4.757, newRatings[3])

	// Winners
	AssertRating(t, 32.012, 3.877, newRatings[4])
	AssertRating(t, 32.132, 2.949, newRatings[5])

	AssertMatchQuality(t, 0.254, calc.CalcMatchQual(gameInfo, teams))
}

func ThreeTeamsOfOneNotDrawn(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())

	team3 := skills.NewTeam()
	team3.AddPlayer(*player3, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2, team3}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2, 3)

	AssertRating(t, 31.675352419172107, 6.6559853776206905, newRatings[*player1])
	AssertRating(t, 25.000000000003912, 6.2078966412243233, newRatings[*player2])
	AssertRating(t, 18.324647580823971, 6.6559853776218318, newRatings[*player3])

	AssertMatchQuality(t, 0.200, calc.CalcMatchQual(gameInfo, teams))
}

func testLoc() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%v:%v", file, line)
	}
	return ""
}

func AssertRating(t *testing.T, expectedMean, expectedStddev float64, actual skills.Rating) {
	if r := actual.Mean(); math.Abs(r-expectedMean) > errorTolerance {
		t.Errorf("actual.Mean = %v, want %v\n%v", r, expectedMean, testLoc())
	}
	if r := actual.Stddev(); math.Abs(r-expectedStddev) > errorTolerance {
		t.Errorf("actual.Stddev = %v, want %v\n%v", r, expectedStddev, testLoc())
	}
}
func AssertMatchQuality(t *testing.T, expectedMatchQual, actualMatchQual float64) {
	if r := actualMatchQual; math.Abs(r-expectedMatchQual) > errorTolerance {
		t.Errorf("actualMatchQual = %v, want %v\n%v", r, expectedMatchQual, testLoc())
	}
}
