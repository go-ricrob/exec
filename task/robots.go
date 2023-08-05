package task

import (
	"fmt"
	"net/url"

	"github.com/go-ricrob/game/board"
)

// Robots flag name constants.
const (
	fnYellowRobot = "ry"
	fnRedRobot    = "rr"
	fnGreenRobot  = "rg"
	fnBlueRobot   = "rb"
	fnSilverRobot = "rs"
)

// Default parameters.
var (
	defYellowRobot = Coordinate{X: 0, Y: 0}
	defRedRobot    = Coordinate{X: 1, Y: 0}
	defGreenRobot  = Coordinate{X: 2, Y: 0}
	defBlueRobot   = Coordinate{X: 3, Y: 0}
	defSilverRobot = Coordinate{X: -1, Y: -1}
)

// Robots represents the coordinates of the robots.
type Robots struct {
	Yellow, Red, Green, Blue, Silver Coordinate
}

// HasSilver returns true if the silver robot is used, else otherwise.
func (r *Robots) HasSilver() bool { return r.Silver.X != -1 && r.Silver.Y != -1 }

func (r *Robots) check(t *Tiles, checkRobotOnSymbol bool) error {
	b := board.New([board.NumTile]string{
		board.TopLeft:     t.TopLeft,
		board.TopRight:    t.TopRight,
		board.BottomLeft:  t.BottomLeft,
		board.BottomRight: t.BottomRight,
	})

	// check robots
	robots := []Coordinate{r.Yellow, r.Red, r.Green, r.Blue}
	if r.Silver.X != -1 && r.Silver.Y != -1 {
		robots = append(robots, r.Silver)
	}
	robotMap := map[Coordinate]bool{}
	for _, robot := range robots {
		if !b.IsValidCoordinate(robot.X, robot.Y) {
			return fmt.Errorf("invalid robot coordinates %v - center field", robot)
		}
		if checkRobotOnSymbol {
			field := b.Field(robot.X, robot.Y)
			if field.Symbol != board.NoSymbol {
				return fmt.Errorf("robot %v sits on symbol %s color %s", robot, field.Symbol, field.Color)
			}
		}
		if _, ok := robotMap[robot]; ok {
			return fmt.Errorf("duplicate robot position %v", robot)
		}
		robotMap[robot] = true
	}
	return nil
}

func (r *Robots) parseURL(u *url.URL, t *Tiles, checkRobotOnSymbol bool) error {
	var err error
	query := u.Query()
	if r.Yellow, err = queryCoordinate(query, fnYellowRobot, defYellowRobot); err != nil {
		return err
	}
	if r.Red, err = queryCoordinate(query, fnRedRobot, defRedRobot); err != nil {
		return err
	}
	if r.Green, err = queryCoordinate(query, fnGreenRobot, defGreenRobot); err != nil {
		return err
	}
	if r.Blue, err = queryCoordinate(query, fnBlueRobot, defBlueRobot); err != nil {
		return err
	}
	if r.Silver, err = queryCoordinate(query, fnSilverRobot, defSilverRobot); err != nil {
		return err
	}

	if err := r.check(t, checkRobotOnSymbol); err != nil {
		return err
	}
	return nil
}
