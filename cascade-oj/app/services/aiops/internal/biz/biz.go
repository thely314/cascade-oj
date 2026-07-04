package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewAnalyzer,
	NewReporter,
	NewPrometheusTool,
	NewLogReaderTool,
	NewLLMClient,
	NewPoller,
)
