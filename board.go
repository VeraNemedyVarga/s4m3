package main

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type TileType struct {
	Sign  string `yaml:"sign"`
	Color string `yaml:"color"`
}

var EMPTY_TILE = TileType{
	Sign:  " ",
	Color: "#000000",
}

var TMP_TILE = TileType{
	Sign:  "_",
	Color: "#000000",
}

type Board struct {
	Width      int          `json:"width"`
	Height     int          `json:"height"`
	Tiles      [][]TileType `json:"tiles"`
	Seed       int64        `json:"seed"`
	ExtraTiles int          `json:"extra_tiles"`
	config     *Config
	rand       *rand.Rand
}

func (b Board) String() string {
	var s string
	for _, row := range b.Tiles {
		for _, tile := range row {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(tile.Color))
			s += style.Render(tile.Sign)
		}
		s += "\n"
	}
	return s
}

func (b Board) WithCursor(cx, cy int) string {
	var s string
	for y, row := range b.Tiles {
		for x, tile := range row {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(tile.Color))
			if x == cx && y == cy {
				style = style.Background(lipgloss.Color("#668800"))
			}
			s += style.Render(tile.Sign)
		}
		if y != len(b.Tiles)-1 {
			s += "\n"
		}
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#aadd88")).
		Padding(1, 2)

	return style.Render(s)
}

func (b Board) floodFill(t, ft TileType, cx, cy int) int {
	if cx < 0 || cx >= b.Width || cy < 0 || cy >= b.Height || t == EMPTY_TILE {
		return 0
	}

	current := b.Tiles[cy][cx]
	if current == t {
		b.Tiles[cy][cx] = ft
		return 1 + b.floodFill(t, ft, cx+1, cy) + b.floodFill(t, ft, cx-1, cy) + b.floodFill(t, ft, cx, cy+1) + b.floodFill(t, ft, cx, cy-1)
	}

	return 0
}

func (b Board) shake() {
	for col := 0; col < b.Width; col++ {
		for row := b.Height - 1; row > 0; row-- {
			if b.Tiles[row][col] == EMPTY_TILE {
				for row2 := row - 1; row2 >= 0; row2-- {
					if b.Tiles[row2][col] != EMPTY_TILE {
						b.Tiles[row][col] = b.Tiles[row2][col]
						b.Tiles[row2][col] = EMPTY_TILE
						break
					}
				}
			}
		}
	}
}

func (b *Board) fillBack() {
	for row := b.Height - 1; row >= 0 && b.ExtraTiles > 0; row-- {
		for col := 0; col < b.Width && b.ExtraTiles > 0; col++ {
			if b.Tiles[row][col] == EMPTY_TILE {
				b.Tiles[row][col] = randomTile(*b.config, *b.rand)
				b.ExtraTiles--
			}
		}
	}
}

func (b *Board) Hit(cx, cy int) int {
	if cx < 0 || cx >= b.Width || cy < 0 || cy >= b.Height {
		return 0
	}

	reftile := b.Tiles[cy][cx]
	csize := b.floodFill(reftile, TMP_TILE, cx, cy)

	if csize > 2 {
		csize := b.floodFill(TMP_TILE, EMPTY_TILE, cx, cy)
		b.shake()
		b.fillBack()
		return 1 + (csize-3)*2
	}

	// cluster too small, revert fill
	b.floodFill(b.Tiles[cy][cx], reftile, cx, cy)
	return 0
}

func (b Board) HasMove() bool {
	for cy, row := range b.Tiles {
		for cx, tile := range row {
			if tile != EMPTY_TILE {
				reftile := b.Tiles[cy][cx]
				csize := b.floodFill(reftile, TMP_TILE, cx, cy)
				b.floodFill(TMP_TILE, reftile, cx, cy)
				if csize > 2 {
					return true
				}
			}
		}
	}
	return false
}

func randomTile(c Config, r rand.Rand) TileType {
	ntypes := len(c.TileTypes)
	return c.TileTypes[r.Intn(ntypes)]
}

func generateBoard(config Config) Board {
	var tiles = make([][]TileType, config.Height)

	var seed int64
	if config.Seed == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = config.Seed
	}

	s := rand.NewSource(seed)
	r := rand.New(s)

	for i := 0; i < config.Height; i++ {
		tiles[i] = make([]TileType, config.Width)
		for j := 0; j < config.Width; j++ {
			tiles[i][j] = randomTile(config, *r)
		}
	}

	return Board{
		Width:      config.Width,
		Height:     config.Height,
		Tiles:      tiles,
		Seed:       seed,
		ExtraTiles: config.ExtraTiles,
		rand:       r,
		config:     &config,
	}
}
