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
	FLY_HIKE_AND_GLIDE     = 35
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

var DefaultActivityTypes = ActivityTypes{
	{
		Id:   0,
		Name: "Undefined - Ground Based",
	},
	{
		Id:   1,
		Name: "Walk - Fitness",
	},
	{
		Id:   2,
		Name: "Walk - Hike/Trek etc",
	},
	{
		Id:   3,
		Name: "Run - Fitness",
	},
	{
		Id:   4,
		Name: "Cycle - Sport/Road etc",
	},
	{
		Id:   5,
		Name: "Cycle - Transport",
	},
	{
		Id:   6,
		Name: "Cycle - Mountain/Off Road etc",
	},
	{
		Id:   7,
		Name: "Motorcycle",
	},
	{
		Id:   8,
		Name: "Boat - Sail",
	},
	{
		Id:   9,
		Name: "Boat - Motor",
	},
	{
		Id:   10,
		Name: "Boat - Kayak/Canoe/Row etc",
	},
	{
		Id:   11,
		Name: "Fly - Glide",
	},
	{
		Id:   12,
		Name: "Fly - Aircraft",
	},
	{
		Id:   13,
		Name: "Ski - Cross Country",
	},
	{
		Id:   14,
		Name: "Ski - Downhill",
	},
	{
		Id:   15,
		Name: "Ski - Roller",
	},
	{
		Id:   16,
		Name: "Ski - Waterski/Wakeboard etc",
	},
	{
		Id:   17,
		Name: "Snowboard",
	},
	{
		Id:   18,
		Name: "Skate - Ice",
	},
	{
		Id:   19,
		Name: "Skate - Roller/Skateboard/Scooter etc",
	},
	{
		Id:   20,
		Name: "Surf - Wave",
	},
	{
		Id:   21,
		Name: "Surf - Kite",
	},
	{
		Id:   22,
		Name: "Surf - Windsurf",
	},
	{
		Id:   23,
		Name: "Swim",
	},
	{
		Id:   24,
		Name: "Drive - Car/Truck/Bus etc",
	},
	{
		Id:   25,
		Name: "Rail - Train",
	},
	{
		Id:   26,
		Name: "Ride - Equestrian/Camel etc",
	},
	{
		Id:   27,
		Name: "Fly - Hang Glide",
	},
	{
		Id:   28,
		Name: "Fly - Sailplane / Glider",
	},
	{
		Id:   29,
		Name: "Fly - Paraglide",
	},
	{
		Id:   30,
		Name: "Fly - UAV / Drone",
	},
	{
		Id:   31,
		Name: "Fly - Bird / Aves",
	},
	{
		Id:   32,
		Name: "Skate - Windskate",
	},
	{
		Id:   33,
		Name: "Undefined - Aerial",
	},
	{
		Id:   34,
		Name: "Fly - Balloon",
	},
	{
		Id:   35,
		Name: "Fly - Hike + Glide",
	},
}

//go:generate go run cmd/generate-activity-types/generate-activity-types.go -o activity_types.go
