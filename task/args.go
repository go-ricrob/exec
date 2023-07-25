package task

import (
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-ricrob/game/board"
)

// Parameter name constants.
const (
	prmTopLeftTile     = "ttl"
	prmTopRightTile    = "ttr"
	prmBottomRightTile = "tbr"
	prmBottomLeftTile  = "tbl"

	prmYellowRobot = "ry"
	prmRedRobot    = "rr"
	prmGreenRobot  = "rg"
	prmBlueRobot   = "rb"
	prmSilverRobot = "rs"

	prmTargetSymbol = "ts"

	prmCheckRobotOnSymbol = "crs"
)

// Default parameters.
var (
	defTopLeftTile     = "A1F"
	defTopRightTile    = "A2F"
	defBottomRightTile = "A3F"
	defBottomLeftTile  = "A4F"

	defYellowRobot = "0,0"
	defRedRobot    = "1,0"
	defGreenRobot  = "2,0"
	defBlueRobot   = "3,0"
	defSilverRobot = "-1,-1"

	defTargetSymbol = "cosmic"

	defCheckRobotOnSymbol = true
)

func parseCoordinate(s string, coord *Coordinate) error {
	p := strings.Split(s, ",")
	if len(p) != 2 {
		return fmt.Errorf("invalid coordinate format: %s", s)
	}
	x64, err := strconv.ParseInt(p[0], 10, 8)
	if err != nil {
		return fmt.Errorf("invalid x coordinate %s - %w", s, err)
	}
	y64, err := strconv.ParseInt(p[1], 10, 8)
	if err != nil {
		return fmt.Errorf("invalid y coordinate %s - %w", s, err)
	}

	coord.X, coord.Y = int(x64), int(y64)
	return nil
}

func parseSymbol(s string, ptr *Symbol) error {
	symbol, ok := symbolMap[s]
	if !ok {
		return fmt.Errorf("invalid symbol %s", s)
	}
	*ptr = symbol
	return nil
}

// Args holds the task arguments.
type Args struct {
	TopLeftTile     string
	TopRightTile    string
	BottomLeftTile  string
	BottomRightTile string

	YellowRobot Coordinate
	RedRobot    Coordinate
	GreenRobot  Coordinate
	BlueRobot   Coordinate
	SilverRobot Coordinate

	TargetSymbol Symbol

	CheckRobotOnSymbol bool
}

// HasSilverRobot returns true if the silver robot is used, else otherwise.
func (a *Args) HasSilverRobot() bool { return a.SilverRobot.X != -1 && a.SilverRobot.Y != -1 }

// CmdArgs returns an argument slice build by task parameters.
func (a *Args) CmdArgs() []string {
	return []string{
		fmt.Sprintf("-%s %s", prmTopLeftTile, a.TopLeftTile),
		fmt.Sprintf("-%s %s", prmTopRightTile, a.TopRightTile),
		fmt.Sprintf("-%s %s", prmBottomRightTile, a.BottomRightTile),
		fmt.Sprintf("-%s %s", prmBottomLeftTile, a.BottomLeftTile),

		fmt.Sprintf("-%s %d,%d", prmYellowRobot, a.YellowRobot.X, a.YellowRobot.Y),
		fmt.Sprintf("-%s %d,%d", prmRedRobot, a.RedRobot.X, a.RedRobot.Y),
		fmt.Sprintf("-%s %d,%d", prmGreenRobot, a.GreenRobot.X, a.GreenRobot.Y),
		fmt.Sprintf("-%s %d,%d", prmBlueRobot, a.BlueRobot.X, a.BlueRobot.Y),
		fmt.Sprintf("-%s %d,%d", prmSilverRobot, a.SilverRobot.X, a.SilverRobot.Y),

		fmt.Sprintf("-%s %s", prmTargetSymbol, a.TargetSymbol),

		fmt.Sprintf("-%s %t", prmCheckRobotOnSymbol, a.CheckRobotOnSymbol),
	}
}

