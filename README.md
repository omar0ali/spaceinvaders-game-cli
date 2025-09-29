# Space Invaders Game

A terminal-based implementation of the classic space invaders game written in Go, using the [tcell](https://github.com/gdamore/tcell) library for handling terminal graphics and input. A simple game where the player controls spaceship with a laser, and shoots and destroy descending waves of alien invaders before they reach the bottom of the screen.

## Install

> [!NOTE]
> The game looks for `config.toml`, make sure that both the game executable and `config.toml` are located in the same folder and that you run the game from that folder. 
>
> Note: Game can run without `config` file. Will use default configuration.
> Another option would be just by cloning this [Repository](#getting-started).
>
> There is an executable can be downloaded from [releases](https://github.com/omar0ali/spaceinvaders-game-cli/releases).

```bash
go install github.com/omar0ali/spaceinvaders-game-cli@latest
```

Run the game after installation

```bash
spaceinvaders-game-cli
```

To install go or clone this repository follow the [steps](#getting-started).

#### Objective
The game is an endless space shooter where players face increasingly difficult waves of alien ships that scale with their level. Each time the player levels up, they can choose an upgrade to improve their spaceship, such as boosting firepower to handle tougher aliens with stronger armor. The objective is to survive as long as possible, destroy alien ships, and push for a higher score while managing health through occasional drop-down health packs that restore the spaceship health.

#### Youtube - Gameplay Demo - v1.1.0
[![Watch the video](https://img.youtube.com/vi/2flPiJvl4vU/0.jpg)](https://www.youtube.com/watch?v=2flPiJvl4vU)

#### Youtube - Gameplay Demo - v1.0.0
[![Watch the video](https://img.youtube.com/vi/DSeU1Lnglsg/0.jpg)](https://www.youtube.com/watch?v=DSeU1Lnglsg)

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
- [X] The game runs without a `config` file. If the file dose not exist, the default configuration will be used.
- [X] Adding variety designs of alien ships attacking. Using json file to load all alien ships designs.
- [X] Will add variety of spaceship selection before starting the game. Currently added just five.
    - Before the game starts, the player will be able to select a ship.

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
max_level = 50
next_level_score = 50
gun_max_cap = 30
gun_max_speed = 85

[aliens]
limit = 1
health = 2
speed = 4
gun_speed = 35
gun_power = 1

[stars] 
limit = 15
speed = 50

[health_drop]
health = 2
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
