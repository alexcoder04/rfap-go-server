package main

import "runtime"

const (
	CONN_RECV_TIMEOUT_SECS = 65 // disconnect client if it sleeps for so long
	MAX_THREADS_WAIT_SECS  = 5  // don't accept new connections for so long if max number reached
	CONT_LEN_INDIC_LENGTH  = 4  // length of the content length indocator in bytes

	connHost = "localhost"
	connPort = "6700"
	connType = "tcp"
)

var SUPPORTED_RFAP_VERSIONS = []uint32{1}

var MAX_CLIENTS = runtime.NumCPU() * 4 // 4 clients per core
