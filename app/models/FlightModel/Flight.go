package FlightModel

type Flight struct {
	Id             	int
	IdPlane			int
	PlaneName 		string
	IdPilot			[]int
	PilotName		[]string
	ArrivalPoint   	string
	DeparturePoint 	string
}
