package settings

const (
	RFAP_VERSION = 0x0003

	CONN_RECV_TIMEOUT_SECS = 65 // disconnect client if it sleeps for so long
	MAX_THREADS_WAIT_SECS  = 5  // don't accept new connections for so long if max number reached
	MAX_BYTES_SEND_AT_ONCE = 1024 * 16

	VERSION_LENGTH        = 2
	CONT_LEN_INDIC_LENGTH = 4 // length of the content length indicator in bytes
	COMMAND_LENGTH        = 4
	CHECKSUM_LENGTH       = 32
)

var SUPPORTED_RFAP_VERSIONS = []uint32{3}
