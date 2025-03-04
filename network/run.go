package network

import (
	"CTng/CA"
	"CTng/Logger"
	"CTng/gossip"
	"CTng/monitor"
	"CTng/util"
	"fmt"
	"time"
)

func StartCA(CID string) {
	path_prefix := "network/ca_testconfig/" + CID
	path_1 := path_prefix + "/CA_public_config.json"
	path_2 := path_prefix + "/CA_private_config.json"
	path_3 := path_prefix + "/CA_crypto_config.json"
	path_4 := "network/ca_testdata/" + CID + "/ca_testdata.json"
	ctx_ca := CA.InitializeCAContext(path_1, path_2, path_3)
	ctx_ca.StoragePath = path_4
	CA.StartCA(ctx_ca)
}

func StartLogger(LID string) {
	path_prefix := "network/logger_testconfig/" + LID
	path_1 := path_prefix + "/Logger_public_config.json"
	path_2 := path_prefix + "/Logger_private_config.json"
	path_3 := path_prefix + "/Logger_crypto_config.json"
	ctx_logger := Logger.InitializeLoggerContext(path_1, path_2, path_3)
	Logger.StartLogger(ctx_logger)
}

func StartMonitor(MID string) {
	path_prefix := "network/monitor_testconfig/" + MID
	path_1 := path_prefix + "/Monitor_public_config.json"
	path_2 := path_prefix + "/Monitor_private_config.json"
	path_3 := path_prefix + "/Monitor_crypto_config.json"
	ctx_monitor := monitor.InitializeMonitorContext(path_1, path_2, path_3, MID)
	// clean up the storage
	ctx_monitor.InitializeMonitorStorage("network/monitor_testdata/")
	// delete all the files in the storage
	ctx_monitor.CleanUpMonitorStorage()
	//ctx_monitor.Mode = 0
	//wait for 60 seconds
	fmt.Println("Delay 60 seconds to start monitor server")
	time.Sleep(60 * time.Second)
	monitor.StartMonitorServer(ctx_monitor)
}

func StartGossiper(GID string) {
	path_prefix := "network/gossiper_testconfig/" + GID
	path_1 := path_prefix + "/Gossiper_public_config.json"
	path_2 := path_prefix + "/Gossiper_private_config.json"
	path_3 := path_prefix + "/Gossiper_crypto_config.json"
	ctx_gossiper := gossip.InitializeGossiperContext(path_1, path_2, path_3, GID)
	ctx_gossiper.StorageDirectory = "network/gossiper_testdata/" + ctx_gossiper.StorageID + "/"
	ctx_gossiper.StorageFile = "gossiper_testdata.json"
	ctx_gossiper.CleanUpGossiperStorage()
	// create the storage directory if not exist
	util.CreateDir(ctx_gossiper.StorageDirectory)
	// create the storage file if not exist
	util.CreateFile(ctx_gossiper.StorageDirectory + ctx_gossiper.StorageFile)
	gossip.StartGossiperServer(ctx_gossiper)
}
