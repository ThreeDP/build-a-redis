package define

const (
	DEFAULPORT = "6379"
)

const (
	BUFFERSIZE           = 1024
	RedisSimpleString    = "+"
	RedisSimpleErrors    = "-"
	RedisIntegers        = ":"
	RedisBulkStrings     = "$"
	RedisArrays          = "*"
	RedisBooleans        = "#"
	RedisDoubles         = ","
	RedisBigNumbers      = "("
	RedisBulkErrors      = "!"
	RedisVerbatimStrings = "="
	RedisMaps            = "%"
	RedisSets            = "~"
	RedisPushes          = ">"
)

var RedisSerialization = []string{
	RedisSimpleString,
	RedisSimpleErrors,
	RedisIntegers,
	RedisBulkStrings,
	RedisArrays,
	RedisBooleans,
	RedisDoubles,
	RedisBigNumbers,
	RedisBulkErrors,
	RedisVerbatimStrings,
	RedisMaps,
	RedisSets,
	RedisPushes,
}