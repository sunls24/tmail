// Code generated by ent, DO NOT EDIT.

package envelope

import (
	"time"
	"tmail/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Envelope {
	return predicate.Envelope(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Envelope {
	return predicate.Envelope(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Envelope {
	return predicate.Envelope(sql.FieldLTE(FieldID, id))
}

// To applies equality check predicate on the "to" field. It's identical to ToEQ.
func To(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldTo, v))
}

// From applies equality check predicate on the "from" field. It's identical to FromEQ.
func From(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldFrom, v))
}

// Subject applies equality check predicate on the "subject" field. It's identical to SubjectEQ.
func Subject(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldSubject, v))
}

// Content applies equality check predicate on the "content" field. It's identical to ContentEQ.
func Content(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldContent, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldCreatedAt, v))
}

// ToEQ applies the EQ predicate on the "to" field.
func ToEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldTo, v))
}

// ToNEQ applies the NEQ predicate on the "to" field.
func ToNEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNEQ(FieldTo, v))
}

// ToIn applies the In predicate on the "to" field.
func ToIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldIn(FieldTo, vs...))
}

// ToNotIn applies the NotIn predicate on the "to" field.
func ToNotIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNotIn(FieldTo, vs...))
}

// ToGT applies the GT predicate on the "to" field.
func ToGT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGT(FieldTo, v))
}

// ToGTE applies the GTE predicate on the "to" field.
func ToGTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGTE(FieldTo, v))
}

// ToLT applies the LT predicate on the "to" field.
func ToLT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLT(FieldTo, v))
}

// ToLTE applies the LTE predicate on the "to" field.
func ToLTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLTE(FieldTo, v))
}

// ToContains applies the Contains predicate on the "to" field.
func ToContains(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContains(FieldTo, v))
}

// ToHasPrefix applies the HasPrefix predicate on the "to" field.
func ToHasPrefix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasPrefix(FieldTo, v))
}

// ToHasSuffix applies the HasSuffix predicate on the "to" field.
func ToHasSuffix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasSuffix(FieldTo, v))
}

// ToEqualFold applies the EqualFold predicate on the "to" field.
func ToEqualFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEqualFold(FieldTo, v))
}

// ToContainsFold applies the ContainsFold predicate on the "to" field.
func ToContainsFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContainsFold(FieldTo, v))
}

// FromEQ applies the EQ predicate on the "from" field.
func FromEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldFrom, v))
}

// FromNEQ applies the NEQ predicate on the "from" field.
func FromNEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNEQ(FieldFrom, v))
}

// FromIn applies the In predicate on the "from" field.
func FromIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldIn(FieldFrom, vs...))
}

// FromNotIn applies the NotIn predicate on the "from" field.
func FromNotIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNotIn(FieldFrom, vs...))
}

// FromGT applies the GT predicate on the "from" field.
func FromGT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGT(FieldFrom, v))
}

// FromGTE applies the GTE predicate on the "from" field.
func FromGTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGTE(FieldFrom, v))
}

// FromLT applies the LT predicate on the "from" field.
func FromLT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLT(FieldFrom, v))
}

// FromLTE applies the LTE predicate on the "from" field.
func FromLTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLTE(FieldFrom, v))
}

// FromContains applies the Contains predicate on the "from" field.
func FromContains(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContains(FieldFrom, v))
}

// FromHasPrefix applies the HasPrefix predicate on the "from" field.
func FromHasPrefix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasPrefix(FieldFrom, v))
}

// FromHasSuffix applies the HasSuffix predicate on the "from" field.
func FromHasSuffix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasSuffix(FieldFrom, v))
}

// FromEqualFold applies the EqualFold predicate on the "from" field.
func FromEqualFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEqualFold(FieldFrom, v))
}

// FromContainsFold applies the ContainsFold predicate on the "from" field.
func FromContainsFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContainsFold(FieldFrom, v))
}

// SubjectEQ applies the EQ predicate on the "subject" field.
func SubjectEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldSubject, v))
}

// SubjectNEQ applies the NEQ predicate on the "subject" field.
func SubjectNEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNEQ(FieldSubject, v))
}

