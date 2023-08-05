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

func usage(name, text string) string {
	return fmt.Sprintf("%s [environment variable %s]", text, envVar(name))
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

func parseFlag(name string, cmdArgs []string, errorHandling flag.ErrorHandling) (*Args, error) {
	a := newArgs()

	fs := flag.NewFlagSet(name, errorHandling)

	fs.StringVar(&a.Tiles.TopLeft, fnTopLeftTile, envString(fnTopLeftTile, defTopLeftTile), usage(fnTopLeftTile, "top left tile"))
	fs.StringVar(&a.Tiles.TopRight, fnTopRightTile, envString(fnTopRightTile, defTopRightTile), usage(fnTopRightTile, "top right tile"))
	fs.StringVar(&a.Tiles.BottomLeft, fnBottomLeftTile, envString(fnBottomLeftTile, defBottomLeftTile), usage(fnBottomLeftTile, "top bottom left tile"))
	fs.StringVar(&a.Tiles.BottomRight, fnBottomRightTile, envString(fnBottomRightTile, defBottomRightTile), usage(fnBottomRightTile, "bottom right tile"))

	fs.TextVar(&a.Robots.Yellow, fnYellowRobot, envCoord(fnYellowRobot, defYellowRobot), usage(fnYellowRobot, "yellow robot position x,y"))
	fs.TextVar(&a.Robots.Red, fnRedRobot, envCoord(fnRedRobot, defRedRobot), usage(fnRedRobot, "red robot position x,y"))
	fs.TextVar(&a.Robots.Green, fnGreenRobot, envCoord(fnGreenRobot, defGreenRobot), usage(fnGreenRobot, "green robot position x,y"))
	fs.TextVar(&a.Robots.Blue, fnBlueRobot, envCoord(fnBlueRobot, defBlueRobot), usage(fnBlueRobot, "blue robot position x,y"))
	fs.TextVar(&a.Robots.Silver, fnSilverRobot, envCoord(fnSilverRobot, defSilverRobot), usage(fnSilverRobot, "silver robot position x,y"))

	fs.TextVar(&a.TargetSymbol, fnTargetSymbol, envSymbol(fnTargetSymbol, defTargetSymbol), usage(fnTargetSymbol, "target symbol like yellowPyramid or cosmic"))
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
