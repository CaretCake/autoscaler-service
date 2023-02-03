package deployments

type Config struct {
	ServersPerHost int    `json:"servers_per_host"`
	TargetFreePct  int    `json:"target_free_pct"`
	Id             string `json:"id"`
}