// Check checks validity and consistency of arguments.
func (a *Args) Check() error {
	// check tiles
	tiles := []string{a.TopLeftTile, a.TopRightTile, a.BottomLeftTile, a.BottomRightTile}
	tileMap := map[string]bool{}
	for _, tile := range tiles {
		if (len(tile)) != 3 || (tile[0] != 'A' && tile[0] != 'B') || (tile[1] < '1' || tile[1] > '4') || (tile[2] != 'F' && tile[2] != 'B') {
			return fmt.Errorf("invalid tile %s", tile)
		}
		if _, ok := tileMap[tile]; ok {
			return fmt.Errorf("duplicate tile %s", tile)
		}
		tileMap[tile] = true
	}

	b := board.New([board.NumTile]string{
		board.TopLeft:     a.TopLeftTile,
		board.TopRight:    a.TopRightTile,
		board.BottomLeft:  a.BottomLeftTile,
		board.BottomRight: a.BottomRightTile,
	})

	// check robots
	robots := []Coordinate{a.YellowRobot, a.RedRobot, a.GreenRobot, a.BlueRobot}
	if a.SilverRobot.X != -1 && a.SilverRobot.Y != -1 {
		robots = append(robots, a.SilverRobot)
	}
	robotMap := map[Coordinate]bool{}
	for _, robot := range robots {
		if !b.IsValidCoordinate(robot.X, robot.Y) {
			return fmt.Errorf("invalid robot coordinates %v - center field", robot)
		}
		if a.CheckRobotOnSymbol {
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

func parseFlag(name string, cmdArgs []string, errorHandling flag.ErrorHandling) (*Args, error) {
	a := new(Args)

	fs := flag.NewFlagSet(name, errorHandling)

	fs.StringVar(&a.TopLeftTile, prmTopLeftTile, defTopLeftTile, "top left tile")
	fs.StringVar(&a.TopRightTile, prmTopRightTile, defTopRightTile, "top right tile")
	fs.StringVar(&a.BottomLeftTile, prmBottomLeftTile, defBottomLeftTile, "top bottom left tile")
	fs.StringVar(&a.BottomRightTile, prmBottomRightTile, defBottomRightTile, "bottom right tile")

	var ry, rr, rg, rb, rs string

	fs.StringVar(&ry, prmYellowRobot, defYellowRobot, "yellow robot position x,y")
	fs.StringVar(&rr, prmRedRobot, defRedRobot, "red robot position x,y")
	fs.StringVar(&rg, prmGreenRobot, defGreenRobot, "green robot position x,y")
	fs.StringVar(&rb, prmBlueRobot, defBlueRobot, "blue robot position x,y")
	fs.StringVar(&rs, prmSilverRobot, defSilverRobot, "silver robot position x,y")

	var ts string
	fs.StringVar(&ts, prmTargetSymbol, defTargetSymbol, "target symbol like yellowPyramid or cosmic")

	fs.BoolVar(&a.CheckRobotOnSymbol, prmCheckRobotOnSymbol, defCheckRobotOnSymbol, "check if robots sit on symbol")

	fs.Parse(cmdArgs)

	if err := parseCoordinate(ry, &a.YellowRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rr, &a.RedRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rg, &a.GreenRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rb, &a.BlueRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rs, &a.SilverRobot); err != nil {
		return nil, err
	}

	if err := parseSymbol(ts, &a.TargetSymbol); err != nil {
		return nil, err
	}

	if err := a.Check(); err != nil {
		return nil, err
	}
	return a, nil
}

type querySet struct {
	values url.Values
}

func newQuerySet(u *url.URL) *querySet { return &querySet{values: u.Query()} }

func (qs *querySet) string(name string, value string) string {
	if v, ok := qs.values[name]; ok {
		return v[0]
	}
	return value
}

func (qs *querySet) bool(name string, value bool) bool {
	if v, ok := qs.values[name]; ok {
		if b, err := strconv.ParseBool(v[0]); err == nil {
			return b
		}
	}
	return value
}

func parseURL(u *url.URL) (*Args, error) {
	a := new(Args)

	qs := newQuerySet(u)

	a.TopLeftTile = qs.string(prmTopLeftTile, defTopLeftTile)
	a.TopRightTile = qs.string(prmTopRightTile, defTopRightTile)
	a.BottomLeftTile = qs.string(prmBottomLeftTile, defBottomLeftTile)
	a.BottomRightTile = qs.string(prmBottomRightTile, defBottomRightTile)

	if err := parseCoordinate(qs.string(prmYellowRobot, defYellowRobot), &a.YellowRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(qs.string(prmRedRobot, defRedRobot), &a.RedRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(qs.string(prmGreenRobot, defGreenRobot), &a.GreenRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(qs.string(prmBlueRobot, defBlueRobot), &a.BlueRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(qs.string(prmSilverRobot, defSilverRobot), &a.SilverRobot); err != nil {
		return nil, err
	}

	if err := parseSymbol(qs.string(prmTargetSymbol, defTargetSymbol), &a.TargetSymbol); err != nil {
		return nil, err
	}

	a.CheckRobotOnSymbol = qs.bool(prmCheckRobotOnSymbol, defCheckRobotOnSymbol)

	if err := a.Check(); err != nil {
		return nil, err
	}

	return a, nil
}
