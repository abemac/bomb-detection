package constants

//ManagerToNodeJSON defines the JSON format of Manager to Node communication
type ManagerToNodeJSON struct {
	NextCheckin   int    `json:"nextCheckin"`
	PerformSample bool   `json:"performSample"`
	AssignedID    uint64 `json:"assignedID"`
}

//NodeToManagerJSON defines the JSON format of Node to Manager communication
type NodeToManagerJSON struct {
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	ID          uint64  `json:"assignedID"`
	SampleValue int     `json:"sampleValue,omitempty"`
	SampleValid bool    `json:"sampleValid"`
}

//DelimJSON delimiter used for json messages
var DelimJSON = []byte{0}

const ID_NOT_ASSIGNED uint64 = 0

var LOG_LEVEL int
