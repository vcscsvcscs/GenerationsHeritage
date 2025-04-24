package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	integration_tests "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/integration-tests"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

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
		ExposedPorts: []string{"5237/tcp"},
		Env: map[string]string{
			"MEMGRAPH_URI": memgraphURI,
			"HTTP_PORT":    ":5237",
		},
		Networks:   []string{net.Name},
		WaitingFor: wait.ForListeningPort("5237/tcp"),
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

	dbAdapterHost, err := dbAdapterC.Host(t.Context())
	require.NoError(t, err)
	dbAdapterPort, err := dbAdapterC.MappedPort(t.Context(), "5237/tcp")
	require.NoError(t, err)
	dbAdapterURI := "http://" + dbAdapterHost + ":" + dbAdapterPort.Port()

	IntegrationTestFlow(dbAdapterURI)(t)
}

func IntegrationTestFlow(dbAdapterURI string) func(t *testing.T) {
	return func(t *testing.T) {
		client := &http.Client{}
		t.Run("CreatePersonByGoogleIdAndGetById", integration_tests.TestPersonGoogle(dbAdapterURI))
		t.Run("CreatePerson", integration_tests.CreatePersonTest(dbAdapterURI, client))
		t.Run("UpdatePerson", integration_tests.UpdatePersonTest(dbAdapterURI, client))
		t.Run("AddInviteCodeToPerson", integration_tests.UpdatePersonWithInviteCodeTest(dbAdapterURI, client))
		t.Run("CreateFamilyTest", integration_tests.CreateAFamilyTest(dbAdapterURI, client))
		t.Run("CreateRelationships", integration_tests.CreateRelationshipsTest(dbAdapterURI, client))
		t.Run("GetRelationship", integration_tests.GetRelationshipTest(dbAdapterURI, client))
		t.Run("SoftDeletePerson", integration_tests.SoftDeletePersonTest(dbAdapterURI, client))
		t.Run("HardDeletePerson", integration_tests.HardDeletePersonTest(dbAdapterURI, client))
		t.Run("GetPersonById", integration_tests.GetPersonById(dbAdapterURI, client))
		t.Run("GetFamilyTreeByIdTest", integration_tests.GetFamilyTreeByIdTest(dbAdapterURI, client))
		t.Run("GetFamilyTreeWithSpousesByIdTest", integration_tests.GetFamilyTreeWithSpousesByIdTest(dbAdapterURI, client))
		t.Run("VerifyRelationships", integration_tests.UpdateRelationshipTest(dbAdapterURI, client))
		t.Run("DeleteRelationship", integration_tests.DeleteRelationshipTest(dbAdapterURI, client))
	}
}
