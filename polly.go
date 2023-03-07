package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func main() {
	Execute()
}

// InfluxAsyncGet polls the Tesla Gen3 Wall Connector and writes results to
// InfluxDB with an async, non-blocking client you supply. You must also
// supply the IP of the wall connector.
func InfluxAsyncGet(writeAPI *api.WriteAPI, config Config) {
	client := *writeAPI

	data := getWC3Data(config.Charger)
	if data == nil {
		return
	}

	// Output a dot (.) for every successful GET against the Wall Connector
	// This helps people like me who need to see something to know it works
	fmt.Printf(".")

	p := influxdb2.NewPoint(
		"hpwc",
		map[string]string{
			"product":  "Gen3 HPWC",
			"vendor":   "Tesla",
			"location": config.Charger.Location,
		},
		data,
		time.Now())
	client.WritePoint(p)
}

func getWC3Data(charger Charger) map[string]interface{} {

	var data map[string]interface{}

	resp, err := http.Get(fmt.Sprintf("http://%s/api/1/vitals", charger.IP))
	if err != nil {
		fmt.Printf("error - during GET of hpwc. Do you have the right IP: %s\n", charger.IP)
		return nil
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	return data
}

// Execute simply runs the totality of polly in your program. It is
// recommended you run this as a goroutine so your program can do
// other things.
func Execute() {
	config := loadConfig()

	if hpwcIP, hpwcIpSet := os.LookupEnv("HPWC_IP"); hpwcIpSet {
		fmt.Println("Setting IP based on Env Variable")
		config.Charger.IP = hpwcIP
	}
	if influxIP, influxIpSet := os.LookupEnv("INFLUX_IP"); influxIpSet {
		config.Db.IP = influxIP
	}
	client := influxdb2.NewClientWithOptions(fmt.Sprintf("http://%s:8086", config.Db.IP), config.AuthToken, influxdb2.DefaultOptions().SetBatchSize(20))
	writeAPI := client.WriteAPI("home", "tesla")

	// The way this is set up, these likely don't get executed on ^C.
	defer client.Close()
	defer writeAPI.Flush()

	// Simple, isn't it?
	for {
		go InfluxAsyncGet(&writeAPI, config)
		time.Sleep(time.Millisecond * 1000)
	}
}
