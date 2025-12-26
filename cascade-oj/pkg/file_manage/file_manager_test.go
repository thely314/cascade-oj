package filemanage

import (
	"reflect"
	"testing"
)

func TestGetJudgeConfig(t *testing.T) {
	fileManager := NewJudgeConfigManager("../../cases")
	expectedConfig := &JudgeConfig{
		Score:               100,
		TimeResourceLimit:   500,
		MemoryResourceLimit: 16,
		CaseGroups: []CaseGroup{
			{
				GroupScore: 100,
				Cases: []TestCase{
					{
						SubScore:           50,
						InputFileLocation:  "1.in",
						AnswerFileLocation: "1.ans",
					},
					{
						SubScore:           50,
						InputFileLocation:  "2.in",
						AnswerFileLocation: "2.ans",
					},
				},
			},
		},
	}
	tests := []struct {
		problemId int64
		expected  *JudgeConfig
	}{
		{
			problemId: 1,
			expected:  expectedConfig,
		},
	}
	for _, tt := range tests {
		judgeConfig, err := fileManager.GetJudgeConfig(tt.problemId)
		if err != nil {
			t.Fatalf("GetJudgeConfig failed: %v", err)
		}
		if !reflect.DeepEqual(judgeConfig, tt.expected) {
			t.Errorf("GetJudgeConfig(%d) = %v, want %v", tt.problemId, judgeConfig, tt.expected)
		}
	}
}

func TestGetCase(t *testing.T) {
	fileManager := NewJudgeConfigManager("../../cases")
	tests := []struct {
		problemId   int64
		inputFile   string
		answerFile  string
		expectedIn  []byte
		expectedAns []byte
	}{
		{
			1,
			"1.in",
			"1.ans",
			[]byte("1 2"),
			[]byte("3"),
		},
	}
	for _, tt := range tests {
		in, ans, err := fileManager.GetCase(tt.problemId, tt.inputFile, tt.answerFile)
		if err != nil {
			t.Fatalf("GetCase failed: %v", err)
		}
		if !reflect.DeepEqual(in, tt.expectedIn) {
			t.Errorf("GetCase(%d, %q, %q) input = %v, want %v", tt.problemId, tt.inputFile, tt.answerFile, in, tt.expectedIn)
		}
		if !reflect.DeepEqual(ans, tt.expectedAns) {
			t.Errorf("GetCase(%d, %q, %q) answer = %v, want %v", tt.problemId, tt.inputFile, tt.answerFile, ans, tt.expectedAns)
		}
	}
}
