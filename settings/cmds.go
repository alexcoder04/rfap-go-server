package settings

const (
	CMD_PING       = 0x00000000
	CMD_DISCONNECT = 0x01000000
	CMD_INFO       = 0xa0010000
	CMD_ERROR      = 0xffffffff

	CMD_FILE_READ   = 0xf0020000
	CMD_FILE_DELETE = 0xf1010000
	CMD_FILE_CREATE = 0xf1020000
	CMD_FILE_COPY   = 0xf1030000
	CMD_FILE_MOVE   = 0xf1040000
	CMD_FILE_WRITE  = 0xf2010000

	CMD_DIRECTORY_READ   = 0xd0020000
	CMD_DIRECTORY_DELETE = 0xd1010000
	CMD_DIRECTORY_CREATE = 0xd1020000
	CMD_DIRECTORY_COPY   = 0xd1030000
	CMD_DIRECTORY_MOVE   = 0xd1040000

	ERROR_OK                  = 0
	ERROR_FILE_NOT_EXISTS     = 1
	ERROR_UNKNOWN             = 2
	ERROR_INVALID_COMMAND     = 3
	ERROR_INVALID_FILE_TYPE   = 4
	ERROR_ACCESS_DENIED       = 5
	ERROR_INVALID_CONTENT_LEN = 6
	ERROR_FILE_EXISTS         = 7
	ERROR_WHILE_STAT          = 8
)
