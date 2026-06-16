package cache

import (
	"context"
	"time"
)

// Data type should impl CacheRepo interface
type CacheRepo interface {
	GetCache(ctx context.Context, k string) ([]byte, error)
	SetCache(ctx context.Context, k string, v interface{}, ttl time.Duration) error
	DeleteCache(ctx context.Context, k ...string) error
}

// cache definitions
const (
	ContestListCacheKey      = "contest_list_cache" // contest_list_cache
	ContestDetailCacheKeyFmt = "contest_cache:%d"   // contest_cache:{contestID}

	ProblemListCacheKeyFmt   = "problem_list_cache:%d" // problem_list_cache:{contestID}
	ProblemDetailCacheKeyFmt = "problem_cache:%d"      // problem_cache:{problemID}

	SubmissionListCacheKeyFmt   = "submission_list_cache:%d:%d:%d" // submission_cache:{userID}:{contestID}:{problemID}
	SubmissionDetailCacheKeyFmt = "submission_cache:%d:%s"         // submission_cache:{userID}:{submissionUUID}

	SelfTestCacheKeyFmt = "self_test_cache:%d:%s" // self_test_cache:{userID}:{selfTestUUID}
)
