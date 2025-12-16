package filemanage

import (
	"os"
	"strconv"

	"github.com/pelletier/go-toml/v2"
)

type JudgeConfigManager struct {
	configLocation string
}

type JudgeConfig struct {
	Score               int
	TimeResourceLimit   int64
	MemoryResourceLimit int64
	CaseGroups          []CaseGroup
}

type TestCase struct {
	SubScore           int
	InputFileLocation  string
	AnswerFileLocation string
}

type CaseGroup struct {
	GroupScore int
	Cases      []TestCase
}

func NewJudgeConfigManager(fileLocation string) *JudgeConfigManager {
	return &JudgeConfigManager{
		configLocation: fileLocation,
	}
}

func (fileManager *JudgeConfigManager) GetJudgeConfig(problemId int64) (*JudgeConfig, error) {
	tomlFileBytes, err := ReadTomlFile(problemId, fileManager.configLocation)
	if err != nil {
		return nil, err
	}
	var config JudgeConfig
	err = toml.Unmarshal(tomlFileBytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (fileManager *JudgeConfigManager) GetCase(problemId int64, inputFileLocation string, ansFileLocation string) (in []byte, ans []byte, err error) {
	baseLocation := fileManager.configLocation + "/" + strconv.FormatInt(problemId, 10) + "/testcase/"
	in, err = os.ReadFile(baseLocation + inputFileLocation)
	if err != nil {
		return nil, nil, err
	}
	ans, err = os.ReadFile(baseLocation + ansFileLocation)
	if err != nil {
		return nil, nil, err
	}
	return in, ans, nil
}
