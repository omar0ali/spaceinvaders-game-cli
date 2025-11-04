package design

type AbilityEffect struct {
	PowerIncrease          int `json:"power_increase"`
	SpeedIncrease          int `json:"speed_increase"`
	CapacityIncrease       int `json:"capacity_increase"`
	CooldownDecrease       int `json:"cooldown_decrease"`
	ReloadCooldownDecrease int `json:"reload_cooldown_decrease"`
	HealthCpacity          int `json:"health_capacity_increase"`
	MaxValue               int `json:"max_value"`
}

type AbilityDesign struct {
	Name        string        `json:"name"`
	Shape       []string      `json:"shape"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Effect      AbilityEffect `json:"effect"`
}

type SpaceshipDesign struct {
	Design
	GunPower          int `json:"gun_power"`
	GunSpeed          int `json:"gun_speed"`
	GunCap            int `json:"gun_cap"`
	GunCooldown       int `json:"gun_cooldown"`
	GunReloadCooldown int `json:"gun_reload_cooldown"`
}
