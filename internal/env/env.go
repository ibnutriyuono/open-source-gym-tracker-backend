package env

import (
	"os"
	"strconv"
)

func GetString(key, callback string) string{
	val, ok := os.LookupEnv(key)
	if !ok {
		return callback
	}

	return val
}

func GetInt(key string, callback int) int{
	val, ok := os.LookupEnv(key)
	if !ok {
		return callback
	}

	intVal, err := strconv.Atoi((val))
	if err != nil {
		return callback
	}

 	return intVal
}

func GetBool(key string, callback bool) bool  {
	val, ok := os.LookupEnv(key)
	if !ok {
		return callback
	}

	boolVal, err := strconv.ParseBool((val))
	if err != nil {
		return callback
	}

 	return boolVal
}