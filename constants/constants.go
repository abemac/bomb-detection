package constants

//ManagerToNodeJSON defines the JSON format of Manager to Node communication
type ManagerToNodeJSON struct {
	NextCheckin    int    `json:"nc"`
	PerformSample  bool   `json:"ps"`
	AssignedID     uint64 `json:"id"`
	BatteryPercent uint8  `json:"bp"`
}

//NodeToManagerJSON defines the JSON format of Node to Manager communication
type NodeToManagerJSON struct {
	Latitude    float64 `json:"la,omitempty"`
	Longitude   float64 `json:"lo,omitempty"`
	ID          uint64  `json:"id"`
	SampleValue int     `json:"s,omitempty"`
	SampleValid bool    `json:"v,omitempty"`
	SuperNode   bool    `json:"sn,omitempty"`
}

//DelimJSON delimiter used for json messages
var DelimJSON = []byte{0}

const ID_NOT_ASSIGNED uint64 = 0

var LOG_LEVEL int
