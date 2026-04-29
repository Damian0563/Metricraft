package external

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCountryOrigin(ip string) (string, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data CountryOrigin
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Country, nil
}
