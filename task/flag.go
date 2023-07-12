package task

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	prmTargetColor  = "tc"
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

	defTargetSymbol = "Cosmic"
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
	"Pyramid": Pyramid,
	"Star":    Star,
	"Moon":    Moon,
	"Saturn":  Saturn,
	"Cosmic":  Cosmic,
}

var targetColorMap = map[string]Color{
	"yellow": Yellow,
	"red":    Red,
	"green":  Green,
	"blue":   Blue,
}

func parseSymbol(s string, ptr *Symbol) error {
	symbol, ok := symbolMap[s]
	if !ok {
		return fmt.Errorf("invalid symbol %s", s)
	}
	*ptr = symbol
	return nil
}

func parseColor(s string, ptr *Color) error {
	color, ok := targetColorMap[s]
	if !ok {
		return fmt.Errorf("invalid color %s", s)
	}
	*ptr = color
	return nil
}

func parseFlags(name string, cmdArgs []string, errorHandling flag.ErrorHandling) (*Task, error) {

	t := newTask()

	fs := flag.NewFlagSet(name, errorHandling)

	fs.StringVar(&t.Args.TopLeftTile, prmTopLeftTile, defTopLeftTile, "top left tile")
	fs.StringVar(&t.Args.TopRightTile, prmTopRightTile, defTopRightTile, "top right tile")
	fs.StringVar(&t.Args.BottomLeftTile, prmBottomLeftTile, defBottomLeftTile, "top bottom left tile")
	fs.StringVar(&t.Args.BottomRightTile, prmBottomRightTile, defBottomRightTile, "bottom right tile")

	var ry, rr, rg, rb, rs string

	fs.StringVar(&ry, prmYellowRobot, defYellowRobot, "yellow robot position x,y")
	fs.StringVar(&rr, prmRedRobot, defRedRobot, "red robot position x,y")
	fs.StringVar(&rg, prmGreenRobot, defGreenRobot, "green robot position x,y")
	fs.StringVar(&rb, prmBlueRobot, defBlueRobot, "blue robot position x,y")
	fs.StringVar(&rs, prmSilverRobot, defSilverRobot, "silver robot position x,y")

	var ts, tc string
	fs.StringVar(&ts, prmTargetSymbol, defTargetSymbol, "target symbol (Pyramid|Star|Moon|Saturn|Cosmic)")
	fs.StringVar(&tc, prmTargetColor, "", "target color (yellow|red|green|blue) - leave empty for symbol Cosmic")

	fs.Parse(cmdArgs)

	if err := parseCoordinate(ry, &t.Args.YellowRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rr, &t.Args.RedRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rg, &t.Args.GreenRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rb, &t.Args.BlueRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rs, &t.Args.SilverRobot); err != nil {
		return nil, err
	}

	if err := parseSymbol(ts, &t.Args.TargetSymbol); err != nil {
		return nil, err
	}

	if t.Args.TargetSymbol != Cosmic { // no color for symbol cosmic
		if err := parseColor(tc, &t.Args.TargetColor); err != nil {
			return nil, err
		}
	}

	if err := t.checkArgs(); err != nil {
		return nil, err
	}
	return t, nil
}

// ParseFlags returns a task object build by command line flags.
func ParseFlags() (*Task, error) { return parseFlags(os.Args[0], os.Args[1:], flag.ExitOnError) }

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
