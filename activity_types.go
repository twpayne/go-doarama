//go:generate go run cmd/generate-activity-types/generate-activity-types.go -o activity_types.go
package doarama

const (
	BOAT_CANOE             = 10
	BOAT_KAYAK             = 10
	BOAT_MOTOR             = 9
	BOAT_ROW               = 10
	BOAT_SAIL              = 8
	CYCLE_MOUNTAIN         = 6
	CYCLE_OFF_ROAD         = 6
	CYCLE_ROAD             = 4
	CYCLE_SPORT            = 4
	CYCLE_TRANSPORT        = 5
	DRIVE_BUS              = 24
	DRIVE_CAR              = 24
	DRIVE_TRUCK            = 24
	FLY_AIRCRAFT           = 12
	FLY_AVES               = 31
	FLY_BALLOON            = 34
	FLY_BIRD               = 31
	FLY_DRONE              = 30
	FLY_GLIDE              = 11
	FLY_GLIDER             = 28
	FLY_HANG_GLIDE         = 27
	FLY_PARAGLIDE          = 29
	FLY_SAILPLANE          = 28
	FLY_UAV                = 30
	MOTORCYCLE             = 7
	RAIL_TRAIN             = 25
	RIDE_CAMEL             = 26
	RIDE_EQUESTRIAN        = 26
	RUN_FITNESS            = 3
	SKATE_ICE              = 18
	SKATE_ROLLER           = 19
	SKATE_SCOOTER          = 19
	SKATE_SKATEBOARD       = 19
	SKATE_WINDSKATE        = 32
	SKI_CROSS_COUNTRY      = 13
	SKI_DOWNHILL           = 14
	SKI_ROLLER             = 15
	SKI_WAKEBOARD          = 16
	SKI_WATERSKI           = 16
	SNOWBOARD              = 17
	SURF_KITE              = 21
	SURF_WAVE              = 20
	SURF_WINDSURF          = 22
	SWIM                   = 23
	UNDEFINED_AERIAL       = 33
	UNDEFINED_GROUND_BASED = 0
	WALK_FITNESS           = 1
	WALK_HIKE              = 2
	WALK_TREK              = 2
)

var (
	ActivityTypes = map[string]int{
		"Boat - Kayak/Canoe/Row etc":            10,
		"Boat - Motor":                          9,
		"Boat - Sail":                           8,
		"Cycle - Mountain/Off Road etc":         6,
		"Cycle - Sport/Road etc":                4,
		"Cycle - Transport":                     5,
		"Drive - Car/Truck/Bus etc":             24,
		"Fly - Aircraft":                        12,
		"Fly - Balloon":                         34,
		"Fly - Bird / Aves":                     31,
		"Fly - Glide":                           11,
		"Fly - Hang Glide":                      27,
		"Fly - Paraglide":                       29,
		"Fly - Sailplane / Glider":              28,
		"Fly - UAV / Drone":                     30,
		"Motorcycle":                            7,
		"Rail - Train":                          25,
		"Ride - Equestrian/Camel etc":           26,
		"Run - Fitness":                         3,
		"Skate - Ice":                           18,
		"Skate - Roller/Skateboard/Scooter etc": 19,
		"Skate - Windskate":                     32,
		"Ski - Cross Country":                   13,
		"Ski - Downhill":                        14,
		"Ski - Roller":                          15,
		"Ski - Waterski/Wakeboard etc":          16,
		"Snowboard":                             17,
		"Surf - Kite":                           21,
		"Surf - Wave":                           20,
		"Surf - Windsurf":                       22,
		"Swim":                                  23,
		"Undefined - Aerial":                    33,
		"Undefined - Ground Based":              0,
		"Walk - Fitness":                        1,
		"Walk - Hike/Trek etc":                  2,
	}
)