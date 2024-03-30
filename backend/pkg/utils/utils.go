package utils

import (
	"fmt"
	"strings"
)

type Filters map[string]string

var sqlInjections = []string{"' or 1=1; --", "' or '1'='1", "select", "drop", "insert", "1=1"}

func (f Filters) Validate() (err error) {
	for k, v := range f {
		for _, val := range sqlInjections {
			kLow, vLow := strings.ToLower(k), strings.ToLower(v)
			if strings.Contains(kLow, val) || strings.Contains(vLow, val) {
				return fmt.Errorf("sql-инъекция")
			}
		}
	}

	return nil
}

func (f Filters) ParseToSql() (str string, err error) {
	var i int
	for k, v := range f {
		if v[0] == '-' {
			str += fmt.Sprintf("%s < '%s'", k, v[1:])
		} else if v[len(v)-1] == '-' {
			str += fmt.Sprintf("%s > '%s'", k, v[:len(v)-1])
		} else if splitted := strings.Split(v, "-"); len(splitted) == 2 {
			str += fmt.Sprintf("%s between '%s' and '%s'", k, splitted[0], splitted[1])
		} else {
			str += fmt.Sprintf("%s = '%s'", k, v)
		}

		if i != len(f)-1 {
			str += " and "
		}
		i++
	}

	return str, nil
}
