package nodesim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/abemac/bomb-detection/constants"
)

func ExecConfigFile(file string, managerIP string) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.E(err)
	}
	var config interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.E(err)
	}
	for _, r := range (config.(map[string]interface{}))["nodes"].([]interface{}) {
		row := r.(map[string]interface{})
		north := int(row["north"].(float64))
		south := int(row["south"].(float64))
		east := int(row["east"].(float64))
		west := int(row["west"].(float64))
		num := uint64(row["num"].(float64))
		supernode := row["supernode"].(bool)
		lat := row["lat"].(float64)
		long := row["long"].(float64)
		//group := row["group"].(bool)
		if supernode {
			CreateSupernodes(num, managerIP, lat, long)
		} else {
			CreateNodes(num, managerIP, north, south, east, west, lat, long)
		}

	}
}
func SplitConfigFile(file string) {
	data, err := ioutil.ReadFile(filepath.Join(constants.DIST_PATH, "assets", "uploads", file))
	if err != nil {
		log.E(err)
	}
	var config interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.E(err)
	}
	basename := file[0 : len(file)-5]
	outdir := filepath.Join(os.Getenv("GOPATH"), "/src/github.com/abemac/bomb-detection/run/config-parts/")
	filenum := 0
	nodecount := uint64(0)
	first := true
	outfile, err := os.Create(filepath.Join(outdir, basename+"-"+strconv.Itoa(filenum)+".json"))
	if err != nil {
		log.E(err.Error())
	}
	fmt.Fprintf(outfile, `{"nodes":[`)
	for _, r := range (config.(map[string]interface{}))["nodes"].([]interface{}) {
		row := r.(map[string]interface{})
		num := uint64(row["num"].(float64))
		nodecount += num
		if nodecount < 10000 {
			bytes, err := json.Marshal(row)
			if err != nil {
				log.E(err.Error())
			}
			if !first {
				fmt.Fprintf(outfile, ",")
			} else {
				first = false
			}
			outfile.Write(bytes)
		} else {
			overflow := nodecount - 10000
			part := num - overflow
			row["num"] = part
			bytes, err := json.Marshal(row)
			if err != nil {
				log.E(err.Error())
			}
			if !first {
				fmt.Fprintf(outfile, ",")
			} else {
				first = false
			}
			outfile.Write(bytes)
			fmt.Fprintf(outfile, `]}`)
			outfile.Close()
			filenum++
			first = true
			nodecount = 0

			for overflow > 10000 {
				row["num"] = 10000
				outfile, err = os.Create(filepath.Join(outdir, basename+"-"+strconv.Itoa(filenum)+".json"))
				if err != nil {
					log.E(err.Error())
				}
				fmt.Fprintf(outfile, `{"nodes":[`)
				bytes, err := json.Marshal(row)
				if err != nil {
					log.E(err.Error())
				}
				outfile.Write(bytes)
				fmt.Fprintf(outfile, `]}`)
				outfile.Close()
				filenum++
				overflow -= 10000
			}
			if overflow > 0 {
				row["num"] = overflow
				outfile, err = os.Create(filepath.Join(outdir, basename+"-"+strconv.Itoa(filenum)+".json"))
				if err != nil {
					log.E(err.Error())
				}
				fmt.Fprintf(outfile, `{"nodes":[`)
				bytes, err := json.Marshal(row)
				if err != nil {
					log.E(err.Error())
				}
				outfile.Write(bytes)
				first = false
				nodecount = overflow
			}

		}

	}
	fmt.Fprintf(outfile, `]}`)
	outfile.Close()

}
