package manager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/pkg/model"
)

const (
	NAMESPACE       = "NAMESPACE"
	LWD_CONF        = "LWD_CONF"
	ZCASH_CONF      = "ZCASH_CONF"
	ENVOY_CONF      = "ENVOY_CONF"
	DEPLOYMENT      = "DEPLOYMENT"
	SERVICE         = "SERVICE"
	INGRESS_STOPPED = "INGRESS_STOPPED"
	INGRESS         = "INGRESS"
	CREDENTIALS     = "CREDENTIALS"

	ZCASH_SVC_PREFIX = "zcashd-svc"
	USERNAME         = "Username"
	PASSWORD         = "Password"

	ZcashConf          = "ZcashConf"
	InstanceProperties = "InstanceProperties"

	LIGHT_WALLET_IMAGE = "Lightwallet"
	ZCASH              = "Zcash"
	METRICS            = "Metrics"
	GRPC               = "GRPC"
	HTTP               = "HTTP"

	ENVOY_PORT   = "envoy"
	SERVICE_PORT = "service"
	METRICS_PORT = "metrics"
	HTTP_PORT    = "http"

	LWD_IMAGE     = "lwd"
	NODE_IMAGE    = "node"
	METRICS_IMAGE = "metrics"

	DEFAULT_ZCASH_CONF        = "default"
	MINER_ZCASH_CONF          = "miner"
	MAINNET_ZCASH_CONF        = "mainnet"
	TESTNET_ZCASH_CONF        = "testnet"
	RPCPORT_ZCASH_PROPERTY    = "rpcport"
	CONNECT_ZCASH_PROPERTY    = "connect"
	MINER_ZCASH_PROPERTY      = "miner"
	RESOURCE_REQUEST_PROPERTY = "request"

	ZCASH_INSTANCE_NAME = "ZcashInstanceName"
	ZCASH_INSTANCE      = "ZcashInstance"
	ZCASH_PORT          = "ZcashPort"
	LOG_LEVEL           = "LogLevel"

	zcashInstanceProperty = "zcashInstance"
	logLevelProperty      = "logLevel"
)

func addLWDInstance(instance *model.Instance, name string, request model.ResourceRequest) {
	lwdInstances := request.Properties["lwdInstance"].([]string)

	found := false
	for _, inst := range lwdInstances {
		if inst == name {
			found = true
		}
	}
	if !found {
		lwdInstances = append(lwdInstances, name)
	}

	request.Properties["lwdInstance"] = lwdInstances
}

func removeLWDInstance(instance *model.Instance, name string, request model.ResourceRequest) {
	lwdInstances := request.Properties["lwdInstance"].([]string)
	for index, inst := range lwdInstances {
		if inst == name {
			request.Properties["lwdInstance"] = append(lwdInstances[:index], lwdInstances[index+1:]...)
			break
		}
	}
}

func addZcashPeer(instance *model.Instance, name string, request model.ResourceRequest) {
	peers := request.Properties["peers"].([]string)
	found := false
	for _, inst := range peers {
		if inst == name {
			found = true
		}
	}
	if !found {
		peers = append(peers, name)
	}

	request.Properties["peers"] = peers
}

func removeZcashPeer(instance *model.Instance, name string, request model.ResourceRequest) {
	peers := request.Properties["peers"].([]string)
	for index, inst := range peers {
		if inst == name {
			request.Properties["peers"] = append(peers[:index], peers[index+1:]...)
			break
		}
	}
}

func createZcashConf(ic *model.BlockchainNodeInfo, miner bool, network model.NetworkType, rpcport string) []model.KVPair {

	var zcashConf = make([]model.KVPair, 0)

	zcashConf = append(zcashConf, ic.Settings[DEFAULT_ZCASH_CONF]...)

	if miner {
		zcashConf = append(zcashConf, ic.Settings[MINER_ZCASH_CONF]...)
	}

	if network == model.NetworkTypeMain {
		zcashConf = append(zcashConf, ic.Settings[MAINNET_ZCASH_CONF]...)
	} else if network == model.NetworkTypeTest {
		zcashConf = append(zcashConf, ic.Settings[TESTNET_ZCASH_CONF]...)
	}

	zcashConf = append(zcashConf, model.KVPair{Key: RPCPORT_ZCASH_PROPERTY, Value: rpcport})
	return zcashConf
}

func getZcashPeers(conf []model.KVPair, rpcport string, namespace string, peers ...model.Instance) []model.KVPair {
	connect := make([]string, 0)
	// peerProperty := ""
	if peers != nil {
		for _, peer := range peers {
			connect = append(connect, peer.Name)
			conf = append(conf, model.KVPair{Key: CONNECT_ZCASH_PROPERTY, Value: getZcashInstanceHost(peer.Name, namespace) + ":" + rpcport})
			//if index > 0 {
			//	peerProperty += ","
			//}
			//peerProperty += peer.Name
		}
	}

	return conf
}

func getZcashInstanceHost(name, namespace string) string {
	return fmt.Sprintf("%s-%s.%s.svc.cluster.local", ZCASH_SVC_PREFIX, name, namespace)
}

func getZcashInstancePort(ctx context.Context) string {
	zcashIc, err := helper.GetBlockchainNodeInfo(ctx, model.InstanceTypeZCASH)
	if err != nil {
		return ""
	}

	return strconv.FormatInt(int64(zcashIc.GetPort(SERVICE_PORT)), 10)
}
