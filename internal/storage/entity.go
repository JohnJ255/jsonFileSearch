package storage

type Celesty struct {
	OrbitType              string  `json:"Orbit_type"`
	ProvisionalPackedDesig string  `json:"Provisional_packed_desig"`
	YearOfPerihelion       int     `json:"Year_of_perihelion"`
	MonthOfPerihelion      int     `json:"Month_of_perihelion"`
	DayOfPerihelion        float64 `json:"Day_of_perihelion"`
	PerihelionDist         float64 `json:"Perihelion_dist"`
	E                      float64 `json:"e"`
	Peri                   float64 `json:"Peri"`
	Node                   float64 `json:"Node"`
	I                      float64 `json:"i"`
	EpochYear              int     `json:"Epoch_year"`
	EpochMonth             int     `json:"Epoch_month"`
	EpochDay               int     `json:"Epoch_day"`
	H                      float64 `json:"H"`
	G                      float64 `json:"G"`
	DesignationAndName     string  `json:"Designation_and_name"`
	Ref                    string  `json:"Ref"`
}
