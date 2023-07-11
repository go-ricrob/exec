package task

/*
func TestFlag(t *testing.T) {
	defTask &Task{
		TopLeftTile: defTopLeftTile,
		TopRightTile: TileID(defTopRightTile),



		map[game.TilePosition]game.TileID{
			game.TopLeft:     game.DefTopLeftTile,
			game.TopRight:    game.DefTopRightTile,
			game.BottomRight: game.DefBottomRightTile,
			game.BottomLeft:  game.DefBottomLeftTile,
		},
		map[game.Color]game.Coordinate{
			game.Yellow: game.DefYellowRobot,
			game.Red:    game.DefRedRobot,
			game.Green:  game.DefGreenRobot,
			game.Blue:   game.DefBlueRobot,
		},
		game.DefTarget,
	)
	if err != nil {
		t.Fatal(err)
	}
	noSilverRobotTask, err := task.New(
		map[game.TilePosition]game.TileID{
			game.TopLeft:     {SetID: 'A', TileNo: 1, Front: true},
			game.TopRight:    {SetID: 'A', TileNo: 2, Front: true},
			game.BottomRight: {SetID: 'A', TileNo: 3, Front: true},
			game.BottomLeft:  {SetID: 'A', TileNo: 4, Front: true},
		},
		map[game.Color]game.Coordinate{
			game.Yellow: {X: 0, Y: 0},
			game.Red:    {X: 0, Y: 1},
			game.Green:  {X: 0, Y: 2},
			game.Blue:   {X: 0, Y: 3},
		},
		game.Targets[game.TnYellowStar],
	)
	if err != nil {
		t.Fatal(err)
	}
	silverRobotTask, err := task.New(
		map[game.TilePosition]game.TileID{
			game.TopLeft:     {SetID: 'A', TileNo: 1, Front: true},
			game.TopRight:    {SetID: 'A', TileNo: 2, Front: true},
			game.BottomRight: {SetID: 'A', TileNo: 3, Front: true},
			game.BottomLeft:  {SetID: 'A', TileNo: 4, Front: true},
		},
		map[game.Color]game.Coordinate{
			game.Yellow: {X: 0, Y: 0},
			game.Red:    {X: 0, Y: 1},
			game.Green:  {X: 0, Y: 2},
			game.Blue:   {X: 0, Y: 3},
			game.Silver: {X: 0, Y: 4},
		},
		game.Targets[game.TnYellowStar],
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args []string
		err  error // expected error
		task *task.Task
	}{
		{"default", []string{}, nil, defTask},

		{"wrong yellow robot position", []string{"-ry", "-1,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}, task.ErrInvalidRobotPosition, nil},
		{"wrong yellow robot position", []string{"-ry", "0,-1", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}, task.ErrInvalidRobotPosition, nil},
		{"wrong yellow robot position", []string{"-ry", "16,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}, task.ErrInvalidRobotPosition, nil},
		{"wrong yellow robot position", []string{"-ry", "0,16", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3"}, task.ErrInvalidRobotPosition, nil},

		{"invalid tile", []string{"-ttl", "A5F"}, task.ErrInvalidTile, nil},
		{"duplicate tile", []string{"-ttr", "A1F"}, task.ErrDuplicateTile, nil},

		{"duplicate robot position", []string{"-ry", "0,1", "-rr", "0,2", "-rg", "0,3", "-rb", "0,2"}, task.ErrDuplicateRobotPosition, nil},

		{"check board without silver robot", []string{"-ry", "0,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3", "-t", "yellowStar"}, nil, noSilverRobotTask},
		{"check board with silver robot", []string{"-ry", "0,0", "-rr", "0,1", "-rg", "0,2", "-rb", "0,3", "-rs", "0,4", "-t", "yellowStar"}, nil, silverRobotTask},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task, err := parseFlags(true, true, test.name, test.args)

			if test.err != nil {
				if !errors.Is(err, test.err) {
					t.Fatalf("error %v - expected %v", err, test.err)
				}
			} else { // case test.expectedErrors == nil && err == nil:
				if !reflect.DeepEqual(task.Tiles(), test.task.Tiles()) {
					t.Fatalf("tiles %v - expected %v", task.Tiles(), test.task.Tiles())
				}
				if !reflect.DeepEqual(task.Robots(), test.task.Robots()) {
					t.Fatalf("robots %v - expected %v", task.Robots(), test.task.Robots())
				}
				if task.Target() != test.task.Target() {
					t.Fatalf("target %v - expected %v", task.Target(), test.task.Target())
				}
			}
		})
	}
}
*/
