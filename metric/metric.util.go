package metric

import (
	"errors"
)

func (dp *MetricDatapoint) Validate() error {
	if dp == nil {
		return errors.New("Datapoint cannot be null.")
	}
	return nil
}
