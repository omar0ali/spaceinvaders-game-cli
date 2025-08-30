# Space Invader Game

A terminal-based implementation of the classic space invaders game written in Go, using the [tcell](https://github.com/gdamore/tcell) library for handling terminal graphics and input. A simple game where the player controls spaceship with a laser, and shoots and destroy descending waves of alien invaders before they reach the bottom of the screen.

![Space Invaders Game](https://raw.githubusercontent.com/omar0ali/spaceinvader-game-cli/refs/heads/main/screenshots/spaceinvader-game-cli.png)

### TODOS
- [X] Create a spaceship placed at the bottom of the screen.
    - [X] Add controls to the spaceship, `Mouse Motion` left and right. 
    - [X] Draw a proper Triangle shape.
    - [X] Shoots laser beam by hitting the `Left Mouse Button`.
- [X] Create a single enemy spaceship. Can be produced by hitting the `Space Bar` key. This will be changed later.
    - [X] Alien ship can move towards the player.
    - [X] Can place aliens ships in a random X position
    - [X] Alien ship automatically falls with a limit of 3 ships at a time.
- [ ] Spaceship (Player) health, each alien ship reaching the end will depletes the health. (just a little bit)
- [X] Stars falling down with different speed. (just to look cool)
- [ ] Logs Window, to keep track whats being rendered in the current window.
- [ ] Show control info at the bottom of the screen.
- [ ] Alien ship should shoot the player with laser beam as well. *(will implement later)

## Getting Started

Clone repository

```bash
git clone https://github.com/omar0ali/spaceinvader-game-cli.git
```

Run the game

```bash
go run .
```
