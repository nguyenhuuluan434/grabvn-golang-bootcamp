package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
Requirement:
Show Grab service details
- Get activated service ids
- Get service details
- Check if subtext is required

GET http://my-json-server.typicode.com/havinhdat/restful-db/services/{id}
{
	"id": $int64,
	"name": $string,
	"subText": $string
}
GET http://my-json-server.typicode.com/havinhdat/restful-db/config/{id}
{
	"id": $string,
	"value": $interface
}
activatedServiceIDs: value is slice of service ids
allowedSubTextServiceIDs: value is slice of service ids
*/

func main() {
	activatedServiceIDs, err := getActivatedServiceIDs()
	if err != nil {
		fmt.Println("<error>")
		return
	}
	services := make([]*Service, 0, len(activatedServiceIDs))
	for _, serviceID := range activatedServiceIDs {
		s, err := getServiceByID(serviceID)
		if err != nil {
			log.Println("WARN: failed to get service by id with error: ", err)
			continue
		}
		services = append(services, s)
	}
	enabledSubTextServiceIDs, err := getEnabledSubTextServiceIDs()
	if err != nil {
		log.Println("WARN: failed to get enabled subtext service ids with error:", err)
		return
	}
	enabledSubTextByServiceID := map[int64]bool{}
	for _, serviceID := range enabledSubTextServiceIDs {
		enabledSubTextByServiceID[serviceID] = true
	}

	for _, service := range services {
		subText := "<omitted>"
		if enabledSubTextByServiceID[service.ID] {
			subText = service.SubText
		}
		fmt.Printf("Service: %s\nSubText: %s\n\n", service.Name, subText)
	}
}

type ListServiceIDsConfig struct {
	ID    string  `json:"id"`
	Value []int64 `json:"value"`
}

type Service struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	SubText string `json:"subText"`
}

func getActivatedServiceIDs() ([]int64, error) {
	resp, err := http.Get("http://my-json-server.typicode.com/havinhdat/restful-db/config/activatedServiceIDs")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()
	configVal := &ListServiceIDsConfig{}
	if err = json.Unmarshal(body, configVal); err != nil {
		return nil, err
	}

	return configVal.Value, nil
}

func getServiceByID(id int64) (*Service, error) {
	resp, err := http.Get(fmt.Sprintf("http://my-json-server.typicode.com/havinhdat/restful-db/services/%d", id))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	service := &Service{}
	if err := json.Unmarshal(body, service); err != nil {
		return nil, err
	}

	return service, nil
}

func getEnabledSubTextServiceIDs() ([]int64, error) {
	resp, err := http.Get("http://my-json-server.typicode.com/havinhdat/restful-db/config/allowedSubTextServiceIDs")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()
	configVal := &ListServiceIDsConfig{}
	if err = json.Unmarshal(body, configVal); err != nil {
		return nil, err
	}

	return configVal.Value, nil
}
