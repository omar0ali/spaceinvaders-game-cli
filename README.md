# Space Invaders Game
A terminal-based implementation of the classic space invaders game written in Go, using the [tcell](https://github.com/gdamore/tcell) library for handling terminal graphics and input.


## Install

> [!IMPORTANT]  
> You must install **Go version 1.24 or later**.


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

#### Youtube - Gameplay Demo - v1.4.5
[![Watch the video](https://img.youtube.com/vi/bGpVdy3bS5w/0.jpg)](https://www.youtube.com/watch?v=bGpVdy3bS5w)

#### Youtube - Gameplay Demo - v1.3.0
[![Watch the video](https://img.youtube.com/vi/6RS4EPAq_ag/0.jpg)](https://www.youtube.com/watch?v=6RS4EPAq_ag)

#### Youtube - Gameplay Demo - v1.2.0
[![Watch the video](https://img.youtube.com/vi/23ziir_IJDY/0.jpg)](https://www.youtube.com/watch?v=23ziir_IJDY)

#### Youtube - Gameplay Demo - v1.1.0
[![Watch the video](https://img.youtube.com/vi/2flPiJvl4vU/0.jpg)](https://www.youtube.com/watch?v=2flPiJvl4vU)

#### Youtube - Gameplay Demo - v1.0.0
[![Watch the video](https://img.youtube.com/vi/DSeU1Lnglsg/0.jpg)](https://www.youtube.com/watch?v=DSeU1Lnglsg)

### Checklist:
- [X] Create a spaceship placed at mouse pointer on the screen.
- [X] Add controls to the spaceship, `Mouse Motion` left, right, up and down. 
- [X] Shoots laser beam by hitting the `Left Mouse Button` or `Space Bar`.
- [X] Create enemy spaceships (Currently have 10 different types).
- [X] Can deploy aliens ships in a random x position
- [X] Stars falling down with different speed. (just to look cool)
- [X] Adding Start Menu Screen UI
- [X] Adding Pause Menu using `p` to pause the game.
- [X] Adding spaceship UI details.
- [X] Show control info at the bottom of the screen.
- [X] Loot boxes to increase health.
    - [X] Health Kits will be dropped every 1 minute.
    - [X] Can be consumed by pressing [F] key.
- [X] Implement a configuration file.
- [X] Adding Timer.
- [X] Alien Ship Stats on the right of the screen.
- [X] Alien ships can shoot the player's spaceship.
- [X] Implementing a better way to level up the player.
    - [X] With upgrade choice /  stat selection. Player can choose to upgrade either (gun power, speed or capacity).
- [X] Now the game runs without a `config` file. If the file dose not exist, the default configuration will be used.
- [X] Added variety of spaceship selection before starting the game. 
    - [X] Before the game starts, the player will be able to select a ship.
- [X] Redesign UI of the spaceship selection.
- [X] Adding status pop up. Will show to keep the player up to date.
- [X] Add another drop down to increase other spaceship stats. 
- [X] Boss Fights every 3 minutes.
- [X] Gun Cool-down.
- [X] Both the player and the enemy can take damage when crashing the ships.
- [X] Gun Reload + cooldown.
- [X] Adding Asteroids.
- [X] New Particle System: Added simple animation of an explosion
- [X] Added meteoroids — when an asteroid explodes, it scatters meteoroids.
    - They damage spaceships, alien ships, and boss ships.
- [X] Spaceship selection UI Redesign. Now can hover on cards of spaceships and show details.
    - Can select a spaceship with Left Mouse Click.

### Controls

| Control               | Action                                           |
|-----------------------|--------------------------------------------------|
| Left Mouse Pressed    | Shoot beams                                      |
| Mouse Movement        | Move the spaceship                               |
| E                     | Consume health kit (increase spaceship's health) |
| R - Right Mouse Click | Reload Gun                                       |
| P                     | Pause the game                                   |
| Ctrl+R                | Restart game                                     |
| Ctrl+Q                | Quit game                                        |

### Default Configuration File
Configuration file added for the player to freely change/update entity's attributes. The config file saved as `config.toml`.

```toml
[dev]
debug = false
fps_counter = false
asteroids = true

[spaceship]
max_level = 50
next_level_score = 100
gun_max_cap = 40
gun_max_speed = 70

[stars] 
limit = 10
speed = 50
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
