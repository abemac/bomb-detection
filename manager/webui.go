package manager

import (
	"fmt"
	"net/http"

	"github.com/abemac/bomb-detection/constants"
)

//WebUI represents the Web interface
type WebUI struct {
	mgr *Manager
}

//NewWebUI creates a new http server for the webserver API and serves the dist directory
func NewWebUI(mgr *Manager) *WebUI {
	w := new(WebUI)
	w.mgr = mgr
	w.setup()
	return w
}

func (w *WebUI) setup() {
	http.Handle("/", http.FileServer(http.Dir(constants.DIST_PATH)))
	log.I("Web UI http server path: ", constants.DIST_PATH)
	http.HandleFunc("/GetNodes", w.handleNodeInfoRequest)
}

//Run the http server
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
