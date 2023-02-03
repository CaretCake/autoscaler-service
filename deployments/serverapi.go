package deployments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "https://8742bce6-176e-43f2-ba49-46349aea041a.mock.pstmn.io"

func DiscoverDeployments() []byte {
	resp, err := http.Get(url + "/discover")
	if err != nil {
		fmt.Printf("error: request failed: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: could not read response body: %s\n", err)
	}

	return body
}

func DeploymentStatus(id string) []byte {
	resp, err := http.Get(url + "/status/" + id)
	if err != nil {
		fmt.Printf("error: request failed: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: could not read response body: %s\n", err)
	}

	return body
}

type ScalingInfo struct {
	DeploymentId string `json:"deployment_id"`
	Delta        int    `json:"delta"`
}

func ScaleDeployment(delta int, id string) []byte {
	srb := ScalingInfo{
		DeploymentId: id,
		Delta:        delta,
	}
	jsonData, err := json.Marshal(srb)
	if err != nil {
		fmt.Printf("error: could not marshal: %s\n", err)
	}

	resp, err := http.Post(url+"/scale", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("error: request failed: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: could not read response body: %s\n", err)
	}

	return body
}
