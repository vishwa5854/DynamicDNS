package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type DomainRecord struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

type DomainsApiResponse struct {
	DomainRecords []DomainRecord `json:"domain_records"`
}

var (
	domain              = "YOUR_DOMAIN_HERE"
	token               = "YOUR_DIGITAL_OCEAN_TOKEN_HERE"
	digitalOceanBaseUrl = "https://api.digitalocean.com/v2/domains/" + domain + "/records/"
)

var headers map[string][]string = map[string][]string{
	"Content-Type":  {"application/json"},
	"Authorization": {"Bearer " + token},
}

func getCurrentIpUsingIpify() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.ipify.org", nil)
	if err != nil {
		return "", fmt.Errorf("error while trying to get current ip address using ipify : %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while trying to get current ip address using ipify : %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error while trying to get current ip address using ipify : %v", err)
	}
	return string(body), nil
}

func updateDomainRecord(domainRecord DomainRecord) error {
	payload, err := json.Marshal(domainRecord)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%v%v", digitalOceanBaseUrl, domainRecord.Id), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header = headers
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	response := DomainRecord{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error while updating %v.%v with new ip address %v : %v", domainRecord.Name, domain, domainRecord.Data, response)
	} else {
		fmt.Printf("Updated %v.%v with latest IP Address %v\n", domainRecord.Name, domain, domainRecord.Data)
	}
	return nil
}

func main() {
	waitGroup := sync.WaitGroup{}
	currentIp, err := getCurrentIpUsingIpify()
	ipChanged := false
	fmt.Println("Current IP Address is : ", currentIp)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := "https://api.digitalocean.com/v2/domains/" + domain + "/records?type=A"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header = headers

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	response := DomainsApiResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		fmt.Printf("failed to decode json : %v", err)
		return
	}

	for _, i := range response.DomainRecords {
		if i.Data != currentIp {
			ipChanged = true
			i.Data = currentIp
			waitGroup.Add(1)

			go func(domainRecord DomainRecord) {
				defer waitGroup.Done()

				err := updateDomainRecord(domainRecord)
				if err != nil {
					fmt.Println(err)
				}
			}(i)
		}
	}
	waitGroup.Wait()

	if !ipChanged {
		fmt.Println("No domain records are found with outdated IP Address")
	}
}
