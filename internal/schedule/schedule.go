package schedule

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"time"
	"tmail/ent/attachment"
	"tmail/ent/envelope"
	"tmail/internal/api"

	"github.com/rs/zerolog/log"
	"github.com/sunls24/gox"
	"github.com/sunls24/gox/cron"
)

type Scheduler struct {
	ctx context.Context
}

func New(ctx context.Context) *Scheduler {
	return &Scheduler{ctx: ctx}
}

func (s *Scheduler) Run() {
	gox.SafeGo(s.cleanUpExpired)
}

func (s *Scheduler) cleanUpExpired() {
	cron.RunRepeat(func() {
		gox.SafeGo(func() {
			removeEmptyDir(api.Config(s.ctx).BaseDir)
		})
		expired := time.Now().Add(-time.Hour * 240)
		list, err := api.DB(s.ctx).Attachment.Query().Where(attachment.HasOwnerWith(envelope.CreatedAtLT(expired))).All(context.Background())
		if err != nil {
			log.Err(err).Msg("Attachment Query")
			return
		}
		for _, a := range list {
			_ = os.Remove(a.Filepath)
		}
		count, err := api.DB(s.ctx).Attachment.Delete().Where(attachment.HasOwnerWith(envelope.CreatedAtLT(expired))).Exec(context.Background())
		if err != nil {
			log.Err(err).Msg("Attachment Delete")
			return
		}
		if count > 0 {
			log.Info().Msgf("clean up attachment %d", count)
		}
		count, err = api.DB(s.ctx).Envelope.Delete().Where(envelope.CreatedAtLT(expired)).Exec(context.Background())
		if err != nil {
			log.Err(err).Msg("Envelope Delete")
			return
		}
		if count > 0 {
			log.Info().Msgf("clean up expired %d", count)
		}
	}, time.Hour*24)
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
		log.Err(err).Msg("removeEmptyDir")
	}
}
