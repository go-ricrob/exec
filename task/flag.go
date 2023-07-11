package task

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ricrob/game/types"
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

func convertSideToBool(b byte) (bool, error) {
	switch b {
	case 'F', 'f':
		return true, nil
	case 'B', 'b':
		return false, nil
	default:
		return false, fmt.Errorf("invalid tile side: %c", b)
	}
}

func parseTileID(s string, tileID *TileID) error {
	if len(s) != 3 {
		return fmt.Errorf("invalid tile length: %s", s)
	}

	i64, err := strconv.ParseInt(s[1:2], 10, 8)
	if err != nil {
		return fmt.Errorf("invalid tile number %s - %w", s, err)
	}
	front, err := convertSideToBool(s[2])
	if err != nil {
		return fmt.Errorf("invalid tile side %s - %w", s, err)
	}

	setID := s[0]
	if setID != 'A' && setID != 'B' {
		return fmt.Errorf("invalid tile set id %s", s)
	}
	tileNo := byte(i64)
	if tileNo < 1 || tileNo > 4 {
		return fmt.Errorf("invalid tile number %s", s)
	}

	tileID.SetID, tileID.TileNo, tileID.Front = setID, tileNo, front
	return nil
}

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
	"Yellow": Yellow,
	"Red":    Red,
	"Green":  Green,
	"Blue":   Blue,
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

// ParseFlags returns a task object build by command line flags.
func ParseFlags() (*Task, error) {

	var ttl, ttr, tbr, tbl string
	flag.StringVar(&ttl, prmTopLeftTile, defTopLeftTile, "top left tile")
	flag.StringVar(&ttr, prmTopRightTile, defTopRightTile, "top right tile")
	flag.StringVar(&tbr, prmBottomRightTile, defBottomRightTile, "bottom right tile")
	flag.StringVar(&tbl, prmBottomLeftTile, defBottomLeftTile, "top bottom left tile")

	var ry, rr, rg, rb, rs string

	flag.StringVar(&ry, prmYellowRobot, defYellowRobot, "yellow robot position x,y")
	flag.StringVar(&rr, prmRedRobot, defRedRobot, "red robot position x,y")
	flag.StringVar(&rg, prmGreenRobot, defGreenRobot, "green robot position x,y")
	flag.StringVar(&rb, prmBlueRobot, defBlueRobot, "blue robot position x,y")
	flag.StringVar(&rs, prmSilverRobot, defSilverRobot, "silver robot position x,y")

	var ts, tc string
	flag.StringVar(&ts, prmTargetSymbol, defTargetSymbol, "target symbol (Pyramid|Star|Moon|Saturn|Cosmic)")
	flag.StringVar(&tc, prmTargetColor, "", "target color (yellow|red|green|blue) - leave empty for symbol Cosmic")

	flag.Parse()

	task := newTask()

	if err := parseTileID(ttl, &task.TopLeftTile); err != nil {
		return nil, err
	}
	if err := parseTileID(ttr, &task.TopRightTile); err != nil {
		return nil, err
	}
	if err := parseTileID(tbr, &task.BottomRightTile); err != nil {
		return nil, err
	}
	if err := parseTileID(tbl, &task.BottomLeftTile); err != nil {
		return nil, err
	}

	// check duplicate tiles
	tiles := []string{ttl, ttr, tbr, tbl}
	tileMap := map[string]bool{}
	for _, tile := range tiles {
		if _, ok := tileMap[tile]; ok {
			return nil, fmt.Errorf("duplicate tile %s", tile)
		}
		tileMap[tile] = true
	}

	if err := parseCoordinate(ry, &task.YellowRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rr, &task.RedRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rg, &task.GreenRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rb, &task.BlueRobot); err != nil {
		return nil, err
	}
	if err := parseCoordinate(rs, &task.SilverRobot); err != nil {
		return nil, err
	}

	if err := parseSymbol(ts, &task.TargetSymbol); err != nil {
		return nil, err
	}

	if err := parseColor(ts, &task.TargetColor); err != nil {
		return nil, err
	}

	// check robots
	robots := []Coordinate{task.YellowRobot, task.RedRobot, task.GreenRobot, task.BlueRobot}
	if task.SilverRobot.X != -1 && task.SilverRobot.Y != -1 {
		robots = append(robots, task.SilverRobot)
	}
	robotMap := map[Coordinate]bool{}
	for _, robot := range robots {
		if types.IsCenterField(robot.X, robot.Y) {
			return nil, fmt.Errorf("invalid robot coordinates %v - center field", robot)
		}
		if !types.IsInRange(robot.X, robot.Y) {
			return nil, fmt.Errorf("invalid robot coordinates %v - field not in range", robot)
		}
		if _, ok := robotMap[robot]; ok {
			return nil, fmt.Errorf("duplicate robot coordinates %v", robot)
		}
		robotMap[robot] = true
	}

	// TODO: check that robot does not sit on symbol

	return task, nil
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
