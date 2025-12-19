package util

import (
	pbGojudge "github.com/criyle/go-judge/pb"
)

// status
const (
	Pending  = "pending"  // 0
	Judging  = "judging"  // 1
	Accepted = "accepted" // 2

	SystemError         = "system_error"          // 3
	CompileError        = "compile_error"         // 4
	RuntimeError        = "runtime_error"         // 5
	WrongAnswer         = "wrong_answer"          // 6
	PresentationError   = "presentation_error"    // 7, not used currently
	TimeLimitExceeded   = "time_limit_exceeded"   // 8
	MemoryLimitExceeded = "memory_limit_exceeded" // 9

	OutputLimitExceeded = "output_limit_exceeded" // 10
	DangerousSyscall    = "dangerous_syscall"     // 11
)

// StatusToString converts i16 status code to string
func StatusToString(status_i16 int16) (res string) {
	switch status_i16 {
	case 0:
		res = Pending
	case 1:
		res = Judging
	case 2:
		res = Accepted
	case 3:
		res = SystemError
	case 4:
		res = CompileError
	case 5:
		res = RuntimeError
	case 6:
		res = WrongAnswer
	case 7:
		res = PresentationError
	case 8:
		res = TimeLimitExceeded
	case 9:
		res = MemoryLimitExceeded
	case 10:
		res = OutputLimitExceeded
	case 11:
		res = DangerousSyscall
	default:
		res = "Unknown"
	}
	return
}

// StatusToInt16 converts string to i16 status code
//
// unknown string returns 3, means system internal logical error
func StatusToInt16(status_string string) (res int16) {
	switch status_string {
	case Pending:
		res = 0
	case Judging:
		res = 1
	case Accepted:
		res = 2
	case SystemError:
		res = 3
	case CompileError:
		res = 4
	case RuntimeError:
		res = 5
	case WrongAnswer:
		res = 6
	case PresentationError:
		res = 7
	case TimeLimitExceeded:
		res = 8
	case MemoryLimitExceeded:
		res = 9
	case OutputLimitExceeded:
		res = 10
	case DangerousSyscall:
		res = 11
	default:
		// system internal logical error, mark as SystemError
		res = 3
	}
	return
}

// ParseGojudgeStatus converts pbGojudge status to cascade internal i16 status
func ParseGojudgeStatus(status pbGojudge.Response_Result_StatusType) (res int16) {
	switch status {
	case pbGojudge.Response_Result_Accepted:
		res = 2
	case pbGojudge.Response_Result_Invalid:
		res = 3
	case pbGojudge.Response_Result_NonZeroExitStatus, pbGojudge.Response_Result_Signalled:
		res = 5
	case pbGojudge.Response_Result_WrongAnswer, pbGojudge.Response_Result_PartiallyCorrect:
		res = 6
	case pbGojudge.Response_Result_TimeLimitExceeded:
		res = 8
	case pbGojudge.Response_Result_MemoryLimitExceeded:
		res = 9
	case pbGojudge.Response_Result_OutputLimitExceeded:
		res = 10
	case pbGojudge.Response_Result_DangerousSyscall:
		res = 11
	default:
		res = 3
	}
	return
}
