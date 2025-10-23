package design

type AsteroidDesign struct {
	MaxLimit  int      `json:"max_limit"`
	MaxSpeed  int      `json:"max_speed"`
	Asteroids []Design `json:"asteroids"`
}
