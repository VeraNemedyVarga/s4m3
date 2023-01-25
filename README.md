# S4M3

Clear the board!
Click on clusters that are bigger than 3 tiles to clear them.

## TUI

Use the arrow keys to move the cursor, press `space` to clear. (Quit: `q`)

Keys:
- arrow keys - move cursor
- `q` / `ctrl+c` - quit
- `space` - clear cluster
- `r` / `R` - restart (same seed / new random seed)

## Web API

You can configure the game to listen on a port (specified in `config.yaml`, which is created at first run if does not exist), and use a web API to play the game.

Default API / WebSocket port is `2345`.

### Get the board

Issue a simple HTTP GET:
```
GET /api/board
```

### Clear a cluster

Post the coordinates to clear (as JSON):
```
POST /api/board

{x: 1, y: 2, restart: false, new: false}
```
