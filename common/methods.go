package common

type METHOD = string

const (
	ORDER              METHOD = "ORD"
	ORDER_RESPONSE     METHOD = "R/ORD"
	PROXY              METHOD = "PRX"
	PROXY_RESPONSE     METHOD = "R/PRX"
	OPTIONS            METHOD = "OPT"
	OPTIONS_RESPONSE   METHOD = "R/OPT"
	HEARTBEAT          METHOD = "HBT"
	HEARTBEAT_RESPONSE METHOD = "R/HBT"
)
