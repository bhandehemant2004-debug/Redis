package protocol

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/obj"
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
func SimpleError(error string)string{
	return "-"+error+"\r\n"
}

func ResponseXRange(entries []*obj.Entry) []byte {
	result := make([]byte, 0)


	result = append(result, '*')
	result = strconv.AppendInt(result, int64(len(entries)), 10)
	result = append(result, '\r', '\n')

	for _, entry := range entries {

		// Each entry is:
		// [ID, [field, value, field, value, ...]]
		result = append(result, '*', '2', '\r', '\n')

		// ID
		result = append(result, '$')
		result = strconv.AppendInt(result, int64(len(entry.Id)), 10)
		result = append(result, '\r', '\n')
		result = append(result, entry.Id...)
		result = append(result, '\r', '\n')

		// Field/value array
		result = append(result, '*')
		result = strconv.AppendInt(result, int64(len(entry.Entries)*2), 10)
		result = append(result, '\r', '\n')

		for field, value := range entry.Entries {

			// field
			result = append(result, '$')
			result = strconv.AppendInt(result, int64(len(field)), 10)
			result = append(result, '\r', '\n')
			result = append(result, field...)
			result = append(result, '\r', '\n')

			// value
			result = append(result, '$')
			result = strconv.AppendInt(result, int64(len(value)), 10)
			result = append(result, '\r', '\n')
			result = append(result, value...)
			result = append(result, '\r', '\n')
		}
	}

	return result
}
func ResponseXRead(keys []string, entries [][]*obj.Entry) []byte {
	result := make([]byte, 0)

	// Outer array = number of streams
	result = append(result, '*')
	result = strconv.AppendInt(result, int64(len(keys)), 10)
	result = append(result, '\r', '\n')

	for i, key := range keys {

		// [key, entries]
		result = append(result, '*', '2', '\r', '\n')

		// key
		result = append(result, '$')
		result = strconv.AppendInt(result, int64(len(key)), 10)
		result = append(result, '\r', '\n')
		result = append(result, key...)
		result = append(result, '\r', '\n')

		streamEntries := entries[i]

		// entries array
		result = append(result, '*')
		result = strconv.AppendInt(result, int64(len(streamEntries)), 10)
		result = append(result, '\r', '\n')

		for _, entry := range streamEntries {

			// [id, fields]
			result = append(result, '*', '2', '\r', '\n')

			// ID
			result = append(result, '$')
			result = strconv.AppendInt(result, int64(len(entry.Id)), 10)
			result = append(result, '\r', '\n')
			result = append(result, entry.Id...)
			result = append(result, '\r', '\n')

			// fields/values
			result = append(result, '*')
			result = strconv.AppendInt(
				result,
				int64(len(entry.Entries)*2),
				10,
			)
			result = append(result, '\r', '\n')

			for field, value := range entry.Entries {

				// field
				result = append(result, '$')
				result = strconv.AppendInt(result, int64(len(field)), 10)
				result = append(result, '\r', '\n')
				result = append(result, field...)
				result = append(result, '\r', '\n')

				// value
				result = append(result, '$')
				result = strconv.AppendInt(result, int64(len(value)), 10)
				result = append(result, '\r', '\n')
				result = append(result, value...)
				result = append(result, '\r', '\n')
			}
		}
	}

	return result
}