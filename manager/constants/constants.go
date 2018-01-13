package constants

//ManagerToNodeJSON defines the JSON format of Manager to Node communication
type ManagerToNodeJSON struct {
	NextCheckin   int  `json:"nextCheckin"`
	PerformSample bool `json:"performSample"`
}

//NodeToManagerJSON defines the JSON format of Node to Manager communication
type NodeToManagerJSON struct {
	SampleValue int     `json:"sampleValue"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude`
}
