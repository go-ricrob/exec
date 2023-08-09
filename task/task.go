// Package task provides types for solver processes.
package task

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Coordinate defines a two dimensional coordinate.
type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) String() string { return fmt.Sprintf("%d,%d", c.X, c.Y) }

// MarshalText implments the encoding/MarshalText interface.
func (c Coordinate) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
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

// UnmarshalText implments the encoding/UnmarshalText interface.
func (c *Coordinate) UnmarshalText(text []byte) (err error) {
	*c, err = parseCoordinate(string(text))
	return err
}

// Symbol is the type of a symbol.
type Symbol byte

// Symbol constants.
const (
	YellowPyramid Symbol = iota
	YellowStar
	YellowMoon
	YellowSaturn

	RedPyramid
	RedStar
	RedMoon
	RedSaturn

	GreenPyramid
	GreenStar
	GreenMoon
	GreenSaturn

	BluePyramid
	BlueStar
	BlueMoon
	BlueSaturn

	Cosmic
)

const (
	strYellowPyramid = "yellowPyramid"
	strYellowStar    = "yellowStar"
	strYellowMoon    = "yellowMoon"
	strYellowSaturn  = "yellowSaturn"

	strRedPyramid = "redPyramid"
	strRedStar    = "redStar"
	strRedMoon    = "redMoon"
	strRedSaturn  = "redSaturn"

	strGreenPyramid = "greenPyramid"
	strGreenStar    = "greenStar"
	strGreenMoon    = "greenMoon"
	strGreenSaturn  = "greenSaturn"

	strBluePyramid = "bluePyramid"
	strBlueStar    = "blueStar"
	strBlueMoon    = "blueMoon"
	strBlueSaturn  = "blueSaturn"

	strCosmic = "cosmic"
)

var symbolStrs = []string{
	strYellowPyramid,
	strYellowStar,
	strYellowMoon,
	strYellowSaturn,
	strRedPyramid,
	strRedStar,
	strRedMoon,
	strRedSaturn,
	strGreenPyramid,
	strGreenStar,
	strGreenMoon,
	strGreenSaturn,
	strBluePyramid,
	strBlueStar,
	strBlueMoon,
	strBlueSaturn,
	strCosmic,
}

var symbolMap = map[string]Symbol{
	strYellowPyramid: YellowPyramid,
	strYellowStar:    YellowStar,
	strYellowMoon:    YellowMoon,
	strYellowSaturn:  YellowSaturn,

	strRedPyramid: RedPyramid,
	strRedStar:    RedStar,
	strRedMoon:    RedMoon,
	strRedSaturn:  RedSaturn,

	strGreenPyramid: GreenPyramid,
	strGreenStar:    GreenStar,
	strGreenMoon:    GreenMoon,
	strGreenSaturn:  GreenSaturn,

	strBluePyramid: BluePyramid,
	strBlueStar:    BlueStar,
	strBlueMoon:    BlueMoon,
	strBlueSaturn:  BlueSaturn,

	strCosmic: Cosmic,
}

func (s Symbol) String() string {
	if int(s) >= len(symbolStrs) {
		panic(fmt.Errorf("invalid symbol %d", s))
	}
	return symbolStrs[s]
}

func parseSymbol(s string) (Symbol, error) {
	symbol, ok := symbolMap[s]
	if !ok {
		return symbol, fmt.Errorf("invalid symbol %s", s)
	}
	return symbol, nil
}

// MarshalText implments the encoding/MarshalText interface.
func (s Symbol) MarshalText() (text []byte, err error) {
	if int(s) >= len(symbolStrs) {
		return nil, fmt.Errorf("invalid symbol %d", s)
	}
	return []byte(symbolStrs[s]), nil
}

// UnmarshalText implments the encoding/UnmarshalText interface.
func (s *Symbol) UnmarshalText(text []byte) (err error) {
	*s, err = parseSymbol(string(text))
	return err
}

// Robot is the type of a robot.
type Robot byte

// Robot constants
const (
	YellowRobot Robot = iota
	RedRobot
	GreenRobot
	BlueRobot
	SilverRobot
)

var robotStrs = []string{"yellowRobot", "redRobot", "greenRobot", "blueRobot", "silverRobot"}

func (r Robot) String() string {
	if int(r) >= len(robotStrs) {
		panic(fmt.Sprintf("invalid robot %d", r))
	}
	return robotStrs[r]
}

// Task contains all game relevant data for a solver to calculate a game solution.
type Task struct {
	Args   *Args
	logger *slog.Logger
	start  time.Time
}

// New returns a new task instance.
func New(args *Args) *Task {
	attrs := []slog.Attr{slog.String("solver", os.Args[0])}
	return &Task{
		Args:   args,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}).WithAttrs(attrs)),
		start:  time.Now(),
	}
}

// NewFromFlag returns a new task instance with arguments parsed from flag.
func NewFromFlag() (*Task, error) {
	args, err := parseFlag(os.Args[0], os.Args[1:], flag.ExitOnError)
	if err != nil {
		return nil, err
	}
	return New(args), nil
}

// Level logs the level and potential additional information provided by the solver.
func (t *Task) Level(level int, args ...any) {
	t.logger.Info("progress", append([]any{"level", level}, args...)...)
}

// Exit signals an error to the caller and exits the task process.
func (t *Task) Exit(err error) {
	t.logger.Error("exit", "duration", time.Since(t.start)/time.Millisecond, "err", err.Error())
	os.Exit(1)
}

// Result signals a task result to the caller.
func (t *Task) Result(moves Moves, args ...any) {
	t.logger.Info("result", append([]any{"duration", time.Since(t.start) / time.Millisecond, "moves", moves}, args...)...)
}

// Move defines a single move of a robot.
type Move struct {
	To    Coordinate `json:"to"`
	Robot Robot      `json:"robot"`
}

// Moves describes a solution of a game as ordered list of robot moves.
type Moves []*Move
