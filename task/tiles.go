package task

import (
	"flag"
	"fmt"
	"net/url"
)

// Tile constants.
const (
	fnTopLeftTile     = "ttl"
	fnTopRightTile    = "ttr"
	fnBottomRightTile = "tbr"
	fnBottomLeftTile  = "tbl"
)

const (
	argTopLeftTile     = "-" + fnTopLeftTile
	argTopRightTile    = "-" + fnTopRightTile
	argBottomRightTile = "-" + fnBottomRightTile
	argBottomLeftTile  = "-" + fnBottomLeftTile
)

// Default parameters.
var (
	defTopLeftTile     = "A1F"
	defTopRightTile    = "A2F"
	defBottomRightTile = "A3F"
	defBottomLeftTile  = "A4F"
)

// Tiles represents the tiles of a board.
type Tiles struct {
	TopLeft, TopRight, BottomLeft, BottomRight string
}

func (t *Tiles) check() error {
	// check tiles
	tiles := []string{t.TopLeft, t.TopRight, t.BottomLeft, t.BottomRight}
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
	return nil
}

func (t *Tiles) addFlag(fs *flag.FlagSet) {
	fs.StringVar(&t.TopLeft, fnTopLeftTile, defTopLeftTile, "top left tile")
	fs.StringVar(&t.TopRight, fnTopRightTile, defTopRightTile, "top right tile")
	fs.StringVar(&t.BottomLeft, fnBottomLeftTile, defBottomLeftTile, "top bottom left tile")
	fs.StringVar(&t.BottomRight, fnBottomRightTile, defBottomRightTile, "bottom right tile")
}

// ParseURL extract the tiles out of a URL query.
func (t *Tiles) ParseURL(u *url.URL) error {
	query := u.Query()
	t.TopLeft = queryString(query, fnTopLeftTile, defTopLeftTile)
	t.TopRight = queryString(query, fnTopRightTile, defTopRightTile)
	t.BottomLeft = queryString(query, fnBottomLeftTile, defBottomLeftTile)
	t.BottomRight = queryString(query, fnBottomRightTile, defBottomRightTile)

	if err := t.check(); err != nil {
		return err
	}
	return nil
}
