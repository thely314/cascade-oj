package mq

import "fmt"

type SelfTestMessage struct {
	UUID       string `json:"uuid,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	ProblemID  int64  `json:"problem_id,omitempty"`
	Code       string `json:"code,omitempty"`
	Language   string `json:"language,omitempty"`
	Input      string `json:"input,omitempty"`
	IsCompiled bool   `json:"is_compiled,omitempty"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	TimeCost   uint64 `json:"time_cost,omitempty"`
	MemoryCost uint64 `json:"memory_cost,omitempty"`
	Token      string `json:"token,omitempty"`
}

func (sm *SelfTestMessage) String() string {
	return fmt.Sprintf(
		"SelfTestMessage{UUID: %s, UserId: %d, ProblemID: %d, Code: %s, Language: %s, Input: %s, IsCompiled: %t, Stdout: %s, Stderr: %s, Time: %d, Memory: %d, Token: %s}",
		sm.UUID,
		sm.UserID,
		sm.ProblemID,
		sm.Code,
		sm.Language,
		sm.Input,
		sm.IsCompiled,
		sm.Stdout,
		sm.Stderr,
		sm.TimeCost,
		sm.MemoryCost,
		sm.Token,
	)
}
