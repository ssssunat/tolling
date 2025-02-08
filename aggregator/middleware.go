package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ssssunat/tolling/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}
func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"func": "Aggregate Distance",
		}).Info("Aggregate Distance")
	}(time.Now())

	err = m.next.AggregateDistance(dist)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}

		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"obuID": obuID,
			"amount": amount,
			"distance": distance,
		}).Info("Calculate Invoice")
	}(time.Now())

	inv, err = m.next.CalculateInvoice(obuID)
	return
}
