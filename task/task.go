// Package task provides types for solver processes.
package task

import (
	"os"
	"time"

	"golang.org/x/exp/slog"
)

// TileID defines the id of a type (tile number & side).
type TileID struct {
	SetID  byte
	TileNo byte
	Front  bool
}

// Coordinate defines a two dimensional coordinate.
type Coordinate struct {
	X int
	Y int
}

// Symbol is the type of a symbol.
type Symbol byte

// Symbol constants.
const (
	Pyramid Symbol = iota
	Star
	Moon
	Saturn
	Cosmic
)

// Color is the type of a color.
type Color byte

// Color constants
const (
	Yellow = iota
	Red
	Green
	Blue
	Silver
	CosmicColor
)

// Task contains all game relevant data for a solver to calculate a game solution.
type Task struct {
	TopLeftTile     TileID
	TopRightTile    TileID
	BottomLeftTile  TileID
	BottomRightTile TileID

	YellowRobot Coordinate
	RedRobot    Coordinate
	GreenRobot  Coordinate
	BlueRobot   Coordinate
	SilverRobot Coordinate

	TargetSymbol Symbol
	TargetColor  Color

	logger   *slog.Logger
	start    time.Time
	progress int
}

// newTask returns a new task instance.
func newTask() *Task {
	return &Task{logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})), start: time.Now()}
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

// Move defines a single move of a robot.
type Move struct {
	To    Coordinate
	Color Color
}

// Moves describes a solution of a game as ordered list of robot moves.
type Moves []*Move
