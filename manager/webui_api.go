package manager

import (
	"fmt"
	"net/http"
	"os"
)

type WebUI struct {
	mgr *Manager
}

func NewWebUI(mgr *Manager) *WebUI {
	w := new(WebUI)
	w.mgr = mgr
	w.setup()
	return w
}

func (w *WebUI) setup() {
	http.Handle("/", http.FileServer(http.Dir(os.Getenv("GOPATH")+"/src/github.com/abemac/bomb-detection/webapp/manager/dist/")))
	http.HandleFunc("/GetNodes", w.handleNodeInfoRequest)

}

func (w *WebUI) Run() {
	http.ListenAndServe(":8080", nil)
}

func (w *WebUI) handleNodeInfoRequest(resp http.ResponseWriter, req *http.Request) {
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
