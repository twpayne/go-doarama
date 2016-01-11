package doarama

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

var DefaultActivityTypes = ActivityTypes{
	{Id: 0, Name: "Undefined - Ground Based"},
	{Id: 1, Name: "Walk - Fitness"},
	{Id: 2, Name: "Walk - Hike/Trek etc"},
	{Id: 3, Name: "Run - Fitness"},
	{Id: 4, Name: "Cycle - Sport/Road etc"},
	{Id: 5, Name: "Cycle - Transport"},
	{Id: 6, Name: "Cycle - Mountain/Off Road etc"},
	{Id: 7, Name: "Motorcycle"},
	{Id: 8, Name: "Boat - Sail"},
	{Id: 9, Name: "Boat - Motor"},
	{Id: 10, Name: "Boat - Kayak/Canoe/Row etc"},
	{Id: 11, Name: "Fly - Glide"},
	{Id: 12, Name: "Fly - Aircraft"},
	{Id: 13, Name: "Ski - Cross Country"},
	{Id: 14, Name: "Ski - Downhill"},
	{Id: 15, Name: "Ski - Roller"},
	{Id: 16, Name: "Ski - Waterski/Wakeboard etc"},
	{Id: 17, Name: "Snowboard"},
	{Id: 18, Name: "Skate - Ice"},
	{Id: 19, Name: "Skate - Roller/Skateboard/Scooter etc"},
	{Id: 20, Name: "Surf - Wave"},
	{Id: 21, Name: "Surf - Kite"},
	{Id: 22, Name: "Surf - Windsurf"},
	{Id: 23, Name: "Swim"},
	{Id: 24, Name: "Drive - Car/Truck/Bus etc"},
	{Id: 25, Name: "Rail - Train"},
	{Id: 26, Name: "Ride - Equestrian/Camel etc"},
	{Id: 27, Name: "Fly - Hang Glide"},
	{Id: 28, Name: "Fly - Sailplane / Glider"},
	{Id: 29, Name: "Fly - Paraglide"},
	{Id: 30, Name: "Fly - UAV / Drone"},
	{Id: 31, Name: "Fly - Bird / Aves"},
	{Id: 32, Name: "Skate - Windskate"},
	{Id: 33, Name: "Undefined - Aerial"},
	{Id: 34, Name: "Fly - Balloon"},
	{Id: 35, Name: "Fly - Hike + Glide"},
}

//go:generate go run cmd/generate-activity-types/generate-activity-types.go -o activity_types.go
