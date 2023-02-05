package main

import (
	"fmt"
	"testing"
)

func TestCalculateDelta(t *testing.T) {
	var tests = []struct {
		c      DeploymentConfig
		s      Status
		expect int
	}{
		{
			c: DeploymentConfig{
				ServersPerHost: 10,
				TargetFreePct:  15,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 10,
				TotalServers: 100,
				FreeServers:  10,
			},
			expect: 1,
		},
		{
			c: DeploymentConfig{
				ServersPerHost: 10,
				TargetFreePct:  15,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 100,
				TotalServers: 1000,
				FreeServers:  100,
			},
			expect: 6,
		},
		{
			c: DeploymentConfig{
				ServersPerHost: 3,
				TargetFreePct:  30,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 200,
				TotalServers: 600,
				FreeServers:  400,
			},
			expect: -104,
		},
		{
			c: DeploymentConfig{
				ServersPerHost: 250,
				TargetFreePct:  10,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 3,
				TotalServers: 750,
				FreeServers:  113,
			},
			expect: 0,
		},
		{
			c: DeploymentConfig{
				ServersPerHost: 1,
				TargetFreePct:  90,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 600,
				TotalServers: 600,
				FreeServers:  541,
			},
			expect: -10,
		},
		{
			c: DeploymentConfig{
				ServersPerHost: 1,
				TargetFreePct:  90,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 600,
				TotalServers: 600,
				FreeServers:  540,
			},
			expect: 0,
		},
		{
			c: DeploymentConfig{
				ServersPerHost: 1,
				TargetFreePct:  90,
				Id:             "",
			},
			s: Status{
				CurrentHosts: 600,
				TotalServers: 600,
				FreeServers:  540,
			},
			expect: 0,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v,%v", tt.c, tt.s)
		t.Run(testname, func(t *testing.T) {
			got := calculateDelta(tt.c, tt.s)
			if got != tt.expect {
				t.Errorf("got: %d, expected: %d", got, tt.expect)
			}
		})
	}
}

// Test all 0 values, alternating 0 values
// Test empty ids
