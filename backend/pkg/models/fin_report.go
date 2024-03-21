package models

import "github.com/google/uuid"

type FinancialReport struct {
	ID      uuid.UUID
	Revenue float32
	Costs   float32
	period  Period
}

type Period struct {
	StartYear    int
	StartQuarter int
	EndYear      int
	EndQuarter   int
}
