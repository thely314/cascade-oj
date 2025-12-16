package gojudge

import (
	"cascade-oj/app/services/judge/internal/conf"
	"strings"
)

type Commands = map[string]Command

type Command struct {
	// env variables for compile and run commands
	Env []string

	// compile command and its arguments
	//
	// already split into string array
	Compile []string
	// run command and its arguments
	//
	// already split into string array
	Run []string

	// source file name
	Source string
	// target file name after compilation
	Target string

	// resource limits for compile
	CompileConfig *conf.ExecConfig
	// resource limits for run
	RunConfig *conf.ExecConfig
}

// construct Commands from configuration
//
// # Warning
//
// configuration needs to ensure no space in file names or arguments
func NewCommand(
	enable []string,
	env map[string][]string,
	compile map[string]string,
	run map[string]string,
	source map[string]string,
	target map[string]string,
	configs map[string]*conf.LanguageConfig) Commands {
	res := make(map[string]Command)
	for _, language := range enable {
		environment, is_ok := env[language]
		if !is_ok {
			environment = env["default"]
		}
		compile_command, is_ok := compile[language]
		if !is_ok {
			// no compile command means interpreted language like python
			compile_command = ""
		}
		run_command, is_ok := run[language]
		if !is_ok {
			run_command = run["default"]
		}
		src, is_ok := source[language]
		if !is_ok {
			src = source["default"]
		}
		tgt, is_ok := target[language]
		if !is_ok {
			tgt = target["default"]
		}
		config, is_ok := configs[language]
		if !is_ok {
			config = configs["default"]
		}

		// split commands into string array
		// Warning: configuration needs to ensure no space in file names or arguments
		compileCmd := strings.Fields(compile_command)
		runCmd := strings.Fields(run_command)

		res[language] = Command{
			Env:           environment,
			Compile:       compileCmd,
			Run:           runCmd,
			Source:        src,
			Target:        tgt,
			CompileConfig: config.Compile,
			RunConfig:     config.Run,
		}
	}
	return res
}
