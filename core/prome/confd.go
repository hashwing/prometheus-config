package prome

import (
	"path/filepath"
	"os"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)


// Config prome config
type Config struct{
	ConfigPath			string
	ScrapeInterval 		string
	EvaluationInterval 	string
	ScrapeTimeout		string
	RulesPath 			string
	RemoteR  			string
	RemoteW				string
	AlertManager		bool
	LabelKey			string
	LabelValue			string				
	Job					*Job
	ShardsSum 			int
	ShardsNum			int
}

// Remote remote r and w
type Remote struct{
	Enabled 	bool
	Url 	string
}

// Job scrape jobs
type Job struct{
	Local 		bool
	Endpoints 	bool
	PushGateway bool
	Service 	bool
	Pod 		bool
	ApiServers 	bool
	Node 		bool
	Cadvisor	bool
}

// CreateConfig create prometheus confgi
func CreateConfig( c Config)error{
	tmpl, err := template.New("config").Parse(promConfigTpl)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(c.ConfigPath),0666)
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(c.ConfigPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	err = tmpl.Execute(fd, c)
	if err != nil {
		return err
	}
	return nil
}

// ReloadEndpoint reload service
func ReloadEndpoint(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to reload Prometheus config: %s", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("Unable to reload the Prometheus config. Endpoint: %s, Reponse StatusCode: %d, Response Body: %s", url, resp.StatusCode, string(respBody))
}