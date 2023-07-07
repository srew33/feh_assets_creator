package tasks

const (
	Move int = iota
	Weapon
	Person
	Skill
	WeaponRefine
	SkillAccessory
	SubscriptionCostume
	HolyGrail
	TransEn
	TransJp
	TransCn
)

var FileTasks = map[int]string{
	Move:   "files/assets/Common/SRPG/Move.json",
	Weapon: "files/assets/Common/SRPG/Weapon.json",
}
var DirectoryTasks = map[int]string{
	Person:              "files/assets/Common/SRPG/Person/",
	Skill:               "files/assets/Common/SRPG/Skill/",
	WeaponRefine:        "files/assets/Common/SRPG/WeaponRefine/",
	SkillAccessory:      "files/assets/Common/SRPG/SkillAccessory/",
	SubscriptionCostume: "files/assets/Common/SubscriptionCostume/",
	HolyGrail:           "files/assets/Common/HolyGrail/",
	TransJp:             "files/assets/JPJA/Message/Data/",
	TransEn:             "files/assets/USEN/Message/Data/",
	TransCn:             "files/assets/TWZH/Message/Data/",
}
