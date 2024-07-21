package walker

import "hybroid/ast"

var variables = map[string]*VariableVal{
	//enums
	"EntityType": {
		Name:    "EntityType",
		Value:   EntityType,
		IsLocal: false,
		IsConst: true,
	},
	"MothershipType": {
		Name:    "MothershipType",
		Value:   MothershipType,
		IsLocal: false,
		IsConst: true,
	},
	"CannonType": {
		Name:    "CannonType",
		Value:   CannonType,
		IsLocal: false,
		IsConst: true,
	},
	"CannonFrequency": {
		Name:    "CannonFrequency",
		Value:   CannonFrequency,
		IsLocal: false,
		IsConst: true,
	},
	"BombType": {
		Name:    "BombType",
		Value:   BombType,
		IsLocal: false,
		IsConst: true,
	},
	"BonusType": {
		Name:    "BonusType",
		Value:   BonusType,
		IsLocal: false,
		IsConst: true,
	},
	"WeaponType": {
		Name:    "WeaponType",
		Value:   WeaponType,
		IsLocal: false,
		IsConst: true,
	},
	"AsteroidSize": {
		Name:    "AsteroidSize",
		Value:   AsteroidSize,
		IsLocal: false,
		IsConst: true,
	},

	//functions
	"NewCustomizableEntity": {
		Name:  "NewCustomizableEntity",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
}

var EntityType = NewEnumVal("EntityType", false,
	"Asteroid",
	"Baf",
	"Inertiac",
	"Mothership",
	"MothershipBullet",
	"RollingCube",
	"RollingSphere",
	"Ufo",
	"Wary",
	"Crowder",
	"CustomizableEntity",
	"Ship",
	"Bomb",
	"BafBlue",
	"BafRed",
	"WaryMissile",
	"UfoBullet",
	"PlayerBullet",
	"BombExplosion",
	"PlayerExplosion",
	"Bonus",
	"FloatingMessage",
	"Pointonium",
	"BonusImplosion",
)

var MothershipType = NewEnumVal("MothershipType", false,
	"Triangle",
	"Square",
	"Pentagon",
	"Hexagon",
	"Septagon",
)

var CannonType = NewEnumVal("CannonType", false,
	"SinglFixed",
	"TicToc",
	"Double",
	"Triple",
	"FourDirections",
	"DoubleSwipe",
	"Hemisphere",
)

var CannonFrequency = NewEnumVal("CannonFrequency", false,
	"Freq30",
	"Freq15",
	"Freq10",
	"Freq7",
	"Freq6",
	"Freq5",
	"Freq3",
	"Freq2",
	"Freq1",
)

var BombType = NewEnumVal("BombType", false,
	"Freeze",
	"Repulsive",
	"Atomize",
	"SmallAtomize",
	"SmallFreeze",
)

var BonusType = NewEnumVal("BonusType", false,
	"Reinstantiation",
	"Shield",
	"Speed",
	"Weapon",
)

var WeaponType = NewEnumVal("WeaponType", false,
	"Bullet",
	"FreezeExplosion",
	"RepulsiveExplosion",
	"AtomizeExplosion",
)

var AsteroidSize = NewEnumVal("AsteroidSize", false,
	"Small",
	"Medium",
	"Large",
	"VeryLarge",
)