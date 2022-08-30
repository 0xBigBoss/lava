// file_pipe.go
package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lavanet/lava/testutil/e2e/proxy"
)

var nodeTest = TestProc{
	filter:           []string{"STARPORT]", "!", "lava_", "ERR_", "panic"},
	expectedEvents:   []string{"🌍", "lava_spec_add", "lava_provider_stake_new", "lava_client_stake_new", "lava_relay_payment"},
	unexpectedEvents: []string{"exit status", "cannot build app", "connection refused", "ERR_client_entries_pairing", "ERR"},
	tests:            events(),
	strict:           false}
var initTest = TestProc{
	filter:           []string{":::", "raw_log", "Error", "error", "panic"},
	expectedEvents:   []string{"init done"},
	unexpectedEvents: []string{"Error"},
	tests:            events(),
	strict:           true}
var providersTest = TestProc{
	filter:           []string{"sent (new/from cache)", "Server", "updated", "server", "error"},
	expectedEvents:   []string{"listening"},
	unexpectedEvents: []string{"ERROR", "refused", "Missing Payment"},
	tests:            events(),
	strict:           true}
var clientTest = TestProc{
	filter:         []string{":::", "reply", "no pairings available", "update", "connect", "rpc", "pubkey", "signal", "Error", "error", "panic"},
	expectedEvents: []string{"update pairing list!", "Client pubkey"},
	// unexpectedEvents: []string{"no pairings available", "error", "Error", "signal: interrupt"},
	unexpectedEvents: []string{"no pairings available", "Error", "signal: interrupt"},
	tests:            events(),
	strict:           true}

func FullFlowTest(t *testing.T) ([]*TestResult, error) {
	readEnvVars()
	prepTest(t)

	// Test Configs
	resetGenesis := true
	init_chain := true
	run_providers_osmosis := true
	run_providers_eth := true
	run_client_osmosis := true
	run_client_eth := true

	start_lava := "killall ignite; killall lavad; cd " + homepath + " && ignite chain serve -v -r  "
	if !resetGenesis {
		start_lava = "lavad start "
	}

	// Start Test Processes
	node := TestProcess("ignite", start_lava, nodeTest)
	await(node, "lava node is running", lava_up, nil, "awaiting for node to proceed...", true)

	if init_chain {
		sleep(2)
		init := TestProcess("init", homepath+"scripts/init.sh", initTest)
		await(init, "get init done", init_done, nil, "awaiting for init to proceed...", true)
	}

	if run_providers_osmosis {
		MOCK_PORT_REST := int64(2031)
		rpcProxyProcessRest := proxy.NewProxy(
			"osmosis_rest",
			MOCK_PORT_REST,
			os.Getenv("OSMO_HOST"),
			proxy.GetMockFilePath("osmosis_rest", ""),
			false,
			true,
			true,
			false,
			0,
			1000,
		)
		srv1 := proxy.StartProxy(rpcProxyProcessRest) // start

		MOCK_PORT_TM := int64(2041)
		rpcProxyProcessRpc := proxy.NewProxy(
			"osmosis_rpc",
			MOCK_PORT_TM,
			os.Getenv("OSMO_HOST"),
			proxy.GetMockFilePath("osmosis_rpc", ""),
			false,
			true,
			true,
			false,
			0,
			1000,
		)
		srv2 := proxy.StartProxy(rpcProxyProcessRpc) // start

		fmt.Println(" ::: Starting Providers Processes [OSMOSIS] ::: ")
		prov_osm := TestProcess("providers_osmosis", homepath+"scripts/osmosis.sh", providersTest)
		// debugOn(prov_osm)
		println(" ::: Providers Processes Started ::: ")
		await(prov_osm, "Osmosis providers ready", providers_ready, nil, "awaiting for providers to listen to proceed...", true)

		if run_client_osmosis {
			sleep(1)
			fmt.Println(" ::: Starting Client Process [OSMOSIS] ::: ")
			clientOsmoRPC := TestProcess("clientOsmoRPC", "lavad test_client COS3 tendermintrpc --from user2", clientTest)
			await(node, "relay payment 1/3 osmosis", found_relay_payment, []interface{}{"Latency: ", 1.0, 1.0, "Sync: ", 1.0, 1.0, "Availability: ", 1.0, 1.0}, "awaiting for OSMOSIS payment to proceed... ", true)
			fmt.Println(" ::: GOT OSMOSIS PAYMENT !!!")
			silent(clientOsmoRPC)
			clientOsmoRest := TestProcess("clientOsmoRest", "lavad test_client COS3 rest --from user2", clientTest)
			await(node, "relay payment 2/3 osmosis", found_relay_payment, []interface{}{"Latency: ", 1.0, 1.0, "Sync: ", 1.0, 1.0, "Availability: ", 1.0, 1.0}, "awaiting for OSMOSIS payment to proceed... ", true)
			fmt.Println(" ::: GOT OSMOSIS PAYMENT !!!")
			silent(clientOsmoRest)
			silent(prov_osm)
		}
		srv1.Shutdown(context.Background())
		srv2.Shutdown(context.Background())
	}
	if run_providers_eth {
		fmt.Println(" ::: Starting Providers Processes [ETH] ::: ")
		prov_eth := TestProcess("providers_eth", homepath+"scripts/eth.sh", providersTest)
		fmt.Println(" ::: Providers Processes Started ::: ")
		await(prov_eth, "ETH providers ready", providers_ready_eth, nil, "awaiting for providers to listen to proceed...", true)

		if run_client_eth {
			sleep(1)
			fmt.Println(" ::: Starting Client Process [ETH] ::: ")
			clientEth := TestProcess("clientEth", "lavad test_client ETH1 jsonrpc --from user1", clientTest)
			await(clientEth, "reply rpc", found_rpc_reply, nil, "awaiting for rpc reply to proceed...", true)
			await(node, "relay payment 3/3 eth", found_relay_payment, []interface{}{"Latency: ", 1.0, 1.0}, "awaiting for ETH payment to proceed...", true)
			fmt.Println(" ::: GOT ETH PAYMENT !!!")
			silent(clientEth)
			silent(prov_eth)
		}
	}

	// FINISHED TEST PROCESSESS
	println("::::::::::::::::::::::::::::::::::::::::::::::")
	awaitErrorsTimeout := 10
	fmt.Println(" ::: wait ", awaitErrorsTimeout, " seconds for potential errors...")
	sleep(awaitErrorsTimeout)
	fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::")
	fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::")
	fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::")

	// Finalize & Display Results
	final := finalizeResults(t)

	return final, nil
}
