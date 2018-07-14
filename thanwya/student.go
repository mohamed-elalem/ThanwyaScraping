package thanwya

type Student struct {
	SeatNumber         int `json:",string"`
	Name               string
	TotalDegree        float32 `json:",string"`
	TotalDegreeAfterHL float32 `json:",string"`
	StudentType        string
	NumberOfFailures   int `json:"NoOfFails,string"`
	Section            int `json:",string"`
	Grade
	School
}
