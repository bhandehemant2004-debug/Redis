package protocol

import (
	"fmt"
	"strconv"
)


func Parser(pending []byte) ([]string, []byte, bool) {
	command := []string{}

	if len(pending) == 0 {
		return command, pending, false
	}
	if pending[0] != '*' {
		return command, pending, false
	}
	end := findCRLF(pending, 1)
	if end == -1 {
		return command, pending, false
	}
	count, err := strconv.Atoi(string(pending[1:end]))
	if err != nil {
		return command, pending, false
	}
	pos := end + 2
	for i := 0; i < count; i++ {
		str, nextPos, ok := parseNextBulkString(pending, pos)
		if !ok {
			return []string{}, pending, false
		}
		command = append(command, str)
		pos = nextPos
	}
	remaining := pending[pos:]
	return command, remaining, true
}
func parseNextBulkString(data []byte, pos int) (string, int, bool) {
	if pos >= len(data) {
		return "", pos, false
	}
	if data[pos] != '$' {
		return "", pos, false
	}
	end := findCRLF(data, pos+1)
	if end == -1 {
		return "", pos, false
	}
	length, err := strconv.Atoi(string(data[pos+1:end]))
	if err != nil {
		return "", pos, false
	}
	dataStart := end + 2
	if dataStart+length+2 > len(data) {
		return "", pos, false
	}
	value := string(data[dataStart : dataStart+length])

	nextPos := dataStart + length + 2

	return value, nextPos, true
}
func findCRLF(data []byte, start int) int {
    for i := start; i+1 < len(data); i++ {
        if data[i] == '\r' && data[i+1] == '\n' {
            return i
        }
    }

    return -1
}


func Response_bulk_string(s string)[]byte{

	var result = make([]byte,0)

	result = append(result, '$')
	result = strconv.AppendInt(result, int64(len(s)), 10)
	result = append(result, '\r', '\n')
	result = append(result, s...)
	result = append(result, '\r', '\n')

	return result
}
func AppendInteger(value int) []byte {
    return []byte(fmt.Sprintf(":%d\r\n", value))
}

func ResponseArray(values []string) string {
    result := fmt.Sprintf("*%d\r\n", len(values))

    for _, value := range values {
        result += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
    }

    return result
}
