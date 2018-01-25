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
	w.mgr.nodesMutex.Lock()
	first := true
	for k, v := range w.mgr.nodes {
		if !first {
			fmt.Fprintf(resp, ",")
		} else {
			first = false
		}
		v.mutex.RLock()
		fmt.Fprintf(resp, `{"id":%d,"lat":%f,"long":%f,"sn":false}`, k, v.Latitude, v.Longitude)
		v.mutex.RUnlock()
	}
	w.mgr.nodesMutex.Unlock()
	w.mgr.supernodesMutex.Lock()
	for k, v := range w.mgr.supernodes {
		fmt.Fprintf(resp, ",")
		v.mutex.RLock()
		fmt.Fprintf(resp, `{"id":%d,"lat":%f,"long":%f,"sn":true}`, k, v.Latitude, v.Longitude)
		v.mutex.RUnlock()
	}
	fmt.Fprintf(resp, "]}")
	w.mgr.supernodesMutex.Unlock()
}

func (n *node) MarshalJSON() ([]byte, error) {
	return []byte{0}, nil
}
