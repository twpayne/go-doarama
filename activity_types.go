package doarama

// Activity types
const (
	BoatCanoe            = 10
	BoatKayak            = 10
	BoatMotor            = 9
	BoatRow              = 10
	BoatSail             = 8
	CycleMountain        = 6
	CycleOffRoad         = 6
	CycleRoad            = 4
	CycleSport           = 4
	CycleTransport       = 5
	DriveBus             = 24
	DriveCar             = 24
	DriveTruck           = 24
	FlyAircraft          = 12
	FlyAves              = 31
	FlyBalloon           = 34
	FlyBird              = 31
	FlyDrone             = 30
	FlyGlide             = 11
	FlyGlider            = 28
	FlyHangGlide         = 27
	FlyHikeAndGlide      = 35
	FlyParaglide         = 29
	FlySailplane         = 28
	FlyUAV               = 30
	Motorcycle           = 7
	RailTrain            = 25
	RideCamel            = 26
	RideEquestrian       = 26
	RunFitness           = 3
	SkateIce             = 18
	SkateRoller          = 19
	SkateScooter         = 19
	SkateSkateboard      = 19
	SkateWindskate       = 32
	SkiCrossCountry      = 13
	SkiDownhill          = 14
	SkiRoller            = 15
	SkiWakeboard         = 16
	SkiWaterski          = 16
	Snowboard            = 17
	SurfKite             = 21
	SurfWave             = 20
	SurfWindsurf         = 22
	Swim                 = 23
	UndefinedAerial      = 33
	UndefinedGroundBased = 0
	WalkFitness          = 1
	WalkHike             = 2
	WalkTrek             = 2
)

// DefaultActivityTypes contains the default activity types.
var DefaultActivityTypes = ActivityTypes{
	{ID: 0, Name: "Undefined - Ground Based"},
	{ID: 1, Name: "Walk - Fitness"},
	{ID: 2, Name: "Walk - Hike/Trek etc"},
	{ID: 3, Name: "Run - Fitness"},
	{ID: 4, Name: "Cycle - Sport/Road etc"},
	{ID: 5, Name: "Cycle - Transport"},
	{ID: 6, Name: "Cycle - Mountain/Off Road etc"},
	{ID: 7, Name: "Motorcycle"},
	{ID: 8, Name: "Boat - Sail"},
	{ID: 9, Name: "Boat - Motor"},
	{ID: 10, Name: "Boat - Kayak/Canoe/Row etc"},
	{ID: 11, Name: "Fly - Glide"},
	{ID: 12, Name: "Fly - Aircraft"},
	{ID: 13, Name: "Ski - Cross Country"},
	{ID: 14, Name: "Ski - Downhill"},
	{ID: 15, Name: "Ski - Roller"},
	{ID: 16, Name: "Ski - Waterski/Wakeboard etc"},
	{ID: 17, Name: "Snowboard"},
	{ID: 18, Name: "Skate - Ice"},
	{ID: 19, Name: "Skate - Roller/Skateboard/Scooter etc"},
	{ID: 20, Name: "Surf - Wave"},
	{ID: 21, Name: "Surf - Kite"},
	{ID: 22, Name: "Surf - Windsurf"},
	{ID: 23, Name: "Swim"},
	{ID: 24, Name: "Drive - Car/Truck/Bus etc"},
	{ID: 25, Name: "Rail - Train"},
	{ID: 26, Name: "Ride - Equestrian/Camel etc"},
	{ID: 27, Name: "Fly - Hang Glide"},
	{ID: 28, Name: "Fly - Sailplane / Glider"},
	{ID: 29, Name: "Fly - Paraglide"},
	{ID: 30, Name: "Fly - UAV / Drone"},
	{ID: 31, Name: "Fly - Bird / Aves"},
	{ID: 32, Name: "Skate - Windskate"},
	{ID: 33, Name: "Undefined - Aerial"},
	{ID: 34, Name: "Fly - Balloon"},
	{ID: 35, Name: "Fly - Hike + Glide"},
}

//go:generate go run cmd/generate-activity-types/generate-activity-types.go -o activity_types.go
