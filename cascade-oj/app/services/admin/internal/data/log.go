package data

import (
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent/systemlog"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type LogRepo struct {
	data *Data
	log  *log.Helper
}

func NewLogRepo(data *Data, logger log.Logger) biz.LogRepo {
	return &LogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (logRepo *LogRepo) GetLogs(ctx context.Context, request biz.LogRequestInfo) ([]*biz.Log, error) {
	entLogs, err := logRepo.data.db.SystemLog.
		Query().
		Select(
			systemlog.FieldID,
			systemlog.FieldLogInfo,
			systemlog.FieldLogTime,
		).
		Offset(int(max((request.Page-1), 0) * request.PageSize)).
		Limit(int(request.PageSize)).All(ctx)
	if err != nil {
		return nil, err
	}
	logs := make([]*biz.Log, 0, len(entLogs))
	for i := 0; i < len(entLogs); i++ {
		logs = append(logs,
			&biz.Log{
				ID:        entLogs[i].ID,
				Message:   entLogs[i].LogInfo,
				TimeStamp: entLogs[i].LogTime,
			})
	}
	return logs, nil
}
