package route

import (
	"cr-product/conf"
	"cr-product/internal/app/model"
	"cr-product/internal/utils"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type HealthCheck struct {
	ServiceName string    `json:"service_name"`
	Version     string    `json:"version"`
	Hostname    string    `json:"hostname"`
	Timelife    time.Time `json:"time_life"`
}

var onceCategory = sync.Once{}
var singleton *HealthCheck

func NewService() {
	onceCategory.Do(func() {
		hostname, _ := os.Hostname()
		singleton = &HealthCheck{
			ServiceName: utils.APPNAME,
			Version:     utils.VERSION,
			Timelife:    time.Now(),
			Hostname:    hostname,
		}
	})

	r := mux.NewRouter()

	r.HandleFunc("/status", CheckHealth).Methods("GET")

	srv := &http.Server{
		Addr:    "0.0.0.0:" + conf.LoadEnv().Port,
		Handler: r,
	}
	err := srv.ListenAndServe()
	if err != nil {
		utils.Log(utils.ERROR_LOG, "", err, "")
	}
}

func CheckHealth(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rsp := model.HealthCheckResponse{
		ServiceName: singleton.ServiceName,
		Version:     singleton.Version,
		Hostname:    singleton.Hostname,
		Timelife:    time.Since(singleton.Timelife).String(),
	}
	json.NewEncoder(rw).Encode(rsp)
}