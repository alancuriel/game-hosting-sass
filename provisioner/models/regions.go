package models

type Region uint8

const (
	INVALID      Region = 0
	US_SEA       Region = 1
	US_WEST      Region = 2
	US_EAST      Region = 3
	US_CENTRAL   Region = 4
	US_SOUTHEAST Region = 5
	US_IAD       Region = 6
)

func (r Region) String() string {
	return [...]string{
		"invalid",
		"us-sea",
		"us-west",
		"us-east",
		"us-central",
		"us-southeast",
		"us-iad",
	}[r]
}
