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

func (ms *MetricSeries) Validate() error {
	if ms == nil {
		return errors.New("Metric series cannot be nil.")
	}

	if ms.DataType != MetricSeries_NUMBER {
		return errors.New("Nothing but number data types are supported.")
	}

	if ms.Title == "" {
		return errors.New("Title is required.")
	}

	if ms.Id == "" {
		return errors.New("ID is required.")
	}

	return nil
}
