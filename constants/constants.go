package constants

//ManagerToNodeJSON defines the JSON format of Manager to Node communication
type ManagerToNodeJSON struct {
	NextCheckin   int     `json:"nc"`
	PerformSample bool    `json:"ps"`
	AssignedID    uint64  `json:"id"`
	ManagerUID    int64   `json:"muid"`
	GoToLat       float64 `json:"gtla,omitempty"`
	GoToLong      float64 `json:"gtlo,omitempty"`
}

//NodeToManagerJSON defines the JSON format of Node to Manager communication
type NodeToManagerJSON struct {
	Latitude       float64 `json:"la,omitempty"`
	Longitude      float64 `json:"lo,omitempty"`
	ID             uint64  `json:"id"`
	SampleValue    int     `json:"s,omitempty"`
	SampleValid    bool    `json:"v,omitempty"`
	SuperNode      bool    `json:"sn,omitempty"`
	BatteryPercent float32 `json:"bp"`
	ManagerUID     int64   `json:"muid"`
}

//DelimJSON delimiter used for json messages
var DelimJSON = []byte{0}

const ID_NOT_ASSIGNED uint64 = 0

var LOG_LEVEL int

var DIST_PATH string