// SubjectIn applies the In predicate on the "subject" field.
func SubjectIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldIn(FieldSubject, vs...))
}

// SubjectNotIn applies the NotIn predicate on the "subject" field.
func SubjectNotIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNotIn(FieldSubject, vs...))
}

// SubjectGT applies the GT predicate on the "subject" field.
func SubjectGT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGT(FieldSubject, v))
}

// SubjectGTE applies the GTE predicate on the "subject" field.
func SubjectGTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGTE(FieldSubject, v))
}

// SubjectLT applies the LT predicate on the "subject" field.
func SubjectLT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLT(FieldSubject, v))
}

// SubjectLTE applies the LTE predicate on the "subject" field.
func SubjectLTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLTE(FieldSubject, v))
}

// SubjectContains applies the Contains predicate on the "subject" field.
func SubjectContains(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContains(FieldSubject, v))
}

// SubjectHasPrefix applies the HasPrefix predicate on the "subject" field.
func SubjectHasPrefix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasPrefix(FieldSubject, v))
}

// SubjectHasSuffix applies the HasSuffix predicate on the "subject" field.
func SubjectHasSuffix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasSuffix(FieldSubject, v))
}

// SubjectEqualFold applies the EqualFold predicate on the "subject" field.
func SubjectEqualFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEqualFold(FieldSubject, v))
}

// SubjectContainsFold applies the ContainsFold predicate on the "subject" field.
func SubjectContainsFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContainsFold(FieldSubject, v))
}

// ContentEQ applies the EQ predicate on the "content" field.
func ContentEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldContent, v))
}

// ContentNEQ applies the NEQ predicate on the "content" field.
func ContentNEQ(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNEQ(FieldContent, v))
}

// ContentIn applies the In predicate on the "content" field.
func ContentIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldIn(FieldContent, vs...))
}

// ContentNotIn applies the NotIn predicate on the "content" field.
func ContentNotIn(vs ...string) predicate.Envelope {
	return predicate.Envelope(sql.FieldNotIn(FieldContent, vs...))
}

// ContentGT applies the GT predicate on the "content" field.
func ContentGT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGT(FieldContent, v))
}

// ContentGTE applies the GTE predicate on the "content" field.
func ContentGTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldGTE(FieldContent, v))
}

// ContentLT applies the LT predicate on the "content" field.
func ContentLT(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLT(FieldContent, v))
}

// ContentLTE applies the LTE predicate on the "content" field.
func ContentLTE(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldLTE(FieldContent, v))
}

// ContentContains applies the Contains predicate on the "content" field.
func ContentContains(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContains(FieldContent, v))
}

// ContentHasPrefix applies the HasPrefix predicate on the "content" field.
func ContentHasPrefix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasPrefix(FieldContent, v))
}

// ContentHasSuffix applies the HasSuffix predicate on the "content" field.
func ContentHasSuffix(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldHasSuffix(FieldContent, v))
}

// ContentEqualFold applies the EqualFold predicate on the "content" field.
func ContentEqualFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldEqualFold(FieldContent, v))
}

// ContentContainsFold applies the ContainsFold predicate on the "content" field.
func ContentContainsFold(v string) predicate.Envelope {
	return predicate.Envelope(sql.FieldContainsFold(FieldContent, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Envelope {
	return predicate.Envelope(sql.FieldLTE(FieldCreatedAt, v))
}

// HasAttachments applies the HasEdge predicate on the "attachments" edge.
func HasAttachments() predicate.Envelope {
	return predicate.Envelope(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AttachmentsTable, AttachmentsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAttachmentsWith applies the HasEdge predicate on the "attachments" edge with a given conditions (other predicates).
func HasAttachmentsWith(preds ...predicate.Attachment) predicate.Envelope {
	return predicate.Envelope(func(s *sql.Selector) {
		step := newAttachmentsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Envelope) predicate.Envelope {
	return predicate.Envelope(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Envelope) predicate.Envelope {
	return predicate.Envelope(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Envelope) predicate.Envelope {
	return predicate.Envelope(sql.NotPredicates(p))
}
