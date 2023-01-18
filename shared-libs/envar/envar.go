package envar

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

type typ interface {
	int | bool | string
}

func GetEnv[T typ](key string, def T) T {
	var res any = def

	envVal := strings.TrimSpace(os.Getenv(key))

	switch reflect.TypeOf(def).Kind() {
	case reflect.String:
		if envVal != "" {
			res = envVal
		}

	case reflect.Bool:
		if envVal != "" {
			val, err := strconv.ParseBool(strings.ToLower(envVal))
			evaluate(err, val, &res)
		}

	case reflect.Int:
		if envVal != "" {
			val, err := strconv.Atoi(envVal)
			evaluate(err, val, &res)
		}
	}

	return res.(T)
}

func evaluate[T typ](err error, val T, res *interface{}) {
	if err == nil {
		*res = val
	}
}
