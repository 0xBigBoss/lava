package relayer

import (
	context "context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/lavanet/lava/relayer/chainproxy"
	"github.com/lavanet/lava/relayer/sentry"
	"github.com/lavanet/lava/relayer/sigs"
	"github.com/lavanet/lava/relayer/testclients"
	"github.com/lavanet/lava/utils"
	"github.com/spf13/pflag"
)

func TestClient(
	ctx context.Context,
	clientCtx client.Context,
	chainID string,
	apiInterface string,
	flagSet *pflag.FlagSet,
) {
	// Every client must preseed
	rand.Seed(time.Now().UnixNano())
	var testErrors error = nil
	if chainID == "ALL" {
		testErrors = TestAllClients(ctx, clientCtx, chainID, flagSet)
	} else {
		testErrors = TestSingleClient(ctx, clientCtx, chainID, apiInterface, flagSet)
	}

	if testErrors != nil {
		log.Fatalln(fmt.Sprintf("%s Client test failed with errors %s", chainID, testErrors))
	} else {
		log.Printf("%s Client test  complete \n", chainID)
	}
}

func TestSingleClient(
	ctx context.Context,
	clientCtx client.Context,
	chainID string,
	apiInterface string,
	flagSet *pflag.FlagSet,
) error {
	//
	sk, _, err := utils.GetOrCreateVRFKey(clientCtx)
	if err != nil {
		log.Fatalln("error: GetOrCreateVRFKey", err)
	}
	// Start sentry
	sentry := sentry.NewSentry(clientCtx, chainID, true, nil, nil, apiInterface, sk, flagSet, 0)
	err = sentry.Init(ctx)
	if err != nil {
		log.Fatalln("error sentry.Init", err)
	}
	go sentry.Start(ctx)
	for sentry.GetBlockHeight() == 0 {
		time.Sleep(1 * time.Second)
	}

	//
	// Node
	chainProxy, err := chainproxy.GetChainProxy("", 1, sentry)
	if err != nil {
		log.Fatalln("error: GetChainProxy", err)
	}

	//
	// Set up a connection to the server.
	log.Println("TestClient connecting")

	keyName, err := sigs.GetKeyName(clientCtx)
	if err != nil {
		log.Fatalln("error: getKeyName", err)
	}

	privKey, err := sigs.GetPrivKey(clientCtx, keyName)
	if err != nil {
		log.Fatalln("error: getPrivKey", err)
	}
	clientKey, _ := clientCtx.Keyring.Key(keyName)
	log.Println("Client pubkey", clientKey.GetPubKey().Address())

	//
	// Run tests
	var testErrors error = nil
	switch chainID {
	case "ETH1", "ETH4", "GTH1":
		testErrors = testclients.EthTests(ctx, chainProxy, privKey)
	case "COS1":
		testErrors = testclients.TerraTests(ctx, chainProxy, privKey, apiInterface)
	case "COS3":
		testErrors = testclients.OsmosisTests(ctx, chainProxy, privKey, apiInterface)
	case "LAV1":
		testErrors = testclients.LavaTests(ctx, chainProxy, privKey, apiInterface, sentry, clientCtx)
	}

	return testErrors
}

func TestAllClients(
	ctx context.Context,
	clientCtx client.Context,
	chainID string,
	flagSet *pflag.FlagSet,
) error {
	testclients.PrintStatusNoticable("Testing All Clients")

	all_chain_ids := []string{"ETH1", "ETH4", "GTH1", "COS1", "COS3", "LAV1"}
	all_errors := []string{}

	sk, _, err := utils.GetOrCreateVRFKey(clientCtx)
	if err != nil {
		log.Fatalln("error: GetOrCreateVRFKey", err)
	}
	keyName, err := sigs.GetKeyName(clientCtx)
	if err != nil {
		log.Fatalln("error: getKeyName", err)
	}

	privKey, err := sigs.GetPrivKey(clientCtx, keyName)
	if err != nil {
		log.Fatalln("error: getPrivKey", err)
	}
	clientKey, _ := clientCtx.Keyring.Key(keyName)
	log.Println("Client pubkey", clientKey.GetPubKey().Address())

	for idx, id := range all_chain_ids {
		// Start sentry
		var apiInterfaceTests string
		switch id {
		case "ETH1", "ETH4", "GTH1":
			apiInterfaceTests = "jsonrpc" // testing jsonrpc for eth.
		default:
			apiInterfaceTests = "rest" // currently testing only rest for osmosis and lava
		}

		sentry := sentry.NewSentry(clientCtx, id, true, nil, nil, apiInterfaceTests, sk, flagSet, 0)
		if idx == 0 {
			err = sentry.Init(ctx)
			if err != nil {
				log.Fatalln("error sentry.Init", err)
			}
		} else {
			err = sentry.Reset(ctx)
			if err != nil {
				log.Fatalln("error sentry.Init", err)
			}
			err = sentry.Init(ctx)
		}
		go sentry.Start(ctx)
		time.Sleep(1 * time.Second)

		// Node
		chainProxy, err := chainproxy.GetChainProxy("", 1, sentry)
		if err != nil {
			log.Fatalln("error: GetChainProxy", err)
		}

		switch id {
		case "ETH1", "ETH4", "GTH1":
			all_errors = append(all_errors, testclients.EthTests(ctx, chainProxy, privKey).Error())
		case "COS1":
			all_errors = append(all_errors, testclients.TerraTests(ctx, chainProxy, privKey, apiInterfaceTests).Error())
		case "COS3":
			all_errors = append(all_errors, testclients.OsmosisTests(ctx, chainProxy, privKey, apiInterfaceTests).Error())
		case "LAV1":
			all_errors = append(all_errors, testclients.LavaTests(ctx, chainProxy, privKey, apiInterfaceTests, sentry, clientCtx).Error())
		}

		testclients.PrintStatusNoticable("Test finished, moving to next test")
	}

	if len(all_errors) > 0 {
		return fmt.Errorf(strings.Join(all_errors, ",\n"))
	}

	return nil
}
