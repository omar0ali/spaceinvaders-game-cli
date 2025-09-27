# Space Invaders Game

A terminal-based implementation of the classic space invaders game written in Go, using the [tcell](https://github.com/gdamore/tcell) library for handling terminal graphics and input. A simple game where the player controls spaceship with a laser, and shoots and destroy descending waves of alien invaders before they reach the bottom of the screen.

## Install

> [!IMPORTANT]
> The game looks for `config.toml`, make sure that both the game executable and config.toml are located in the same folder and that you run the game from that folder. Will fix this issue later once the game is finished. Another option would be just by cloning this [Repository](#getting-started).

```bash
go install github.com/omar0ali/spaceinvaders-game-cli@latest
```

To install go or clone this repository follow the [steps](#getting-started).

#### Objective
The game will start at wave 1, so it will be waves were each wave will increase the number of alien ships and power. Number of waves is endless, the higher it gets the higher the score. Also can loot boxes to get extra health or increase power of the beams to destroy the alien ships quicker.

![Space Invaders Game (Start Menu)](https://raw.githubusercontent.com/omar0ali/spaceinvaders-game-cli/refs/heads/main/screenshots/spaceinvaders-game-cli-startmenu.png)

![Space Invaders Game (Gameplay 1)](https://raw.githubusercontent.com/omar0ali/spaceinvaders-game-cli/refs/heads/main/screenshots/spaceinvaders-game-cli-gameplay.png)

![Space Invaders Game (Pause)](https://raw.githubusercontent.com/omar0ali/spaceinvaders-game-cli/refs/heads/main/screenshots/spaceinvaders-game-cli-pause.png)

![Space Invaders Game (Game Over)](https://raw.githubusercontent.com/omar0ali/spaceinvaders-game-cli/refs/heads/main/screenshots/spaceinvaders-game-cli-gameover.png)

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
    - [X] Adding spaceship details
- [X] Spaceship (Player) health, each alien ship reaching the end will depletes the health. (just a little bit)
- [X] Show control info at the bottom of the screen.
- [X] Loot boxes to increase health.
    - [X] Health Drop Packs will increase every minute by one and to use a health pack need to press `[H]` key to deploy and must destroy the health box to obtain.
- [X] Implement a configuration file.
- [X] Timer
- [X] Spaceship stats on the left of screen.
- [X] Alien Ship Stats on the right of the screen.
- [X] Alien ships can shoot the player's spaceship.
- [X] Spaceship now can move around the whole screen using the mouse.
- [X] Better way to level up the player.
    - [X] With upgrade choice /  stat selection. Player can choose to upgrade either (gun power, speed or capacity).
- [ ] Must find a way to run the game without `config` file. Using default values.
- [ ] Final Boss fight when reaching either time limit or a given score. After that the game will repeat but with higher difficulty.

### Controls

| Control            | Action                                |
|-------------------|----------------------------------------|
| Left Mouse Click / Space Bar | Shoot beams                 |
| Mouse Movement     | Move the spaceship                    |
| F                  | Drop health kit (increase spaceship health) |
| P                  | Pause the game                        |
| S                  | Start the game or Upgrade Gun Speed   |
| A                  | Upgrade Gun Power                     |
| D                  | Upgrade Gun Capacity                  |
| C                  | Increase heatlh

### Default Configuration File
Configuration file added for the player to freely change/update entity's attributes. The config file saved as `config.toml`.

```toml
[spaceship]
health = 15
max_level = 10
next_level_score = 300
gun_max_cap = 6
gun_cap = 3
gun_speed = 40
gun_max_speed = 80
gun_power = 2

[aliens]
limit = 1
health = 5
speed = 3
gun_speed = 35
gun_power = 1

[stars] 
limit = 15
speed = 50

[health_drop]
health = 6
speed = 3
limit = 1
max_drop = 5
start = 1

```

## Getting Started

> [!NOTE]
> Ensure go is installed [Install Golang](https://go.dev/doc/install)

Clone repository

```bash
git clone https://github.com/omar0ali/spaceinvaders-game-cli.git
```

Run the game

```bash
go run .
```
