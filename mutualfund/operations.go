package mutualfund

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (m *Handler) GetLatestNavData(symbol string) (MutualFund, error) {
	var resData MutualFund
	reqURL := fmt.Sprintf(MutualFundLatestBaseURL, symbol)
	m.log.Info(fmt.Sprintf("Request URL: %s", reqURL))
	resp, err := http.Get(reqURL)
	if err != nil {
		return resData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resData, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resData, err
	}
	m.log.Info(fmt.Sprintf("Response body: %v", body))

	if err := json.Unmarshal(body, &resData); err != nil {
		return resData, err
	}
	m.log.Info(fmt.Sprintf("Response navData: %v", resData))

	return resData, nil
}
