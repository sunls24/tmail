package schedule

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
	"tmail/config"
	"tmail/ent"
	"tmail/ent/attachment"
	"tmail/ent/envelope"

	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	db  *ent.Client
	cfg *config.Config
}

func New(db *ent.Client, cfg *config.Config) *Scheduler {
	return &Scheduler{db: db, cfg: cfg}
}

func (s *Scheduler) Run() {
	go s.cleanUpExpired()
}

func (s *Scheduler) cleanUpExpired() {
	run(func() {
		go removeEmptyDir(s.cfg.BaseDir)
		expired := time.Now().Add(-time.Hour * 240)
		list, err := s.db.Attachment.Query().Where(attachment.HasOwnerWith(envelope.CreatedAtLT(expired))).All(context.TODO())
		if err != nil {
			log.Err(err).Msg("Attachment Query")
			return
		}
		for _, a := range list {
			_ = os.Remove(a.Filepath)
		}
		count, err := s.db.Attachment.Delete().Where(attachment.HasOwnerWith(envelope.CreatedAtLT(expired))).Exec(context.TODO())
		if err != nil {
			log.Err(err).Msg("Attachment Delete")
			return
		}
		if count > 0 {
			log.Info().Msgf("clean up attachment %d", count)
		}
		count, err = s.db.Envelope.Delete().Where(envelope.CreatedAtLT(expired)).Exec(context.TODO())
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

func run(fn func(), dur time.Duration) {
	for {
		select {
		case <-time.Tick(dur):
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Error().Msg(fmt.Sprint(err))
					}
				}()
				fn()
			}()
		}
	}
}
