package utils

func RetError(metadata HeaderMetadata, errorCode int, errorMsg string) HeaderMetadata {
	metadata.ErrorCode = errorCode
	metadata.ErrorMessage = errorMsg
	return metadata
}

func ConcatBytes(arrays ...[]byte) []byte {
	result := make([]byte, 0)
	for _, i := range arrays {
		result = append(result, i...)
	}
	return result
}

func Uint32ArrayContains(array []uint32, element uint32) bool {
	for _, i := range array {
		if i == element {
			return true
		}
	}
	return false
}
