package mq

import (
	"fmt"
	"time"
)

type SubmissionMessage struct {
	UUID         string    `json:"uuid,omitempty"`
	UserID       int64     `json:"user_id,omitempty"`
	ProblemID    int64     `json:"problem_id,omitempty"`
	ProblemSetID int64     `json:"problem_set_id,omitempty"`
	Code         string    `json:"code,omitempty"`
	Status       int16     `json:"status,omitempty"`
	Score        int       `json:"Score,omitempty"`
	CreateTime   time.Time `json:"create_time"`
	TimeCost     uint64    `json:"total_time_cost,omitempty"`
	MemoryCost   uint64    `json:"max_memory_cost,omitempty"`
	Language     string    `json:"language,omitempty"`
	Stderr       string    `json:"stderr,omitempty"`
	CaseVersion  int16     `json:"case_version,omitempty"`
	Token        string    `json:"token,omitempty"`
}

func (sm *SubmissionMessage) String() string {
	return fmt.Sprintf(
		"SubmissionMessage{UUID: %s, UserId: %d, ProblemID: %d, ProblemSetID: %d, Code: %s, Status: %d, Score: %d, CreateTime: %s, TimeCost: %d, MemoryCost: %d, Language: %s, Stderr: %s, CaseVersion: %d, Token: %s}",
		sm.UUID,
		sm.UserID,
		sm.ProblemID,
		sm.ProblemSetID,
		sm.Code,
		sm.Status,
		sm.Score,
		sm.CreateTime,
		sm.TimeCost,
		sm.MemoryCost,
		sm.Language,
		sm.Stderr,
		sm.CaseVersion,
		sm.Token,
	)
}
