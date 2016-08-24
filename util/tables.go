package util

import (
	"strings"
)

func TableNameForMetricId(id string) string {
	return strings.Replace(id, ".", "_", -1)
}
