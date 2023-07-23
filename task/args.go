package task

import (
	"flag"
	"fmt"
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

	prmNoSymbolCheck = "s"
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

	defNoSymbolCheck = false
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

var symbolMap = map[string]Symbol{
	"yellowPyramid": YellowPyramid,
	"yellowStar":    YellowStar,
	"yellowMoon":    YellowMoon,
	"yellowSaturn":  YellowSaturn,

	"redPyramid": RedPyramid,
	"redStar":    RedStar,
	"redMoon":    RedMoon,
	"redSaturn":  RedSaturn,

	"greenPyramid": GreenPyramid,
	"greenStar":    GreenStar,
	"greenMoon":    GreenMoon,
	"greenSaturn":  GreenSaturn,

	"bluePyramid": BluePyramid,
	"blueStar":    BlueStar,
	"blueMoon":    BlueMoon,
	"blueSaturn":  BlueSaturn,

	"cosmic": Cosmic,
}

func parseSymbol(s string, ptr *Symbol) error {
	symbol, ok := symbolMap[s]
	if !ok {
		return fmt.Errorf("invalid symbol %s", s)
	}
	*ptr = symbol
	return nil
}

// Check flags for arguments
const (
	NoSymbolCheck = 1 << iota
)

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
}

// HasSilverRobot returns true if the silver robot is used, else otherwise.
func (a *Args) HasSilverRobot() bool { return a.SilverRobot.X != -1 && a.SilverRobot.Y != -1 }

// Check checks validity and consistency of arguments.
func (a *Args) Check(flag int) error {
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
		if flag&NoSymbolCheck == 0 {
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

func parseArgs(name string, cmdArgs []string, errorHandling flag.ErrorHandling) (*Args, error) {
	a := new(Args)

	fs := flag.NewFlagSet(name, errorHandling)

	var noSymbolCheck bool
	fs.BoolVar(&noSymbolCheck, prmNoSymbolCheck, defNoSymbolCheck, "do not check if robots sit on symbol")

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

	var flag int
	if noSymbolCheck {
		flag += NoSymbolCheck
	}

	if err := a.Check(flag); err != nil {
		return nil, err
	}
	return a, nil
}

/*
func makeFlag(prm string) string { return "-" + prm }

// Args returns an argument slice build by task parameters.
func Args(t *task.Task) []string {
	tiles := t.Tiles()
	robots := t.Robots()
	return []string{
		makeFlag(game.TopLeftTilePrm), tiles[game.TopLeft].String(),
		makeFlag(game.TopRightTilePrm), tiles[game.TopRight].String(),
		makeFlag(game.BottomRightTilePrm), tiles[game.BottomRight].String(),
		makeFlag(game.BottomLeftTilePrm), tiles[game.BottomLeft].String(),

		makeFlag(game.YellowRobotPrm), robots[game.Yellow].String(),
		makeFlag(game.RedRobotPrm), robots[game.Red].String(),
		makeFlag(game.GreenRobotPrm), robots[game.Green].String(),
		makeFlag(game.BlueRobotPrm), robots[game.Blue].String(),
		makeFlag(game.SilverRobotPrm), robots[game.Silver].String(),

		makeFlag(game.TargetPrm), t.Target().String(),
	}
}
*/
