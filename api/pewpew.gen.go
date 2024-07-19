// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package api

type EntityType int

const (
  ASTEROID EntityType = iota
  BAF
  INERTIAC
  MOTHERSHIP
  MOTHERSHIP_BULLET
  ROLLING_CUBE
  ROLLING_SPHERE
  UFO
  WARY
  CROWDER
  CUSTOMIZABLE_ENTITY
  SHIP
  BOMB
  BAF_BLUE
  BAF_RED
  WARY_MISSILE
  UFO_BULLET
  PLAYER_BULLET
  BOMB_EXPLOSION
  PLAYER_EXPLOSION
  BONUS
  FLOATING_MESSAGE
  POINTONIUM
  BONUS_IMPLOSION
)

type MothershipType int

const (
  THREE_CORNERS MothershipType = iota
  FOUR_CORNERS
  FIVE_CORNERS
  SIX_CORNERS
  SEVEN_CORNERS
)

type CannonType int

const (
  SINGLE CannonType = iota
  TIC_TOC
  DOUBLE
  TRIPLE
  FOUR_DIRECTIONS
  DOUBLE_SWIPE
  HEMISPHERE
)

type CannonFrequency int

const (
  FREQ_30 CannonFrequency = iota
  FREQ_15
  FREQ_10
  FREQ_7_5
  FREQ_6
  FREQ_5
  FREQ_3
  FREQ_2
  FREQ_1
)

type BombType int

const (
  FREEZE BombType = iota
  REPULSIVE
  ATOMIZE
  SMALL_ATOMIZE
  SMALL_FREEZE
)

type BonusType int

const (
  REINSTANTIATION BonusType = iota
  SHIELD
  SPEED
  WEAPON
)

type WeaponType int

const (
  BULLET WeaponType = iota
  FREEZE_EXPLOSION
  REPULSIVE_EXPLOSION
  ATOMIZE_EXPLOSION
)

type AsteroidSize int

const (
  SMALL AsteroidSize = iota
  MEDIUM
  LARGE
  VERY_LARGE
)

