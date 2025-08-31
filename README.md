# Space Invaders Game

A terminal-based implementation of the classic space invaders game written in Go, using the [tcell](https://github.com/gdamore/tcell) library for handling terminal graphics and input. A simple game where the player controls spaceship with a laser, and shoots and destroy descending waves of alien invaders before they reach the bottom of the screen.

### Game Status
Currently the game is endless, only alien ship can be destroyed. And 3 ships can be spawned at a time.

#### Objective
The game will start at wave 1, so it will be waves were each wave will increase the number of alien ships and power. Number of waves is endless, the higher it gets the higher the score. Also can loot boxes to get extra health or increase power of the beams to destroy the alien ships quicker.

![Space Invaders Game](https://raw.githubusercontent.com/omar0ali/spaceinvader-game-cli/refs/heads/main/screenshots/spaceinvader-game-cli.png)

### Checklist:
- [X] Create a spaceship placed at the bottom of the screen.
    - [X] Add controls to the spaceship, `Mouse Motion` left and right. 
    - [X] Draw a proper Triangle shape.
    - [X] Shoots laser beam by hitting the `Left Mouse Button`.
- [X] Create a single enemy spaceship. Can be produced by hitting the `Space Bar` key. This will be changed later.
    - [X] Alien ship can move towards the player.
    - [X] Can place aliens ships in a random X position
    - [X] Alien ship automatically falls with a limit of 3 ships at a time.
- [X] Stars falling down with different speed. (just to look cool)
- [X] Adding Start Screen UI
- [X] Adding Pause Menu using `p` to pause the game.
- [ ] Spaceship (Player) health, each alien ship reaching the end will depletes the health. (just a little bit)
- [X] Show control info at the bottom of the screen.
- [ ] Alien ship should shoot the player with laser beam as well. *(will implement later)
- [ ] Logs Window, to keep track whats being rendered in the current window.
- [ ] Implement a configuration file.
- [ ] Timer


### Controls
For now mouse `leftoo click` to shoot beams, and moving the mouse left and right will be moving the spaceship.

## Getting Started

Clone repository

```bash
git clone https://github.com/omar0ali/spaceinvaders-game-cli.git
```

Run the game

```bash
go run .
```
