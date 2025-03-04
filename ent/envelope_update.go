// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tmail/ent/envelope"
	"tmail/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// EnvelopeUpdate is the builder for updating Envelope entities.
type EnvelopeUpdate struct {
	config
	hooks    []Hook
	mutation *EnvelopeMutation
}

// Where appends a list predicates to the EnvelopeUpdate builder.
func (eu *EnvelopeUpdate) Where(ps ...predicate.Envelope) *EnvelopeUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetTo sets the "to" field.
func (eu *EnvelopeUpdate) SetTo(s string) *EnvelopeUpdate {
	eu.mutation.SetTo(s)
	return eu
}

// SetNillableTo sets the "to" field if the given value is not nil.
func (eu *EnvelopeUpdate) SetNillableTo(s *string) *EnvelopeUpdate {
	if s != nil {
		eu.SetTo(*s)
	}
	return eu
}

// SetFrom sets the "from" field.
func (eu *EnvelopeUpdate) SetFrom(s string) *EnvelopeUpdate {
	eu.mutation.SetFrom(s)
	return eu
}

// SetNillableFrom sets the "from" field if the given value is not nil.
func (eu *EnvelopeUpdate) SetNillableFrom(s *string) *EnvelopeUpdate {
	if s != nil {
		eu.SetFrom(*s)
	}
	return eu
}

// SetSubject sets the "subject" field.
func (eu *EnvelopeUpdate) SetSubject(s string) *EnvelopeUpdate {
	eu.mutation.SetSubject(s)
	return eu
}

// SetNillableSubject sets the "subject" field if the given value is not nil.
func (eu *EnvelopeUpdate) SetNillableSubject(s *string) *EnvelopeUpdate {
	if s != nil {
		eu.SetSubject(*s)
	}
	return eu
}

