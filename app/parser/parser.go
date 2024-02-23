package parser

import (
	"strconv"
	"strings"
)

type RedisProtocolParser struct {
	Idx int
}

func BulkStringEncode(str string) []byte {
	return []byte("$" + strconv.Itoa(len(str)) + "\r\n" + str + "\r\n")
}

func (rpp *RedisProtocolParser) BulkStringDecode(pieces []string) (interface{}, error) {
	size, err := strconv.Atoi(pieces[rpp.Idx][1:])
	if err != nil {
		return nil, err
	}
	if size < 0 {
		return nil, nil
	}
	rpp.Idx += 1
	return pieces[rpp.Idx][:size], nil
}

func ArrayEncode(array []string) []byte {
	str := "*" + strconv.Itoa(len(array)) + "\r\n"
	for _, v := range array {
		str += string(BulkStringEncode(v))
	}
	return []byte(str)
}

func (rpp *RedisProtocolParser) arrayParser(pieces []string) (interface{}, error) {
	var array []string
	size, err := strconv.Atoi(pieces[rpp.Idx][1:])
	if err != nil {
		return nil, err
	}
	for i := 0; i < size; i++{
		rpp.Idx += 1
		item, err := rpp.defineDataType(pieces)
		if err != nil {
			return item, err
		}
		it, ok := item.(string)
		if !ok {
			return nil, err
		}
		array = append(array, it)
	}
	return array, nil
}

func (rpp *RedisProtocolParser) defineDataType(pieces []string) (interface{}, error) {
	switch {
		case strings.HasPrefix(pieces[rpp.Idx], "$"):
			return rpp.BulkStringDecode(pieces)
		case strings.HasPrefix(pieces[rpp.Idx], "*"):
			return rpp.arrayParser(pieces)
	}
	return nil, nil
}

func (rpp *RedisProtocolParser) ParserProtocol(str string) (interface{}, error) {
	pieces := strings.Split(str, "\r\n")
	pieces = pieces[:len(pieces) -1]
	return rpp.defineDataType(pieces)
}