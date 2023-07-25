// Package task provides types for solver processes.
package task

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

// Coordinate defines a two dimensional coordinate.
type Coordinate struct {
	X int
	Y int
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
		panic(fmt.Sprintf("invalid symbol %d", s))
	}
	return symbolStrs[s]
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

// Task contains all game relevant data for a solver to calculate a game solution.
type Task struct {
	Args   *Args
	logger *slog.Logger
	start  time.Time
}

// New returns a new task instance.
func New(args *Args) *Task {
	return &Task{
		Args:   args,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})),
		start:  time.Now(),
	}
}

// NewByFlag returns a new task instance with arguments parsed from flags.
func NewByFlag() (*Task, error) {
	args, err := parseFlag(os.Args[0], os.Args[1:], flag.ExitOnError)
	if err != nil {
		return nil, err
	}
	return New(args), nil
}

// NewByURL returns a new task instance with arguments parsed from an url.
func NewByURL(u *url.URL) (*Task, error) {
	args, err := parseURL(u)
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
	To    Coordinate
	Robot Robot
}

// Moves describes a solution of a game as ordered list of robot moves.
type Moves []*Move