// SetContent sets the "content" field.
func (eu *EnvelopeUpdate) SetContent(s string) *EnvelopeUpdate {
	eu.mutation.SetContent(s)
	return eu
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (eu *EnvelopeUpdate) SetNillableContent(s *string) *EnvelopeUpdate {
	if s != nil {
		eu.SetContent(*s)
	}
	return eu
}

// SetCreatedAt sets the "created_at" field.
func (eu *EnvelopeUpdate) SetCreatedAt(t time.Time) *EnvelopeUpdate {
	eu.mutation.SetCreatedAt(t)
	return eu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (eu *EnvelopeUpdate) SetNillableCreatedAt(t *time.Time) *EnvelopeUpdate {
	if t != nil {
		eu.SetCreatedAt(*t)
	}
	return eu
}

// Mutation returns the EnvelopeMutation object of the builder.
func (eu *EnvelopeUpdate) Mutation() *EnvelopeMutation {
	return eu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EnvelopeUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EnvelopeUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EnvelopeUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EnvelopeUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eu *EnvelopeUpdate) check() error {
	if v, ok := eu.mutation.To(); ok {
		if err := envelope.ToValidator(v); err != nil {
			return &ValidationError{Name: "to", err: fmt.Errorf(`ent: validator failed for field "Envelope.to": %w`, err)}
		}
	}
	if v, ok := eu.mutation.From(); ok {
		if err := envelope.FromValidator(v); err != nil {
			return &ValidationError{Name: "from", err: fmt.Errorf(`ent: validator failed for field "Envelope.from": %w`, err)}
		}
	}
	if v, ok := eu.mutation.Subject(); ok {
		if err := envelope.SubjectValidator(v); err != nil {
			return &ValidationError{Name: "subject", err: fmt.Errorf(`ent: validator failed for field "Envelope.subject": %w`, err)}
		}
	}
	if v, ok := eu.mutation.Content(); ok {
		if err := envelope.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "Envelope.content": %w`, err)}
		}
	}
	return nil
}

func (eu *EnvelopeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := eu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(envelope.Table, envelope.Columns, sqlgraph.NewFieldSpec(envelope.FieldID, field.TypeInt))
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.To(); ok {
		_spec.SetField(envelope.FieldTo, field.TypeString, value)
	}
	if value, ok := eu.mutation.From(); ok {
		_spec.SetField(envelope.FieldFrom, field.TypeString, value)
	}
	if value, ok := eu.mutation.Subject(); ok {
		_spec.SetField(envelope.FieldSubject, field.TypeString, value)
	}
	if value, ok := eu.mutation.Content(); ok {
		_spec.SetField(envelope.FieldContent, field.TypeString, value)
	}
	if value, ok := eu.mutation.CreatedAt(); ok {
		_spec.SetField(envelope.FieldCreatedAt, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{envelope.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EnvelopeUpdateOne is the builder for updating a single Envelope entity.
type EnvelopeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EnvelopeMutation
}

// SetTo sets the "to" field.
func (euo *EnvelopeUpdateOne) SetTo(s string) *EnvelopeUpdateOne {
	euo.mutation.SetTo(s)
	return euo
}

// SetNillableTo sets the "to" field if the given value is not nil.
func (euo *EnvelopeUpdateOne) SetNillableTo(s *string) *EnvelopeUpdateOne {
	if s != nil {
		euo.SetTo(*s)
	}
	return euo
}

// SetFrom sets the "from" field.
func (euo *EnvelopeUpdateOne) SetFrom(s string) *EnvelopeUpdateOne {
	euo.mutation.SetFrom(s)
	return euo
}

// SetNillableFrom sets the "from" field if the given value is not nil.
func (euo *EnvelopeUpdateOne) SetNillableFrom(s *string) *EnvelopeUpdateOne {
	if s != nil {
		euo.SetFrom(*s)
	}
	return euo
}

// SetSubject sets the "subject" field.
func (euo *EnvelopeUpdateOne) SetSubject(s string) *EnvelopeUpdateOne {
	euo.mutation.SetSubject(s)
	return euo
}

// SetNillableSubject sets the "subject" field if the given value is not nil.
func (euo *EnvelopeUpdateOne) SetNillableSubject(s *string) *EnvelopeUpdateOne {
	if s != nil {
		euo.SetSubject(*s)
	}
	return euo
}

// SetContent sets the "content" field.
func (euo *EnvelopeUpdateOne) SetContent(s string) *EnvelopeUpdateOne {
	euo.mutation.SetContent(s)
	return euo
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (euo *EnvelopeUpdateOne) SetNillableContent(s *string) *EnvelopeUpdateOne {
	if s != nil {
		euo.SetContent(*s)
	}
	return euo
}

// SetCreatedAt sets the "created_at" field.
func (euo *EnvelopeUpdateOne) SetCreatedAt(t time.Time) *EnvelopeUpdateOne {
	euo.mutation.SetCreatedAt(t)
	return euo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (euo *EnvelopeUpdateOne) SetNillableCreatedAt(t *time.Time) *EnvelopeUpdateOne {
	if t != nil {
		euo.SetCreatedAt(*t)
	}
	return euo
}

// Mutation returns the EnvelopeMutation object of the builder.
func (euo *EnvelopeUpdateOne) Mutation() *EnvelopeMutation {
	return euo.mutation
}

// Where appends a list predicates to the EnvelopeUpdate builder.
func (euo *EnvelopeUpdateOne) Where(ps ...predicate.Envelope) *EnvelopeUpdateOne {
	euo.mutation.Where(ps...)
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EnvelopeUpdateOne) Select(field string, fields ...string) *EnvelopeUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Envelope entity.
func (euo *EnvelopeUpdateOne) Save(ctx context.Context) (*Envelope, error) {
	return withHooks(ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EnvelopeUpdateOne) SaveX(ctx context.Context) *Envelope {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EnvelopeUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EnvelopeUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (euo *EnvelopeUpdateOne) check() error {
	if v, ok := euo.mutation.To(); ok {
		if err := envelope.ToValidator(v); err != nil {
			return &ValidationError{Name: "to", err: fmt.Errorf(`ent: validator failed for field "Envelope.to": %w`, err)}
		}
	}
	if v, ok := euo.mutation.From(); ok {
		if err := envelope.FromValidator(v); err != nil {
			return &ValidationError{Name: "from", err: fmt.Errorf(`ent: validator failed for field "Envelope.from": %w`, err)}
		}
	}
	if v, ok := euo.mutation.Subject(); ok {
		if err := envelope.SubjectValidator(v); err != nil {
			return &ValidationError{Name: "subject", err: fmt.Errorf(`ent: validator failed for field "Envelope.subject": %w`, err)}
		}
	}
	if v, ok := euo.mutation.Content(); ok {
		if err := envelope.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "Envelope.content": %w`, err)}
		}
	}
	return nil
}

func (euo *EnvelopeUpdateOne) sqlSave(ctx context.Context) (_node *Envelope, err error) {
	if err := euo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(envelope.Table, envelope.Columns, sqlgraph.NewFieldSpec(envelope.FieldID, field.TypeInt))
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Envelope.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, envelope.FieldID)
		for _, f := range fields {
			if !envelope.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != envelope.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.To(); ok {
		_spec.SetField(envelope.FieldTo, field.TypeString, value)
	}
	if value, ok := euo.mutation.From(); ok {
		_spec.SetField(envelope.FieldFrom, field.TypeString, value)
	}
	if value, ok := euo.mutation.Subject(); ok {
		_spec.SetField(envelope.FieldSubject, field.TypeString, value)
	}
	if value, ok := euo.mutation.Content(); ok {
		_spec.SetField(envelope.FieldContent, field.TypeString, value)
	}
	if value, ok := euo.mutation.CreatedAt(); ok {
		_spec.SetField(envelope.FieldCreatedAt, field.TypeTime, value)
	}
	_node = &Envelope{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{envelope.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
