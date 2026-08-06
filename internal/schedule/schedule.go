package schedule

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"
	"tmail/ent"
	"tmail/ent/attachment"
	"tmail/ent/envelope"
	"tmail/internal/api"

	"github.com/sunls24/gox"
	"github.com/sunls24/gox/cron"
)

type Scheduler struct {
	ctx context.Context
}

const cleanUpBatchSize = 500

func New(ctx context.Context) *Scheduler {
	return &Scheduler{ctx: ctx}
}

func (s *Scheduler) Run() {
	gox.SafeGo(s.cleanUpExpired)
}

func (s *Scheduler) cleanUpExpired() {
	cron.RunRepeat(func() {
		defer removeEmptyDir(api.Config(s.ctx).BaseDir)
		expired := time.Now().Add(-time.Hour * 240)
		attachmentCount, err := cleanUpAttachments(s.ctx, expired)
		if err != nil {
			slog.Error("Attachment cleanup", "err", err)
			return
		}
		if attachmentCount > 0 {
			slog.Info("clean up attachment", "count", attachmentCount)
		}

		envelopeCount, err := cleanUpEnvelopes(s.ctx, expired)
		if err != nil {
			slog.Error("Envelope cleanup", "err", err)
			return
		}
		if envelopeCount > 0 {
			slog.Info("clean up expired", "count", envelopeCount)
		}
	}, time.Hour*24)
}

func cleanUpAttachments(ctx context.Context, expired time.Time) (int, error) {
	queryCtx := context.Background()
	cursor := ""
	total := 0
	for {
		query := api.DB(ctx).Attachment.Query().
			Where(attachment.HasOwnerWith(envelope.CreatedAtLT(expired))).
			Order(ent.Asc(attachment.FieldID)).
			Limit(cleanUpBatchSize)
		if cursor != "" {
			query.Where(attachment.IDGT(cursor))
		}
		list, err := query.All(queryCtx)
		if err != nil {
			return total, err
		}
		if len(list) == 0 {
			return total, nil
		}
		cursor = list[len(list)-1].ID

		ids := make([]string, 0, len(list))
		for _, a := range list {
			err = os.Remove(a.Filepath)
			if err == nil || os.IsNotExist(err) {
				ids = append(ids, a.ID)
				continue
			}
			slog.Error("Attachment remove", "id", a.ID, "filepath", a.Filepath, "err", err)
		}
		if len(ids) > 0 {
			count, err := api.DB(ctx).Attachment.Delete().
				Where(attachment.IDIn(ids...)).
				Exec(queryCtx)
			if err != nil {
				return total, err
			}
			total += count
		}
	}
}

func cleanUpEnvelopes(ctx context.Context, expired time.Time) (int, error) {
	queryCtx := context.Background()
	total := 0
	for {
		list, err := api.DB(ctx).Envelope.Query().
			Select(envelope.FieldID).
			Where(envelope.CreatedAtLT(expired), envelope.Not(envelope.HasAttachments())).
			Limit(cleanUpBatchSize).
			All(queryCtx)
		if err != nil {
			return total, err
		}
		if len(list) == 0 {
			return total, nil
		}

		ids := make([]int, 0, len(list))
		for _, e := range list {
			ids = append(ids, e.ID)
		}
		count, err := api.DB(ctx).Envelope.Delete().
			Where(envelope.IDIn(ids...)).
			Exec(queryCtx)
		if err != nil {
			return total, err
		}
		total += count
	}
}

func removeEmptyDir(baseDir string) {
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() || path == baseDir {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			if err = os.Remove(path); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		slog.Error("removeEmptyDir", "err", err)
	}
}
