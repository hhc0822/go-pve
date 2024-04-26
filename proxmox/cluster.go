package proxmox

import (
	"fmt"
)

type Cluster struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Nodes   int64  `json:"nodes"`
	Quorate int64  `json:"quorate"`
	Type    string `json:"type"`
	Version int64  `json:"version"`
	Ip      string `json:"ip"`
	Level   string `json:"level"`
	Local   int64  `json:"local"`
	NodeId  int64  `json:"nodeid"`
	Online  int64  `json:"online"`
}

func (c *Client) GetClusterStatus(response interface{}) error {
	url := fmt.Sprintf("/cluster/status")
	resp, err := c.session.Get(url, nil, nil)
	if err != nil {
		return err
	}
	return handleResponse(resp, response)
}
