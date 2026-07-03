package data

import (
	"context"
	"strconv"
	"strings"

	"archive/zip"
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/problemset_includes"
	"cascade-oj/ent/problemtemplate"
	"cascade-oj/pkg/mq"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pelletier/go-toml/v2"
)

type ProblemRepo struct {
	data *Data
	log  *log.Helper
}

type TestCaseConfig struct {
	Score               int         `toml:"Score"`
	TimeResourceLimit   int64       `toml:"TimeResourceLimit"`
	MemoryResourceLimit int64       `toml:"MemoryResourceLimit"`
	CaseGroups          []CaseGroup `toml:"CaseGroups"`
}

type CaseGroup struct {
	GroupScore int    `toml:"GroupScore"`
	Cases      []Case `toml:"Cases"`
}

type Case struct {
	SubScore           int    `toml:"SubScore"`
	InputFileLocation  string `toml:"InputFileLocation"`
	AnswerFileLocation string `toml:"AnswerFileLocation"`
}

func NewProblemRepo(data *Data, logger log.Logger) biz.ProblemRepo {
	return &ProblemRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func mapEntStatusToBizStatus(entStatus problem.UseStatus) biz.ProblemStatus {
	switch entStatus {
	case problem.UseStatusUnavailable:
		return biz.ProblemStatusUnavailable
	case problem.UseStatusAvailable:
		return biz.ProblemStatusAvailable
	case problem.UseStatusUsing:
		return biz.ProblemStatusUsing
	case problem.UseStatusDeleted:
		return biz.ProblemStatusDeleted
	default:
		return biz.ProblemStatusUnavailable
	}
}

func (problemRepo *ProblemRepo) GetProblems(ctx context.Context, contestID int64) ([]*biz.Problem, error) {
	if contestID == 0 {
		allProblems, err := problemRepo.data.db.Problem.Query().
			Where(problem.UseStatusNEQ(problem.UseStatusDeleted)).
			All(ctx)
		if err != nil {
			return nil, err
		}
		res := make([]*biz.Problem, 0, len(allProblems))
		for _, p := range allProblems {
			res = append(res, &biz.Problem{
				ID:            p.ID,
				Title:         p.Title,
				Description:   p.Description,
				TimeLimitMs:   int32(p.TimeLimitMs),
				MemoryLimitKB: int32(p.MemoryLimitKB),
				Status:        mapEntStatusToBizStatus(p.UseStatus),
			})
		}
		return res, nil
	}

	entProblemIncludes, err := problemRepo.data.db.ProblemSet_Includes.
		Query().
		WithProblem().
		Where(problemset_includes.ProblemSetIDEQ(contestID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	problems := make([]*biz.Problem, 0, len(entProblemIncludes))
	for _, item := range entProblemIncludes {
		p := item.Edges.Problem
		problems = append(problems, &biz.Problem{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			TimeLimitMs:   int32(p.TimeLimitMs),
			MemoryLimitKB: int32(p.MemoryLimitKB),
			Status:        mapEntStatusToBizStatus(p.UseStatus),
		})
	}
	return problems, nil
}

func (problemRepo *ProblemRepo) GetSingleProblem(ctx context.Context, problemID int64) (*biz.DetailedProblem, error) {
	entProblem, err := problemRepo.data.db.Problem.
		Query().
		Where(problem.IDEQ(problemID)).
		WithCreator().
		WithTemplates().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	bizTemplates := make([]*biz.ProblemTemplate, 0, len(entProblem.Edges.Templates))
	for _, t := range entProblem.Edges.Templates {
		bizTemplates = append(bizTemplates, &biz.ProblemTemplate{
			Name:    t.Name,
			Content: t.Content,
		})
	}

	return &biz.DetailedProblem{
		Problem: biz.Problem{
			ID:            entProblem.ID,
			Title:         entProblem.Title,
			TimeLimitMs:   int32(entProblem.TimeLimitMs),
			MemoryLimitKB: int32(entProblem.MemoryLimitKB),
			Status:        mapEntStatusToBizStatus(entProblem.UseStatus),
			Templates:     bizTemplates,
		},
		CreatorUsername: entProblem.Edges.Creator.Username,
		Description:     entProblem.Description,
	}, nil
}

func (problemRepo *ProblemRepo) PostProblem(ctx context.Context, creatorID int64, info biz.ProblemCreateInfo) (int64, error) {
	return WithTx(ctx, problemRepo.data.db, func(tx *ent.Tx) (int64, error) {
		// 1. 存主表
		newProblem, err := tx.Problem.Create().
			SetTitle(info.Title).
			SetTimeLimitMs(int(info.TimeLimitMs)).
			SetMemoryLimitKB(int(info.MemoryLimitKB)).
			SetDescription(info.Description).
			SetCreatorID(creatorID).
			SetJudgeConfigID(1).
			SetUseStatus(problem.UseStatusUnavailable).
			Save(ctx)
		if err != nil {
			return 0, err
		}

		// 2. 循环存子表模板 (支持多文件)
		for _, t := range info.Templates {
			_, err = tx.ProblemTemplate.Create().
				SetName(t.Name).
				SetContent(t.Content).
				SetProblem(newProblem).
				Save(ctx)
			if err != nil {
				return 0, err
			}
		}
		return newProblem.ID, nil
	})
}

func (problemRepo *ProblemRepo) PutProblem(ctx context.Context, info biz.ProblemEditInfo) (bool, error) {
	_, err := WithTx(ctx, problemRepo.data.db, func(tx *ent.Tx) (int64, error) {
		// 1. 更新主表
		err := tx.Problem.UpdateOneID(info.ID).
			SetTitle(info.Title).
			SetTimeLimitMs(int(info.TimeLimitMs)).
			SetMemoryLimitKB(int(info.MemoryLimitKB)).
			SetDescription(info.Description).
			Exec(ctx)
		if err != nil {
			return 0, err
		}

		// 2. 更新模板：先删除，再插入
		_, err = tx.ProblemTemplate.Delete().
			Where(problemtemplate.HasProblemWith(problem.ID(info.ID))).
			Exec(ctx)
		if err != nil {
			return 0, err
		}

		for _, t := range info.Templates {
			_, err = tx.ProblemTemplate.Create().
				SetName(t.Name).
				SetContent(t.Content).
				SetProblemID(info.ID).
				Save(ctx)
			if err != nil {
				return 0, err
			}
		}
		return 0, nil
	})

	if err != nil {
		return false, err
	}

	// delete user cache for problem
	err = problemRepo.data.publishProblemInvalidation(ctx, &mq.ProblemCacheMsg{
		ProblemID: info.ID,
		Scale:     "single",
	})
	if err != nil {
		problemRepo.log.Errorf("[cache] failed to publish problem cache invalidation message for problem single: %v", err)
	}

	return true, nil
}

// 增加过滤条件，防止逻辑误复活
func (problemRepo *ProblemRepo) DeleteProblem(ctx context.Context, problemID int64) (bool, error) {
	err := problemRepo.data.db.Problem.Update().
		Where(problem.IDEQ(problemID)).
		SetUseStatus(problem.UseStatusDeleted).
		Exec(ctx)

	if err != nil {
		return false, err
	}

	// delete user cache for problem
	err = problemRepo.data.publishProblemInvalidation(ctx, &mq.ProblemCacheMsg{
		ProblemID: problemID,
		Scale:     "single",
	})
	if err != nil {
		problemRepo.log.Errorf("[cache] failed to publish problem cache invalidation message for problem single: %v", err)
	}

	return true, nil
}

func (problemRepo *ProblemRepo) PublishProblem(ctx context.Context, problemID int64) (bool, error) {
	err := problemRepo.data.db.Problem.Update().
		Where(problem.IDEQ(problemID), problem.UseStatusNEQ(problem.UseStatusDeleted)). // 防御：已删除题目不可发布
		SetUseStatus(problem.UseStatusAvailable).
		Exec(ctx)
	return err == nil, err
}

func (problemRepo *ProblemRepo) DisableProblem(ctx context.Context, problemID int64) (bool, error) {
	err := problemRepo.data.db.Problem.Update().
		Where(problem.IDEQ(problemID), problem.UseStatusNEQ(problem.UseStatusDeleted)). // 防御
		SetUseStatus(problem.UseStatusUnavailable).
		Exec(ctx)
	return err == nil, err
}

func (r *ProblemRepo) SaveTestCases(ctx context.Context, problemID int64, file io.Reader) error {
	// 1. 准备目录
	basePath := filepath.Join("cases", strconv.FormatInt(problemID, 10), "testcase")
	_ = os.RemoveAll(basePath)
	_ = os.MkdirAll(basePath, 0755)

	// 2. 保存临时 zip
	zipPath := filepath.Join(basePath, "upload.zip")
	tmpFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	_, _ = io.Copy(tmpFile, file)
	tmpFile.Close()

	// 3. 解压文件
	zr, err := zip.OpenReader(zipPath) // 改名 zr，避免与接收者 r 冲突
	if err != nil {
		return err
	}
	defer zr.Close()

	var inputFiles []string
	for _, f := range zr.File {
		fpath := filepath.Join(basePath, f.Name) // 使用 fpath

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, f.Mode())
			continue
		}

		// 写入解压后的文件
		dstFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}
		_, _ = io.Copy(dstFile, srcFile)

		dstFile.Close()
		srcFile.Close()

		if filepath.Ext(f.Name) == ".in" {
			inputFiles = append(inputFiles, f.Name)
		}
	}

	// 4. 生成 config.toml
	config := TestCaseConfig{
		Score:               100,
		TimeResourceLimit:   500,
		MemoryResourceLimit: 16,
	}

	group := CaseGroup{GroupScore: 100}
	for _, in := range inputFiles {
		ans := strings.TrimSuffix(in, ".in") + ".ans"
		group.Cases = append(group.Cases, Case{
			SubScore:           100 / len(inputFiles),
			InputFileLocation:  in,
			AnswerFileLocation: ans,
		})
	}
	config.CaseGroups = append(config.CaseGroups, group)

	tomlData, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(basePath, "config.toml"), tomlData, 0644)
}
