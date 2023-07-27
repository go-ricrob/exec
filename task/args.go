package task

import (
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Parameter name constants.
const (
	fnTargetSymbol       = "ts"
	fnCheckRobotOnSymbol = "crs"
)

// Default parameters.
var (
	defTargetSymbol       = Cosmic
	defCheckRobotOnSymbol = true
)

// Args holds the task arguments.
type Args struct {
	Tiles        Tiles
	Robots       Robots
	TargetSymbol Symbol

	CheckRobotOnSymbol bool
}

// CmdArgs returns an argument slice build by task parameters.
func (a *Args) CmdArgs() []string {
	return []string{
		fmt.Sprintf("-%s %s", fnTopLeftTile, a.Tiles.TopLeft),
		fmt.Sprintf("-%s %s", fnTopRightTile, a.Tiles.TopRight),
		fmt.Sprintf("-%s %s", fnBottomRightTile, a.Tiles.BottomRight),
		fmt.Sprintf("-%s %s", fnBottomLeftTile, a.Tiles.BottomLeft),

		fmt.Sprintf("-%s %d,%d", fnYellowRobot, a.Robots.Yellow.X, a.Robots.Yellow.Y),
		fmt.Sprintf("-%s %d,%d", fnRedRobot, a.Robots.Red.X, a.Robots.Red.Y),
		fmt.Sprintf("-%s %d,%d", fnGreenRobot, a.Robots.Green.X, a.Robots.Green.Y),
		fmt.Sprintf("-%s %d,%d", fnBlueRobot, a.Robots.Blue.X, a.Robots.Blue.Y),
		fmt.Sprintf("-%s %d,%d", fnSilverRobot, a.Robots.Silver.X, a.Robots.Silver.Y),

		fmt.Sprintf("-%s %s", fnTargetSymbol, a.TargetSymbol),

		fmt.Sprintf("-%s %t", fnCheckRobotOnSymbol, a.CheckRobotOnSymbol),
	}
}

func parseSymbol(s string) (Symbol, error) {
	symbol, ok := symbolMap[s]
	if !ok {
		return symbol, fmt.Errorf("invalid symbol %s", s)
	}
	return symbol, nil
}

func parseCoordinate(s string) (Coordinate, error) {
	coord := Coordinate{}
	p := strings.Split(s, ",")
	if len(p) != 2 {
		return coord, fmt.Errorf("invalid coordinate format: %s", s)
	}
	x64, err := strconv.ParseInt(p[0], 10, 8)
	if err != nil {
		return coord, fmt.Errorf("invalid x coordinate %s - %w", s, err)
	}
	y64, err := strconv.ParseInt(p[1], 10, 8)
	if err != nil {
		return coord, fmt.Errorf("invalid y coordinate %s - %w", s, err)
	}
	coord.X, coord.Y = int(x64), int(y64)
	return coord, nil
}

func parseFlag(name string, cmdArgs []string, errorHandling flag.ErrorHandling) (*Args, error) {
	a := new(Args)

	fs := flag.NewFlagSet(name, errorHandling)

	a.Tiles.addFlag(fs)
	a.Robots.addFlag(fs)

	var err error
	a.TargetSymbol = defTargetSymbol
	fs.Func(fnTargetSymbol, "target symbol like yellowPyramid or cosmic", func(s string) error {
		a.TargetSymbol, err = parseSymbol(s)
		return err
	})

	fs.BoolVar(&a.CheckRobotOnSymbol, fnCheckRobotOnSymbol, defCheckRobotOnSymbol, "check if robots sit on symbol")

	fs.Parse(cmdArgs)

	if err := a.Tiles.check(); err != nil {
		return nil, err
	}
	if err := a.Robots.check(&a.Tiles, a.CheckRobotOnSymbol); err != nil {
		return nil, err
	}
	return a, nil
}

func queryString(values url.Values, name string, value string) string {
	if v, ok := values[name]; ok {
		return v[0]
	}
	return value
}

func querySymbol(values url.Values, name string, value Symbol) (Symbol, error) {
	if v, ok := values[name]; ok {
		return parseSymbol(v[0])
	}
	return value, nil
}

func queryCoordinate(values url.Values, name string, value Coordinate) (Coordinate, error) {
	if v, ok := values[name]; ok {
		return parseCoordinate(v[0])
	}
	return value, nil
}

func queryBool(values url.Values, name string, value bool) (bool, error) {
	if v, ok := values[name]; ok {
		return strconv.ParseBool(v[0])
	}
	return value, nil
}

func parseURL(u *url.URL) (*Args, error) {
	a := new(Args)

	query := u.Query()

	if err := a.Tiles.ParseURL(u); err != nil {
		return nil, err
	}

	var err error
	if a.CheckRobotOnSymbol, err = queryBool(query, fnCheckRobotOnSymbol, defCheckRobotOnSymbol); err != nil {
		return nil, err
	}

	if err := a.Robots.parseURL(u, &a.Tiles, a.CheckRobotOnSymbol); err != nil {
		return nil, err
	}

	if a.TargetSymbol, err = querySymbol(query, fnTargetSymbol, defTargetSymbol); err != nil {
		return nil, err
	}

	return a, nil
}
