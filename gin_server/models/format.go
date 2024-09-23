package models

import "time"

type FormatData struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	MachineID string    `json:"machine_id"`
	Time      time.Time `json:"time"`
	Judgement bool      `json:"judgement"`
	Traceability
}

type Traceability []TraceabilityDetail

func (t *Traceability) AddTraceabilityDetail(name string, data any) {
	*t = append(*t, TraceabilityDetail{name, data})
}

type TraceabilityDetail struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}
