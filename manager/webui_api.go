package manager

import (
	"encoding/json"
)

type WebUIReq struct {
	InfoWanted string `json:"infoWanted"`
}
type WebUIResp struct {
	Nodes map[uint64]*node `json:"nodes"`
}
type WebAPI struct {
	mgr *Manager
}

func NewWebAPI(mgr *Manager) *WebAPI {
	w := new(WebAPI)
	w.mgr = mgr
	return w
}
func (w *WebAPI) handleMessage(bytes []byte) []byte {
	message := new(WebUIReq)
	json.Unmarshal(bytes[:len(bytes)-1], message)
	log.D("Received: ", *message)

	if message.InfoWanted == "nodes" {
		response := new(WebUIResp)
		response.Nodes = map[uint64]*node{
			0: nil,
			2: nil,
			3: nil,
			4: nil,
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			panic(err.Error())
		}
		log.D("Sending: ", *response)
		return responseBytes
	}
	return nil
}
