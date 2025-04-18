package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration(t *testing.T) {
	net, err := network.New(t.Context())
	if err != nil {
		t.Logf("failed to create network: %s", err)

		return
	}
	defer func() {
		if err := net.Remove(t.Context()); err != nil { //nolint:govet // ignore shadowing
			t.Logf("failed to remove network: %s", err)
		}
	}()

	DBreq := testcontainers.ContainerRequest{
		Name:         "memgraph",
		Hostname:     "memgraph",
		ExposedPorts: []string{"7687/tcp", "7444/tcp"},
		Image:        "memgraph/memgraph-mage:latest",
		Cmd: []string{
			"--log-level=TRACE",
			"--storage-mode=ON_DISK_TRANSACTIONAL",
			"--storage-snapshot-interval-sec=86400",
			"--storage-snapshot-retention-count=60",
			"--storage-property-store-compression-enabled=true",
			"--storage-property-store-compression-level=mid",
			"--storage-snapshot-on-exit=true",
		},
		Env: map[string]string{
			"MEMGRAPH_PASSWORD": "memgraph",
			"MEMGRAPH_USER":     "memgraph",
		},
		Networks:   []string{net.Name},
		WaitingFor: wait.ForListeningPort("7687/tcp"),
	}

	memgraphC, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: DBreq,
		Started:          true,
	})
	defer testcontainers.CleanupContainer(t, memgraphC)
	if err != nil {
		t.Log(memgraphC.Logs(t.Context()))
	}
	require.NoError(t, err)

	memgraphHost, err := memgraphC.ContainerIP(t.Context())
	require.NoError(t, err)

	memgraphURI := "bolt://" + memgraphHost + ":7687" // + memgraphPort.Port()

	DBAdapterReq := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "./",
			Dockerfile: "integration-test.dockerfile",
		},
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"MEMGRAPH_URI": memgraphURI,
			"HTTP_PORT":    ":8080",
		},
		Networks:   []string{net.Name},
		WaitingFor: wait.ForLog("Starting server"),
	}

	dbAdapterC, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: DBAdapterReq,
		Started:          true,
	})
	defer testcontainers.CleanupContainer(t, dbAdapterC)
	if err != nil {
		aliasis, naerr := memgraphC.NetworkAliases(t.Context())
		if naerr != nil {
			t.Logf("failed to get network aliases: %s", naerr)
		} else {
			t.Log("Memgraph Network Aliases: ", aliasis)
		}
		t.Log("Memgraph Running: ", memgraphC.IsRunning())
		t.Log(dbAdapterC.Logs(t.Context()))
	}

	require.NoError(t, err)
}
