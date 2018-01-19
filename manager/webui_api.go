package manager

import (
	"fmt"
	"net/http"
)

type WebAPI struct {
	mgr *Manager
}

func NewWebAPI(mgr *Manager) *WebAPI {
	w := new(WebAPI)
	w.mgr = mgr
	w.setup()
	return w
}

// func (w *WebAPI) handleMessage(bytes []byte) []byte {
// 	//message := new(WebUIReq)
// 	json.Unmarshal(bytes[:len(bytes)-1], message)
// 	log.D("Received: ", *message)

// 	if message.InfoWanted == "nodes" {
// 		response := new(WebUIResp)
// 		response.Nodes = map[uint64]*node{
// 			0: nil,
// 			2: nil,
// 			3: nil,
// 			4: nil,
// 		}
// 		responseBytes, err := json.Marshal(response)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		log.D("Sending: ", *response)
// 		return responseBytes
// 	}
// 	return nil
// }

func (w *WebAPI) setup() {
	http.Handle("/", http.FileServer(http.Dir("./webapp/manager/dist/")))
	http.HandleFunc("/GetNodes", w.handleNodeInfoRequest)

}

func (w *WebAPI) Run() {
	http.ListenAndServe(":8080", nil)
}

func (w *WebAPI) handleNodeInfoRequest(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, `{"nodes":[`)
	w.mgr.mapMutex.RLock()
	first := true
	for k, v := range w.mgr.nodes {
		if !first {
			fmt.Fprintf(resp, ",")
		} else {
			first = false
		}
		fmt.Fprintf(resp, `{"id":%d,"latitude":%f,"longitude":%f}`, k, v.Latitude, v.Longitude)
	}
	fmt.Fprintf(resp, "]}")
	w.mgr.mapMutex.RUnlock()
}

func (n *node) MarshalJSON() ([]byte, error) {
	return []byte{0}, nil
}
