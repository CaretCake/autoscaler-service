package deployments

type Status struct {
	CurrentHosts int `json:"current_hosts"`
	TotalServers int `json:"total_servers"`
	FreeServers  int `json:"free_servers"`
}
