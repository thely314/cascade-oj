package gojudge

import (
	"context"
	"errors"
	"strconv"

	pbGojudge "github.com/criyle/go-judge/pb"
)

type GoJudge struct {
	Client *pbGojudge.ExecutorClient

	// supported languages and their commands
	Commands *Commands
}

// perform code compilation
//
// The first return value is the fileID of the compiled target file,
// fileID will be used in case judge, by requestCached(fileID string)
//
// *pbGojudge.Response_Result.Time is ns, *pbGojudge.Response_Result.Memory is Bytes
func (gj *GoJudge) Compile(code []byte, language string, requestUUID string) (string, *pbGojudge.Response_Result, error) {
	command, is_ok := (*gj.Commands)[language]
	if !is_ok {
		return "", nil, errors.New("unsupported language: " + language)
	}

	// no compile command means interpreted language like python
	if len(command.Compile) == 0 {
		// upload source code as target file
		fileID, err := gj.AddFile(command.Target, code)
		if err != nil {
			return "", nil, err
		}
		return fileID, nil, nil
	}

	tempResult, err := (*gj.Client).Exec(context.Background(), &pbGojudge.Request{
		RequestID: requestUUID,
		// https://github.com/criyle/go-judge/blob/v1.8.5/README.cn.md
		Cmd: []*pbGojudge.Request_CmdType{
			{
				Args: command.Compile,
				Env:  command.Env,
				// define stdin, stdout, stderr files
				Files: []*pbGojudge.Request_File{
					requestMemory([]byte{}),
					requestPipe("stdout", command.RunConfig.StdoutMaxSize),
					requestPipe("stderr", command.RunConfig.StderrMaxSize),
				},
				// ms -> ns
				CpuTimeLimit: command.CompileConfig.CpuTimeLimit * 1000 * 1000,
				// CPU usage rate limit
				// Linux only, 1000 means one core 100%
				CpuRateLimit: command.CompileConfig.CpuRateLimit,
				// waiting time limit in ms -> ns
				ClockTimeLimit: command.CompileConfig.ClockTimeLimit * 1000 * 1000,
				// MemoryLimit MiB -> Bytes
				MemoryLimit: command.CompileConfig.MemoryLimit * 1024 * 1024,
				// threads limit, 是线程数不是进程数
				// go-judge 你无敌了
				ProcLimit: command.CompileConfig.ProcLimit,
				// copy files into sandbox before execution
				CopyIn: map[string]*pbGojudge.Request_File{
					command.Source: requestMemory(code),
				},
				// copy out files from sandbox after execution
				CopyOut: []*pbGojudge.Request_CmdCopyOutFile{
					requestCopyOutFile("stdout"),
					requestCopyOutFile("stderr"),
				},
				// return cached fileIDs after execution
				// files can be downloaded using /file/:fileId
				CopyOutCached: []*pbGojudge.Request_CmdCopyOutFile{
					requestCopyOutFile(command.Target),
				},
			},
		},
	})

	result, err := handleExecError(requestUUID, tempResult, err)
	if err != nil {
		return "", result, err
	}
	return result.FileIDs[command.Target], result, nil
}

