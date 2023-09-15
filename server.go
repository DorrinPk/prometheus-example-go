package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	http.HandleFunc("/metrics", metricsHandler)

	log.Println("Server started on 8080")
	http.ListenAndServe(":8080", nil)

}

func metricsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	cmd := exec.Command("/usr/local/bin/metrics")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		return
	}

	defer cmd.Wait()

	metricsData, err := io.ReadAll(stdout)
	if err != nil {
		log.Println("Error reading from stdout:", err)
		return
	}

	reg := prometheus.NewRegistry()

	totalIPsInUse := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_ips_in_use",
			Help: "Total IPs in use, labeled with the MAC address and hostname",
		},
		[]string{"mac_address", "hostname"},
	)

	reg.MustRegister(totalIPsInUse)

	lines := strings.Split(string(metricsData), "\n")
	for _, line := range lines {
		field := strings.Fields(line)
		if len(field) >= 3 {
			host := field[0]
			ipCountString := field[1]

			hostInfo := strings.Split(host, "(")
			if len(hostInfo) != 2 {
				log.Println("Error parsing MAC address and hostname: ", line)
				continue
			}

			macAddress := hostInfo[0]
			hostname := strings.TrimSuffix(hostInfo[1], ")")

			ipCount, err := strconv.Atoi(ipCountString)
			if err != nil {
				log.Println("error parsing IP count:", err)
				continue
			}

			totalIPsInUse.WithLabelValues(macAddress, hostname).Set(float64(ipCount))
		}
	}
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	handler.ServeHTTP(w, r)

}
