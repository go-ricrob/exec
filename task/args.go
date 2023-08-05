package task

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const envPrefix = "RICROB_" // environment prefix.

func envVar(name string) string { return envPrefix + strings.ToUpper(name) }

func argVar(name string) string { return "-" + name }

func envString(name, def string) string {
	if value, ok := os.LookupEnv(envVar(name)); ok {
		return value
	}
	return def
}

func envBool(name string, def bool) bool {
	if strValue, ok := os.LookupEnv(envVar(name)); ok {
		if value, err := strconv.ParseBool(strValue); err == nil {
			return value
		}
	}
	return def
}

func envCoord(name string, def Coordinate) Coordinate {
	if strValue, ok := os.LookupEnv(envVar(name)); ok {
		if value, err := parseCoordinate(strValue); err == nil {
			return value
		}
	}
	return def
}

func envSymbol(name string, def Symbol) Symbol {
	if strValue, ok := os.LookupEnv(envVar(name)); ok {
		if value, err := parseSymbol(strValue); err == nil {
			return value
		}
	}
	return def
}

func usage(key, text string) string {
	return fmt.Sprintf("%s (environment variable %s)", text, envVar(key))
}

const (
	fnTargetSymbol       = "ts"
	fnCheckRobotOnSymbol = "crs"
)

var (
	defTargetSymbol       = Cosmic
	defCheckRobotOnSymbol = true
)

// Args holds the task arguments.
type Args struct {
	Tiles        *Tiles
	Robots       *Robots
	TargetSymbol Symbol

	CheckRobotOnSymbol bool
}

func newArgs() *Args {
	return &Args{Tiles: new(Tiles), Robots: new(Robots)}
}

// CmdArgs returns an argument slice build by task parameters.
func (a *Args) CmdArgs() []string {
	s := []string{
		argVar(fnTopLeftTile), a.Tiles.TopLeft,
		argVar(fnTopRightTile), a.Tiles.TopRight,
		argVar(fnBottomRightTile), a.Tiles.BottomRight,
		argVar(fnBottomLeftTile), a.Tiles.BottomLeft,

		argVar(fnYellowRobot), a.Robots.Yellow.String(),
		argVar(fnRedRobot), a.Robots.Red.String(),
		argVar(fnGreenRobot), a.Robots.Green.String(),
		argVar(fnBlueRobot), a.Robots.Blue.String(),
		argVar(fnSilverRobot), a.Robots.Silver.String(),

		argVar(fnTargetSymbol), a.TargetSymbol.String(),
	}
	if a.CheckRobotOnSymbol {
		s = append(s, argVar(fnCheckRobotOnSymbol))
	}
	return s
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
	a := newArgs()

	fs := flag.NewFlagSet(name, errorHandling)

	a.Tiles.addFlag(fs)
	a.Robots.addFlag(fs)

	var err error
	a.TargetSymbol = envSymbol(fnTargetSymbol, defTargetSymbol)
	fs.Func(fnTargetSymbol, usage(fnTargetSymbol, "target symbol like yellowPyramid or cosmic"), func(s string) error {
		a.TargetSymbol, err = parseSymbol(s)
		return err
	})

	fs.BoolVar(&a.CheckRobotOnSymbol, fnCheckRobotOnSymbol, envBool(fnCheckRobotOnSymbol, defCheckRobotOnSymbol), usage(fnCheckRobotOnSymbol, "check if robots sit on symbol"))

	fs.Parse(cmdArgs)

	if err := a.Tiles.check(); err != nil {
		return nil, err
	}
	if err := a.Robots.check(a.Tiles, a.CheckRobotOnSymbol); err != nil {
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

// ParseURL parses and returns arguments base on an URL query.
func ParseURL(u *url.URL) (*Args, error) {
	a := new(Args)

	query := u.Query()

	if err := a.Tiles.ParseURL(u); err != nil {
		return nil, err
	}

	var err error
	if a.CheckRobotOnSymbol, err = queryBool(query, fnCheckRobotOnSymbol, defCheckRobotOnSymbol); err != nil {
		return nil, err
	}

	if err := a.Robots.parseURL(u, a.Tiles, a.CheckRobotOnSymbol); err != nil {
		return nil, err
	}

	if a.TargetSymbol, err = querySymbol(query, fnTargetSymbol, defTargetSymbol); err != nil {
		return nil, err
	}

	return a, nil
}
