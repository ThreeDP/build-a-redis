package parser

import (
	"reflect"
	"testing"
)



func TestParserRedisProtocol(t *testing.T) {
	t.Run("Test parser the bulk string $4hello\\r\\n", func(t *testing.T) {
		rrp := RedisProtocolParser{Idx:0}
		input := "$5\r\nhello\r\n"
		expected := "hello"
		
		res, _ := rrp.ParserProtocol(input)
		
		if res != expected {
			t.Errorf("Expected '%s', but has '%s'\n", expected, res)
		}
	})
	
	t.Run("Test parser the bulk string $3war\\r\\n", func(t *testing.T) {
		rrp := RedisProtocolParser{Idx:0}
		input := "$3\r\nwar\r\n"
		expected := "war"

		res, _ := rrp.ParserProtocol(input)

		if res != expected {
			t.Errorf("Expected '%s', but has '%s'\n", expected, res)
		}
	})

	t.Run("Test parser the bulk string $0\\r\\n", func(t *testing.T) {
		rrp := RedisProtocolParser{Idx:0}
		input := "$0\r\n\r\n"
		expected := ""

		res, _ := rrp.ParserProtocol(input)

		if res != expected {
			t.Errorf("Expected '%s', but has '%s'\n", expected, res)
		}
	})

	t.Run("Test parser the bulk string $-1\\r\\n", func(t *testing.T) {
		rrp := RedisProtocolParser{Idx:0}
		input := "$-1\r\n"
		var expected interface{}

		res, _ := rrp.ParserProtocol(input)
		
		if res != nil {
			t.Errorf("Expected '%v', but has '%v'\n", expected, res)
		}
	})

	t.Run("Test parser the bulk string *2\\r\\n$3\\r\\ngod\\r\\n$3\\r\\nbad\\r\\n", func(t *testing.T) {
		rrp := RedisProtocolParser{Idx:0}
		input := "*2\r\n$3\r\ngod\r\n$3\r\nbad\r\n"
		expected := []string {"god", "bad"}

		res, _ := rrp.ParserProtocol(input)

		res2, ok := res.([]string)
		if !ok {
			t.Error("Incorrect data type.")
		}
		if !reflect.DeepEqual([]string(res2), expected) {
			t.Errorf("Expected '%v', but has '%v'\n", expected, res)
		}
	})

	t.Run("Test parser the bulk string *3\\r\\n$3\\r\\ngod\\r\\n$2\\r\\nof\r\n$3\\r\\nwar\\r\\n", func(t *testing.T) {
		rrp := RedisProtocolParser{Idx:0}
		input := "*3\r\n$3\r\ngod\r\n$2\r\nof\r\n$3\r\nwar\r\n"
		expected := []string {"god", "of", "war"}

		res, _ := rrp.ParserProtocol(input)

		res2, ok := res.([]string)
		if !ok {
			t.Error("Incorrect data type.")
		}
		if !reflect.DeepEqual([]string(res2), expected) {
			t.Errorf("Expected '%v', but has '%v'\n", expected, res)
		}
	})
}

func TestBulkStringEncode(t *testing.T) {
	t.Run("Test encode the string hello", func(t *testing.T) {
		input := "hello"
		expected := []byte("$5\r\nhello\r\n")

		res := BulkStringEncode(input)

		if string(res) != string(expected) {
			t.Errorf("Expected '%s', but has '%s'\n", string(expected), string(res))
		}
	})
}

func TestArrayEncode(t *testing.T) {
	t.Run("Test encode the array of strings", func(t *testing.T) {
		input := []string{"god", "of", "war"}
		expected := []byte("*3\r\n$3\r\ngod\r\n$2\r\nof\r\n$3\r\nwar\r\n")

		res := ArrayEncode(input)

		if string(res) != string(expected) {
			t.Errorf("Expected '%s', but has '%s'\n", string(expected), string(res))
		}
	})
}