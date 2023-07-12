// Package task provides types for solver processes.
package task

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ricrob/game/board"
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
	_ Symbol = iota // no symbol
	Pyramid
	Star
	Moon
	Saturn
	Cosmic
)

// Color is the type of a color.
type Color byte

// Color constants
const (
	_ Color = iota // no color
	Yellow
	Red
	Green
	Blue
	Silver
	CosmicColor
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
	TargetColor  Color
}

// Task contains all game relevant data for a solver to calculate a game solution.
type Task struct {
	Args     *Args
	logger   *slog.Logger
	start    time.Time
	progress int
}

// newTask returns a new task instance.
func newTask() *Task {
	return &Task{
		Args:   new(Args),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})),
		start:  time.Now(),
	}
}

// IncrProgress increments progress and signals the estimated task completion in percent.
func (t *Task) IncrProgress(percent int) { t.SetProgress(t.progress + percent) }

// SetProgress sets progress and signals the estimated task completion in percent.
func (t *Task) SetProgress(percent int) {
	t.progress = percent
	if t.progress > 100 {
		t.progress = 100
	}
	t.logger.Info("progress", "percent", percent)
}

// Exit signals an error to the caller and exits the task process.
func (t *Task) Exit(err error) {
	t.logger.Error("exit", "duration", time.Since(t.start)/time.Millisecond, "err", err.Error())
	os.Exit(1)
}

// Result signals a task result to the caller.
func (t *Task) Result(moves Moves, args ...any) {
	logArgs := []any{
		"duration",
		time.Since(t.start) / time.Millisecond,
		"moves",
		moves,
	}
	t.logger.Info("result", "duration", append(logArgs, args...))
}

func (t *Task) checkArgs() error {
	// check tiles
	tiles := []string{t.Args.TopLeftTile, t.Args.TopRightTile, t.Args.BottomLeftTile, t.Args.BottomRightTile}
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

	b := board.New(map[board.TilePosition]string{
		board.TopLeft:     t.Args.TopLeftTile,
		board.TopRight:    t.Args.TopRightTile,
		board.BottomLeft:  t.Args.BottomLeftTile,
		board.BottomRight: t.Args.BottomRightTile,
	})

	// check robots
	robots := []Coordinate{t.Args.YellowRobot, t.Args.RedRobot, t.Args.GreenRobot, t.Args.BlueRobot}
	if t.Args.SilverRobot.X != -1 && t.Args.SilverRobot.Y != -1 {
		robots = append(robots, t.Args.SilverRobot)
	}
	robotMap := map[Coordinate]bool{}
	for _, robot := range robots {
		if !b.IsValidCoordinate(robot.X, robot.Y) {
			return fmt.Errorf("invalid robot coordinates %v - center field", robot)
		}
		field := b.Field(robot.X, robot.Y)
		if field.Symbol() != board.NoSymbol {
			return fmt.Errorf("robot %v sits on symbol %s color %s", robot, field.Symbol(), field.Color())
		}
		if _, ok := robotMap[robot]; ok {
			return fmt.Errorf("duplicate robot position %v", robot)
		}
		robotMap[robot] = true
	}
	return nil
}

// Move defines a single move of a robot.
type Move struct {
	To    Coordinate
	Color Color
}

// Moves describes a solution of a game as ordered list of robot moves.
type Moves []*Move
