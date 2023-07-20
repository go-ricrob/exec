// Package task provides types for solver processes.
package task

import (
	"flag"
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
	args, err := parseArgs(os.Args[0], os.Args[1:], flag.ExitOnError)
	if err != nil {
		return nil, err
	}
	return New(args), nil
}

// Level logs the level and potential additional information provided by the solver.
func (t *Task) Level(level int, args ...any) {
	t.logger.Info("progress", append([]any{"level", level}, args...))
}

// Exit signals an error to the caller and exits the task process.
func (t *Task) Exit(err error) {
	t.logger.Error("exit", "duration", time.Since(t.start)/time.Millisecond, "err", err.Error())
	os.Exit(1)
}

// Result signals a task result to the caller.
func (t *Task) Result(moves Moves, args ...any) {
	t.logger.Info("result", "duration", append([]any{"duration", time.Since(t.start) / time.Millisecond, "moves", moves}, args...))
}

// Move defines a single move of a robot.
type Move struct {
	To    Coordinate
	Color Color
}

// Moves describes a solution of a game as ordered list of robot moves.
type Moves []*Move
