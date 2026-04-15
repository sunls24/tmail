package api

import (
	"net/http"
	"strconv"
	"time"
	"tmail/ent"
	"tmail/ent/attachment"
	"tmail/ent/envelope"
	"tmail/ent/predicate"
	"tmail/internal/pubsub"
)

func Fetch(ctx *Context) error {
	to := ctx.QueryParam("to")
	if to == "" {
		return ctx.Bad("not found to address")
	}
	admin := to == ctx.AdminAddress
	since := time.Time{}
	if !admin {
		sinceStr := ctx.QueryParam("since")
		if sinceStr != "" {
			ts, err := strconv.ParseInt(sinceStr, 10, 64)
			if err == nil && ts >= 0 {
				since = time.Unix(ts, 0)
			}
		}
	}
	query := ctx.ent.Envelope.Query().
		Select(envelope.FieldID, envelope.FieldTo, envelope.FieldFrom, envelope.FieldSubject, envelope.FieldCreatedAt).
		Order(ent.Desc(envelope.FieldID))
	if !admin {
		wheres := []predicate.Envelope{envelope.To(to)}
		if !since.IsZero() {
			wheres = append(wheres, envelope.CreatedAtGTE(since))
		}
		query.Where(wheres...)
	} else {
		query.Limit(100)
	}
	list, err := query.All(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, list)
}

type MailDetail struct {
	Content     string             `json:"content"`
	Attachments []AttachmentDetail `json:"attachments"`
}

type AttachmentDetail struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

func FetchDetail(ctx *Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.Bad("not found id param")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Bad("invalid id param: " + idStr)
	}
	e, err := ctx.ent.Envelope.Query().
		Select(envelope.FieldContent).
		Where(envelope.ID(id)).
		Only(ctx.Request().Context())
	if ent.IsNotFound(err) {
		return ctx.Badf("envelope %d not found", id)
	}
	if err != nil {
		return err
	}
	dbAttachments, _ := e.QueryAttachments().All(ctx.Request().Context())
	attachments := make([]AttachmentDetail, 0, len(dbAttachments))
	for _, a := range dbAttachments {
		attachments = append(attachments, AttachmentDetail{
			ID:       a.ID,
			Filename: a.Filename,
		})
	}

	return ctx.JSON(http.StatusOK, MailDetail{
		Content:     e.Content,
		Attachments: attachments,
	})
}

func FetchLatest(ctx *Context) error {
	to := ctx.QueryParam("to")
	if to == "" {
		return ctx.Bad("not found to address")
	}
	admin := to == ctx.AdminAddress
	if !admin {
		idStr := ctx.QueryParam("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return ctx.Bad("invalid id param: " + idStr)
		}
		e, err := ctx.ent.Envelope.Query().
			Select(envelope.FieldID, envelope.FieldTo, envelope.FieldFrom, envelope.FieldSubject, envelope.FieldCreatedAt).
			Where(envelope.IDGT(id), envelope.To(to)).
			Order(ent.Asc(envelope.FieldID)).
			First(ctx.Request().Context())
		if err == nil {
			return ctx.JSON(http.StatusOK, e)
		}
		if !ent.IsNotFound(err) {
			return err
		}
	} else {
		to = pubsub.SubAll
	}

	ch, cancel := pubsub.Subscribe(to)
	defer cancel()
	select {
	case e := <-ch:
		return ctx.JSON(http.StatusOK, e)
	case <-time.After(time.Minute):
		return ctx.NoContent(http.StatusNoContent)
	case <-ctx.Request().Context().Done():
		return nil
	}
}

func Download(ctx *Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.Bad("not found id param")
	}

	a, err := ctx.ent.Attachment.Query().Where(attachment.ID(id)).First(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.Attachment(a.Filepath, a.Filename)
}
