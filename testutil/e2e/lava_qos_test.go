package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lavanet/lava/testutil/e2e/proxy"
)

func TestLavaQoS(t *testing.T) {
	timeout := time.After(11 * time.Minute)
	done := make(chan bool)
	go func() {
		// run lava testing
		LavaTestQoSflow(t)
		// Test finished on time !
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}

func LavaTestQoSflow(t *testing.T) {
	finalresults, err := QoSFlowTest(t)
	wrapGoTest(t, finalresults, err)
}

func QoSFlowTest(t *testing.T) ([]*TestResult, error) {
	readEnvVars()
	prepTest(t)

	// Test Configs
	resetGenesis := true
	init_chain := true
	run_providers_osmosis := true
	run_client_osmosis := true

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
			1500,
			100,
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
			1500,
			100,
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
			await(node, "relay payment 1/3 osmosis", found_relay_payment, []interface{}{"Latency: ", 0.5, 0.8, "Sync: ", 1.0, 1.0, "Availability: ", 1.0, 1.0}, "awaiting for OSMOSIS payment to proceed... ", true)
			fmt.Println(" ::: GOT OSMOSIS PAYMENT !!!")
			silent(clientOsmoRPC)
			clientOsmoRest := TestProcess("clientOsmoRest", "lavad test_client COS3 rest --from user2", clientTest)
			await(node, "relay payment 2/3 osmosis", found_relay_payment, []interface{}{"Latency: ", 0.5, 0.8, "Sync: ", 1.0, 1.0, "Availability: ", 1.0, 1.0}, "awaiting for OSMOSIS payment to proceed... ", true)
			fmt.Println(" ::: GOT OSMOSIS PAYMENT !!!")
			silent(clientOsmoRest)
			silent(prov_osm)
		}
		srv1.Shutdown(context.Background())
		srv2.Shutdown(context.Background())
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
