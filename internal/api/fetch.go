package api

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"tmail/ent"
	"tmail/ent/attachment"
	"tmail/ent/envelope"
	"tmail/ent/predicate"

	"github.com/sunls24/gox"
	"github.com/sunls24/gox/notifier"
	"github.com/sunls24/gox/server"
)

const subAll = "all"

const (
	defaultFetchLimit = 30
	maxFetchLimit     = 100
)

type ReqFetch struct {
	To       string `query:"to"`
	Since    string `query:"since"`
	BeforeID int    `query:"before_id"`
	Limit    int    `query:"limit"`
}

func Fetch(ctx context.Context, req ReqFetch) ([]*ent.Envelope, error) {
	if req.To == "" {
		return nil, server.BadParam()
	}
	limit := req.Limit
	if limit == 0 {
		limit = defaultFetchLimit
	}
	if limit < 1 || limit > maxFetchLimit || req.BeforeID < 0 {
		return nil, server.BadParam()
	}
	admin := req.To == Config(ctx).AdminAddress
	since := time.Time{}
	if !admin {
		if req.Since != "" {
			ts, err := strconv.ParseInt(req.Since, 10, 64)
			if err == nil && ts >= 0 {
				since = time.Unix(ts, 0)
			}
		}
	}
	query := DB(ctx).Envelope.Query().
		Select(envelope.FieldID, envelope.FieldTo, envelope.FieldFrom, envelope.FieldSubject, envelope.FieldCreatedAt).
		Order(ent.Desc(envelope.FieldID))
	if req.BeforeID > 0 {
		query.Where(envelope.IDLT(req.BeforeID))
	}
	if !admin {
		wheres := []predicate.Envelope{envelope.To(req.To)}
		if !since.IsZero() {
			wheres = append(wheres, envelope.CreatedAtGTE(since))
		}
		query.Where(wheres...)
	}
	query.Limit(limit + 1)
	list, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

type MailDetail struct {
	Content     string             `json:"content"`
	Attachments []AttachmentDetail `json:"attachments"`
}

type AttachmentDetail struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type ReqFetchDetail struct {
	ID int    `param:"id"`
	To string `query:"to"`
}

func FetchDetail(ctx context.Context, req ReqFetchDetail) (*MailDetail, error) {
	if req.To == "" {
		return nil, server.BadParam()
	}
	e, err := DB(ctx).Envelope.Query().
		Select(envelope.FieldContent).
		Where(envelope.ID(req.ID), envelope.To(req.To)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, server.ErrMsg("envelope not found")
	}
	if err != nil {
		return nil, err
	}
	dbAttachments, err := e.QueryAttachments().All(ctx)
	if err != nil {
		return nil, err
	}
	attachments := gox.Map(dbAttachments, func(a *ent.Attachment) AttachmentDetail {
		return AttachmentDetail{
			ID:       a.ID,
			Filename: a.Filename,
		}
	})

	return &MailDetail{
		Content:     e.Content,
		Attachments: attachments,
	}, nil
}

type ReqFetchLatest struct {
	To string `query:"to"`
	ID string `query:"id"`
}

func FetchLatest(ctx context.Context, req ReqFetchLatest) (*server.Reply, error) {
	if req.To == "" {
		return nil, server.BadParam()
	}
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		return nil, server.BadParam()
	}

	to := req.To
	admin := to == Config(ctx).AdminAddress
	wheres := []predicate.Envelope{envelope.IDGT(id)}
	if !admin {
		wheres = append(wheres, envelope.To(to))
	} else {
		to = subAll
	}
	e, err := DB(ctx).Envelope.Query().
		Select(envelope.FieldID, envelope.FieldTo, envelope.FieldFrom, envelope.FieldSubject, envelope.FieldCreatedAt).
		Where(wheres...).
		Order(ent.Asc(envelope.FieldID)).
		First(ctx)
	if err == nil {
		return server.OK(e), nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}

	ch, cancel := notifier.Wait(to)
	defer cancel()
	select {
	case v := <-ch:
		return server.OK(v), nil
	case <-time.After(time.Minute):
		return server.StatusCode(http.StatusNoContent), nil
	case <-ctx.Done():
		return server.Handled(), nil
	}
}

type ReqDownload struct {
	ID string `param:"id"`
	To string `query:"to"`
}

func Download(ctx context.Context, req ReqDownload) (*server.Reply, error) {
	if req.ID == "" || req.To == "" {
		return nil, server.BadParam()
	}

	a, err := DB(ctx).Attachment.Query().
		Where(attachment.ID(req.ID), attachment.HasOwnerWith(envelope.To(req.To))).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, server.ErrMsg("attachment not found")
	}
	if err != nil {
		return nil, err
	}

	ec := server.EchoContext(ctx)
	if err = ec.Attachment(a.Filepath, a.Filename); err != nil {
		return nil, err
	}

	return server.Handled(), nil
}
