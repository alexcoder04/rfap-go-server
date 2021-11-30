package main

func Concat(arrays ...[]byte) []byte {
	result := make([]byte, 0)
	for _, i := range arrays {
		result = append(result, i...)
	}
	return result
}
