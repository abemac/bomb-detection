package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/nodesim"
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
	http.HandleFunc("/StartSim", w.handleStartSim)
	http.HandleFunc("/StopSim", w.handleStopSim)

}

//Run the http server
func (w *WebUI) Run() {
	http.ListenAndServe(":8080", nil)
}
func simRunning() bool {
	cmd := "docker ps | grep nodesim | wc -l"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.E(err.Error())
	}
	if strings.TrimSpace(string(out)) == "0" {
		return false
	}
	return true
}
func (w *WebUI) handleStopSim(resp http.ResponseWriter, req *http.Request) {

	cmd := "docker ps | grep nodesim | awk '{print $NF}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.E(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(resp, "ERROR:", err.Error())
	} else {
		for _, container := range strings.Split(string(out), "\n") {
			if container != "" {
				log.I(container)
				err := exec.Command("docker", "kill", container).Run()
				if err != nil {
					log.E(err.Error())
					resp.WriteHeader(http.StatusInternalServerError)
					resp.Header().Set("Content-Type", "text/plain")
					fmt.Fprintln(resp, "ERROR:", err.Error())
				}
				resp.Header().Set("Content-Type", "text/plain")
				fmt.Fprintln(resp, "SUCCESS")
			}
		}
		w.mgr.nodesMutex.Lock()
		w.mgr.supernodesMutex.Lock()
		w.mgr.nodes = make(map[uint64]*node)
		w.mgr.supernodes = make(map[uint64]*supernode)
		w.mgr.nodesMutex.Unlock()
		w.mgr.supernodesMutex.Unlock()
	}
}
func (w *WebUI) handleStartSim(resp http.ResponseWriter, req *http.Request) {
	qp := req.URL.Query()
	if filename, ok := qp["filename"]; ok {
		if simRunning() {
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(resp, "ERROR: simulation already running")
			return
		}
		nodesim.SplitConfigFile(filename[0])
		outdir := filepath.Join(os.Getenv("GOPATH"), "/src/github.com/abemac/bomb-detection/run/config-parts/")
		basename := filename[0][0 : len(filename[0])-5]
		files, err := filepath.Glob(filepath.Join(outdir, basename+"-*.json"))

		for index, file := range files {
			log.I(filepath.Base(file))
			cmd := exec.Command("docker", "run", "--rm", "--name", "nodesim"+strconv.Itoa(index),
				"-v", os.Getenv("GOPATH")+":/go", "-e", "GOPATH=/go", "-e", "LOG_LEVEL=4",
				"-e", "NODE_CONFIG_FILE="+filepath.Base(file), "nodesim")
			err := cmd.Start()
			if err != nil {
				log.E(err.Error())
			}
		}
		if err != nil {
			log.E(err)
		}
	} else {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(resp, "ERROR: must specify query parameter filename")
	}
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
			os.Mkdir(filepath.Join(constants.DIST_PATH, "assets"), 0775)
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
		fmt.Fprintf(resp, `{"id":%d,"lat":%f,"long":%f,"sn":false,"bp":%f,"sv":%d}`, k, v.Latitude, v.Longitude, v.BatteryPercentage, v.Value)
		v.mutex.RUnlock()
	}
	w.mgr.nodesMutex.Unlock()
	w.mgr.supernodesMutex.Lock()
	for k, v := range w.mgr.supernodes {
		fmt.Fprintf(resp, ",")
		v.mutex.RLock()
		fmt.Fprintf(resp, `{"id":%d,"lat":%f,"long":%f,"sn":true,"bp":1.0,"sv":%d}`, k, v.Latitude, v.Longitude, v.Value)
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
