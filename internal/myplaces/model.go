package myplaces

type MyPlaces struct {
	AflId                     int64   `json:"afl_id"`
	AccountId                 int64   `json:"account_id"`
	Title                     string  `json:"title"`
	Latitude                  float32 `json:"latitude"`
	Longitude                 float32 `json:"longitude"`
	City                      string  `json:"city"`
	Province                  string  `json:"province"`
	PostalCode                string  `json:"postal_code"`
	Status                    bool    `json:"status"`
	Crime                     bool    `json:"crime"`
	TrafficAccident           bool    `json:"traffic_accident"`
	MedicalEmergency          bool    `json:"medical_emergency"`
	FireIncident              bool    `json:"fire_incident"`
	Vandalism                 bool    `json:"vandalism"`
	SuspiciousActivity        bool    `json:"suspicious_activity"`
	InfrastructureIssues      bool    `json:"infrastructure_issues"`
	ExtremeWeather            bool    `json:"extreme_weather"`
	CommunityEvents           bool    `json:"community_events"`
	DangerousWildlifeSighting bool    `json:"dangerous_wildlife_sighting"`
	PositiveActions           bool    `json:"positive_actions"`
	LostPet                   bool    `json:"lost_pet"`
	Radius                    int     `json:"radius"`
}