// perform a case judge
//
// used for self-test (when problemId = -1) and each test case from problem
func (gj *GoJudge) CaseJudge(
	problemId int64,
	caseName string,
	input []byte,
	language string,
	targetID string,
	requestUUID string,
	cpuTimeLimit uint64,
	clockTimeLimit uint64,
	memoryLimit uint64) (*pbGojudge.Response_Result, error) {
	command, is_ok := (*gj.Commands)[language]
	if !is_ok {
		return nil, errors.New("unsupported language: " + language)
	}

	var fileList []*pbGojudge.Request_File
	if problemId != -1 {
		// normal test case
		// get input from local file
		fileList = []*pbGojudge.Request_File{
			requestLocal(problemId, caseName),
			requestPipe("stdout", command.RunConfig.StdoutMaxSize),
			requestPipe("stderr", command.RunConfig.StderrMaxSize),
		}
	} else {
		// self-test case
		// get input from memory
		fileList = []*pbGojudge.Request_File{
			requestMemory(input),
			requestPipe("stdout", command.RunConfig.StdoutMaxSize),
			requestPipe("stderr", command.RunConfig.StderrMaxSize),
		}
	}

	tempResult, err := (*gj.Client).Exec(context.Background(), &pbGojudge.Request{
		RequestID: requestUUID,
		// https://github.com/criyle/go-judge/blob/v1.8.5/README.cn.md
		Cmd: []*pbGojudge.Request_CmdType{
			{
				Args:           command.Run,
				Env:            command.Env,
				Files:          fileList,
				CpuTimeLimit:   min(command.RunConfig.CpuTimeLimit, cpuTimeLimit) * 1000 * 1000,
				CpuRateLimit:   command.RunConfig.CpuRateLimit,
				ClockTimeLimit: min(command.RunConfig.ClockTimeLimit, clockTimeLimit) * 1000 * 1000,
				MemoryLimit:    min(command.RunConfig.MemoryLimit, memoryLimit) * 1024 * 1024,
				ProcLimit:      command.RunConfig.ProcLimit,
				CopyIn: map[string]*pbGojudge.Request_File{
					command.Target: requestCached(targetID),
				},
			},
		},
	})

	result, err := handleExecError(requestUUID, tempResult, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (gj *GoJudge) AddFile(file string, content []byte) (fileID string, err error) {
	fileid, err := (*gj.Client).FileAdd(context.Background(), &pbGojudge.FileContent{
		Name:    file,
		Content: content,
	})
	if err != nil {
		return "", err
	}
	return fileid.FileID, nil
}

func (gj *GoJudge) DeleteFile(fileID string) error {
	_, err := (*gj.Client).FileDelete(context.Background(), &pbGojudge.FileID{
		FileID: fileID,
	})
	if err != nil {
		return err
	}
	return nil
}

func handleExecError(requestUUID string, response *pbGojudge.Response, err error) (*pbGojudge.Response_Result, error) {
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, errors.New("request " + requestUUID + ": no results returned from go-judge")
	}

	// check for execution errors
	result := response.Results[0]
	if result.Error != "" {
		return result, errors.New("request " + requestUUID + ": " + result.Error)
	}
	if result.FileError != nil {
		errMessage := "request " + requestUUID + ":\n"
		for _, fileError := range result.FileError {
			errMessage += "file error: " + fileError.Message + "\n"
		}
		return result, errors.New(errMessage)
	}

	// Check for non-zero exit status
	if result.ExitStatus != 0 {
		return result, errors.New("request " + requestUUID + ": non-zero exit status: " + string(result.ExitStatus))
	}
	return result, nil
}

func requestMemory(content []byte) *pbGojudge.Request_File {
	return &pbGojudge.Request_File{
		File: &pbGojudge.Request_File_Memory{
			Memory: &pbGojudge.Request_MemoryFile{
				Content: content,
			},
		},
	}
}

// request local file to upload into sandbox
func requestLocal(problemId int64, caseName string) *pbGojudge.Request_File {
	return &pbGojudge.Request_File{
		File: &pbGojudge.Request_File_Local{
			Local: &pbGojudge.Request_LocalFile{
				Src: "/cases/" + strconv.FormatInt(problemId, 10) + "/testcase/" + caseName,
			},
		},
	}
}

func requestPipe(name string, max int64) *pbGojudge.Request_File {
	return &pbGojudge.Request_File{
		File: &pbGojudge.Request_File_Pipe{
			Pipe: &pbGojudge.Request_PipeCollector{
				Name: name,
				Max:  max,
			},
		},
	}
}

func requestCopyOutFile(name string) *pbGojudge.Request_CmdCopyOutFile {
	return &pbGojudge.Request_CmdCopyOutFile{
		Name: name,
	}
}

func requestCached(fileID string) *pbGojudge.Request_File {
	return &pbGojudge.Request_File{
		File: &pbGojudge.Request_File_Cached{
			Cached: &pbGojudge.Request_CachedFile{
				FileID: fileID,
			},
		},
	}
}
