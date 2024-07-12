package utils

import "time"

type DateFactory interface {
	Now() time.Time
}

type nowDateFactory struct{}

type fixedDateFactory struct {
	date time.Time
}

func NewNowDateFactory() DateFactory {
	return &nowDateFactory{}
}

func NewFixedDateFactory(date time.Time) DateFactory {
	return &fixedDateFactory{date: date}
}

func (f *nowDateFactory) Now() time.Time {
	return time.Now()
}

func (f *fixedDateFactory) Now() time.Time {
	return f.date
}
