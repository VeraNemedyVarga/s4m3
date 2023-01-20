package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TileTypes []TileType `yaml:"tile_types"`
	Seed      int64      `yaml:"seed"`
	Width     int        `yaml:"width"`
	Height    int        `yaml:"height"`
	Addr      string     `yaml:"addr"`
}

var defaultConfig = Config{
	TileTypes: []TileType{
		{
			Sign:  "*",
			Color: "#ffff00",
		},
		{
			Sign:  "X",
			Color: "#88ff00",
		},
		{
			Sign:  "O",
			Color: "#0088ff",
		},
	},
	Seed:   0,
	Width:  20,
	Height: 20,
}

type model struct {
	config   Config
	board    Board
	gameOver bool
	cx       int
	cy       int
	points   int
	sub      chan webHitMsg
}

type webHitMsg interface {
	getResp() chan model
}

func waitForWebHit(sub chan webHitMsg) tea.Cmd {
	return func() tea.Msg {
		return webHitMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return waitForWebHit(m.sub)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			m.points += m.board.Hit(m.cx, m.cy)
		case "up":
			m.cy--
		case "down":
			m.cy++
		case "left":
			m.cx--
		case "right":
			m.cx++
		}
	case webHitMsg:
		switch msg.(type) {
		case WebHit:
			x, y := msg.(WebHit).getCoords()
			m.points += m.board.Hit(x, y)
			msg.getResp() <- m
			return m, waitForWebHit(m.sub)
		case WebGet:
			msg.getResp() <- m
			return m, waitForWebHit(m.sub)
		}
	}

	if m.cx < 0 {
		m.cx = 0
	}
	if m.cx >= m.config.Width {
		m.cx = m.config.Width - 1
	}
	if m.cy < 0 {
		m.cy = 0
	}
	if m.cy >= m.config.Height {
		m.cy = m.config.Height - 1
	}

	return m, nil
}

func (m model) View() string {
	return m.board.WithCursor(m.cx, m.cy) + fmt.Sprintf("Points: %d", m.points)
}

func initialModel(cfg Config) model {
	return model{
		config: cfg,
		board:  generateBoard(cfg),
		sub:    make(chan webHitMsg),
	}
}

func saveConfig(path string, cfg Config) error {
	ycfg, err := yaml.Marshal(cfg)

	if err != nil {
		log.Println(err)
		return err
	}

	err2 := ioutil.WriteFile(path, ycfg, 0644)

	if err2 != nil {
		log.Println(err2)
		return err2
	}

	log.Println("Configuration saved to", path)
	return nil
}

func readConfig(path string) (Config, error) {
	var cfg Config

	if path == "" {
		path = "config.yaml"
	}

	ycfg, err := ioutil.ReadFile(path)

	if err != nil {
		// error reading config file, use default config and try to save it
		errS := saveConfig(path, defaultConfig)
		return defaultConfig, errS
	}

	err2 := yaml.Unmarshal(ycfg, &cfg)

	if err2 != nil {
		return defaultConfig, err2
	}

	return cfg, nil
}

func main() {
	configPath := flag.String("c", "config.yaml", "path to config file")

	cfg, err := readConfig(*configPath)

	if err != nil {
		log.Println(err)
	}

	m := initialModel(cfg)

	if cfg.Addr != "" {
		go initApi(cfg, m.sub)
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("could not start game: %v", err)
		os.Exit(1)
	}

}
