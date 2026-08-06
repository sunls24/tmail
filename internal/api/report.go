package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"tmail/ent"

	"github.com/jhillyerd/enmime/v2"
	"github.com/sunls24/gox"
	"github.com/sunls24/gox/notifier"
	"github.com/sunls24/gox/server"
)

func Report(ctx context.Context) (*server.Reply, error) {
	ec := server.EchoContext(ctx)
	to := ec.QueryParam("to")
	if to == "" {
		return nil, server.BadParam()
	}
	envelope, err := enmime.ReadEnvelope(ec.Request().Body)
	if err != nil {
		return nil, err
	}
	subject := envelope.GetHeader("subject")
	from := envelope.GetHeader("from")
	if from == "" {
		return server.OK(nil), nil
	}
	content := envelope.HTML
	if content == "" {
		content = envelope.Text
	}

	slog.Debug("Report", "to", to, "from", from, "subject", subject)
	e, err := DB(ctx).Envelope.Create().
		SetTo(to).
		SetFrom(from).
		SetSubject(subject).
		SetContent(content).
		Save(ctx)
	if err == nil {
		notifyEnvelope := envelopeSummary(e)
		attachmentCtx := context.WithoutCancel(ctx)
		gox.SafeGo(func() {
			saveAttachment(attachmentCtx, envelope.Attachments, to, e.ID)
			notifier.Notify(e.To, notifyEnvelope)
			notifier.Notify(subAll, notifyEnvelope)
		})
	}
	return server.OK(nil), err
}

func envelopeSummary(e *ent.Envelope) *ent.Envelope {
	return &ent.Envelope{
		ID:        e.ID,
		To:        e.To,
		From:      e.From,
		Subject:   e.Subject,
		CreatedAt: e.CreatedAt,
	}
}

func saveAttachment(ctx context.Context, attachments []*enmime.Part, to string, ownerID int) {
	const maxSize = 200000000 // 200M
	if len(attachments) == 0 {
		return
	}

	cfg := Config(ctx)
	dir := filepath.Join(cfg.BaseDir, gox.MD5(to)[:16])
	if err := os.MkdirAll(dir, 0o755); err != nil {
		slog.Error("MkdirAll", "err", err)
		return
	}

	for i, a := range attachments {
		if a.FileName == "" || len(a.Content) > maxSize {
			continue
		}

		name := gox.MD5(fmt.Sprintf("%d:%d:%s", ownerID, i, a.FileName))
		fp := filepath.Join(dir, name)
		slog.Info("Attachment", "filename", a.FileName, "filepath", fp)
		if err := os.WriteFile(fp, a.Content, 0o644); err != nil {
			slog.Error("WriteFile", "err", err)
			continue
		}

		_, err := DB(ctx).Attachment.Create().
			SetID(filepath.Base(dir) + name[:6] + gox.RandStr(4)).
			SetFilename(a.FileName).
			SetFilepath(fp).
			SetContentType(a.ContentType).
			SetOwnerID(ownerID).
			Save(ctx)
		if err != nil {
			_ = os.Remove(fp)
			slog.Error("Attachment Save", "err", err)
		}
	}
}
