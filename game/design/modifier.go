package design

type ModifierDesign struct {
	Design
	ModifyHealth            int  `json:"modify_health"`
	ModifyLevel             bool `json:"modify_level"`
	ModifyGunPower          int  `json:"modify_gun_power"`
	ModifyGunCap            int  `json:"modify_gun_cap"`
	ModifyGunSpeed          int  `json:"modify_gun_speed"`
	ModifyGunCoolDown       int  `json:"modify_gun_cooldown"`
	ModifyGunReloadCoolDown int  `json:"modify_gun_reload_cooldown"`
	MaxValue                int  `json:"max_value"`
}
