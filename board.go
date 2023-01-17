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

type Board struct {
	Width  int
	Height int
	Tiles  [][]TileType
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
		s += "\n"
	}
	return s
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
	ntypes := len(config.TileTypes)

	for i := 0; i < config.Height; i++ {
		tiles[i] = make([]TileType, config.Width)
		for j := 0; j < config.Width; j++ {
			tiles[i][j] = config.TileTypes[r.Intn(ntypes)]
		}
	}

	return Board{
		Width:  config.Width,
		Height: config.Height,
		Tiles:  tiles,
	}
}
