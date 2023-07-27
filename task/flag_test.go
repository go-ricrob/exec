package task

import (
	"flag"
	"reflect"
	"testing"
)

func TestFlag(t *testing.T) {

	t.Run("negative test", func(t *testing.T) {
		tests := []struct {
			name    string
			cmdArgs []string
		}{

			{"invalid yellow robot position", []string{"-ry", "-1,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"invalid yellow robot position", []string{"-ry", "0,-1", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"outside board yellow robot position", []string{"-ry", "16,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"outside board yellow robot position", []string{"-ry", "0,16", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"center field yellow robot position", []string{"-ry", "7,7", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"center field yellow robot position", []string{"-ry", "8,8", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"center field yellow robot position", []string{"-ry", "7,7", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"center field yellow robot position", []string{"-ry", "8,8", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},
			{"symbol yellow robot position", []string{"-ry", "7,10", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}},

			{"invalid tile", []string{"-ttl", "A5F"}},
			{"duplicate tile", []string{"-ttr", "A1F"}},

			{"duplicate robot position", []string{"-ry", "0,1", "-rr", "0,2", "-rg", "0,3", "-rb", "0,2"}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if _, err := parseFlag(test.name, test.cmdArgs, flag.ContinueOnError); err != nil {
					t.Log(err)
				} else {
					t.Fatal("error expected")
				}
			})
		}
	})

	t.Run("positive test", func(t *testing.T) {
		tests := []struct {
			name    string
			cmdArgs []string
			args    *Args
		}{
			{
				"default",
				[]string{},
				&Args{
					Tiles: Tiles{
						TopLeft:     defTopLeftTile,
						TopRight:    defTopRightTile,
						BottomLeft:  defBottomLeftTile,
						BottomRight: defBottomRightTile,
					},
					Robots: Robots{
						Yellow: Coordinate{X: 0, Y: 0},
						Red:    Coordinate{X: 1, Y: 0},
						Green:  Coordinate{X: 2, Y: 0},
						Blue:   Coordinate{X: 3, Y: 0},
						Silver: Coordinate{X: -1, Y: -1},
					},
					TargetSymbol:       Cosmic,
					CheckRobotOnSymbol: true,
				},
			},
			{
				"check board without silver robot",
				[]string{"-ry", "0,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3", "-ts", "yellowStar"},
				&Args{
					Tiles: Tiles{
						TopLeft:     defTopLeftTile,
						TopRight:    defTopRightTile,
						BottomLeft:  defBottomLeftTile,
						BottomRight: defBottomRightTile,
					},
					Robots: Robots{
						Yellow: Coordinate{X: 0, Y: 0},
						Red:    Coordinate{X: 0, Y: 1},
						Green:  Coordinate{X: 0, Y: 2},
						Blue:   Coordinate{X: 0, Y: 3},
						Silver: Coordinate{X: -1, Y: -1},
					},
					TargetSymbol:       YellowStar,
					CheckRobotOnSymbol: true,
				},
			},
			{
				"check board with silver robot",
				[]string{"-ry", "0,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3", "-rs", "0,4", "-ts", "yellowStar"},
				&Args{
					Tiles: Tiles{
						TopLeft:     defTopLeftTile,
						TopRight:    defTopRightTile,
						BottomLeft:  defBottomLeftTile,
						BottomRight: defBottomRightTile,
					},
					Robots: Robots{
						Yellow: Coordinate{X: 0, Y: 0},
						Red:    Coordinate{X: 0, Y: 1},
						Green:  Coordinate{X: 0, Y: 2},
						Blue:   Coordinate{X: 0, Y: 3},
						Silver: Coordinate{X: 0, Y: 4},
					},
					TargetSymbol:       YellowStar,
					CheckRobotOnSymbol: true,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				args, err := parseFlag(test.name, test.cmdArgs, flag.ContinueOnError)
				if err != nil {
					t.Fatal(err)
				}
				//if *task.Args != *test.args {
				if !reflect.DeepEqual(args, test.args) {
					t.Fatalf("task %v - expected %v", args, test.args)
				}
			})
		}
	})
}
