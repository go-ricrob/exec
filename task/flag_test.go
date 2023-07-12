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
				if _, err := parseArgs(test.name, test.cmdArgs, flag.ContinueOnError); err != nil {
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
					TopLeftTile:     defTopLeftTile,
					TopRightTile:    defTopRightTile,
					BottomLeftTile:  defBottomLeftTile,
					BottomRightTile: defBottomRightTile,
					YellowRobot:     Coordinate{X: 0, Y: 0},
					RedRobot:        Coordinate{X: 1, Y: 0},
					GreenRobot:      Coordinate{X: 2, Y: 0},
					BlueRobot:       Coordinate{X: 3, Y: 0},
					SilverRobot:     Coordinate{X: -1, Y: -1},
					TargetSymbol:    Cosmic,
				},
			},
			{
				"check board without silver robot",
				[]string{"-ry", "0,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3", "-ts", "Star", "-tc", "yellow"},
				&Args{
					TopLeftTile:     defTopLeftTile,
					TopRightTile:    defTopRightTile,
					BottomLeftTile:  defBottomLeftTile,
					BottomRightTile: defBottomRightTile,
					YellowRobot:     Coordinate{X: 0, Y: 0},
					RedRobot:        Coordinate{X: 0, Y: 1},
					GreenRobot:      Coordinate{X: 0, Y: 2},
					BlueRobot:       Coordinate{X: 0, Y: 3},
					SilverRobot:     Coordinate{X: -1, Y: -1},
					TargetSymbol:    Star,
					TargetColor:     Yellow,
				},
			},
			{
				"check board with silver robot",
				[]string{"-ry", "0,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3", "-rs", "0,4", "-ts", "Star", "-tc", "yellow"},
				&Args{
					TopLeftTile:     defTopLeftTile,
					TopRightTile:    defTopRightTile,
					BottomLeftTile:  defBottomLeftTile,
					BottomRightTile: defBottomRightTile,
					YellowRobot:     Coordinate{X: 0, Y: 0},
					RedRobot:        Coordinate{X: 0, Y: 1},
					GreenRobot:      Coordinate{X: 0, Y: 2},
					BlueRobot:       Coordinate{X: 0, Y: 3},
					SilverRobot:     Coordinate{X: 0, Y: 4},
					TargetSymbol:    Star,
					TargetColor:     Yellow,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				args, err := parseArgs(test.name, test.cmdArgs, flag.ContinueOnError)
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
