package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

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
	http.HandleFunc("/UploadConfig", w.handleUpload)
	http.HandleFunc("/GetConfig", w.handleGetConfig)

}

//Run the http server
func (w *WebUI) Run() {
	http.ListenAndServe(":8080", nil)
}
func (w *WebUI) handleGetConfig(resp http.ResponseWriter, req *http.Request) {
	qp := req.URL.Query()
	if filename, ok := qp["filename"]; ok {
		//Return file
		data, err := ioutil.ReadFile(filepath.Join(constants.DIST_PATH, "assets", "uploads", filename[0]))
		if err != nil {
			log.E(err)
		}
		data, err = prettyprint(data)
		if err != nil {
			log.E(err)
		}
		fmt.Fprintf(resp, string(data))
	} else {
		//Return list of files
		files, err := filepath.Glob(filepath.Join(constants.DIST_PATH, "assets", "uploads", "*"))
		if err != nil {
			log.E(err)
		}
		fmt.Fprintf(resp, `{"files":[`)
		var first = true
		for _, file := range files {
			if !first {
				fmt.Fprintf(resp, ",")
			} else {
				first = false
			}
			fmt.Fprintf(resp, `"%s"`, filepath.Base(file))
		}
		fmt.Fprintf(resp, "]}")

	}

}
func (w *WebUI) handleUpload(resp http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		data := make([]byte, req.ContentLength)
		req.Body.Read(data)
		fmt.Println(string(data))
		_, err := os.Stat(filepath.Join(constants.DIST_PATH, "assets", "uploads"))
		if err != nil {
			os.Mkdir(filepath.Join(constants.DIST_PATH, "assets", "uploads"), 0775)
		}
		_, err = os.Stat(filepath.Join(constants.DIST_PATH, "assets", "uploads", req.Header["Filename"][0]))
		if err == nil {
			log.E("File already exists")
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(resp, "ERROR: file already exists. Please Rename")

		} else {
			file, err := os.Create(filepath.Join(constants.DIST_PATH, "assets", "uploads", req.Header["Filename"][0]))
			if err != nil {
				log.E(err)
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Header().Set("Content-Type", "text/plain")
				fmt.Fprintln(resp, err, "ERROR uploading config file")
			} else {
				_, err = file.Write(data)
				err = file.Close()
				if err == nil {
					resp.WriteHeader(http.StatusAccepted)
					resp.Header().Set("Content-Type", "text/plain")
					fmt.Fprintln(resp, "SUCCESS")
				} else {
					resp.WriteHeader(http.StatusInternalServerError)
					resp.Header().Set("Content-Type", "text/plain")
					fmt.Fprintln(resp, err, "ERROR uploading config file")
				}
			}

		}
		return
	}
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
		fmt.Fprintf(resp, `{"id":%d,"lat":%f,"long":%f,"sn":false,"bp":%f}`, k, v.Latitude, v.Longitude, v.BatteryPercentage)
		v.mutex.RUnlock()
	}
	w.mgr.nodesMutex.Unlock()
	w.mgr.supernodesMutex.Lock()
	for k, v := range w.mgr.supernodes {
		fmt.Fprintf(resp, ",")
		v.mutex.RLock()
		fmt.Fprintf(resp, `{"id":%d,"lat":%f,"long":%f,"sn":true,"bp":1.0}`, k, v.Latitude, v.Longitude)
		v.mutex.RUnlock()
	}
	fmt.Fprintf(resp, "]}")
	w.mgr.supernodesMutex.Unlock()
}
func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
