// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tmail/ent/envelope"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// EnvelopeCreate is the builder for creating a Envelope entity.
type EnvelopeCreate struct {
	config
	mutation *EnvelopeMutation
	hooks    []Hook
}

// SetTo sets the "to" field.
func (ec *EnvelopeCreate) SetTo(s string) *EnvelopeCreate {
	ec.mutation.SetTo(s)
	return ec
}

// SetFrom sets the "from" field.
func (ec *EnvelopeCreate) SetFrom(s string) *EnvelopeCreate {
	ec.mutation.SetFrom(s)
	return ec
}

// SetSubject sets the "subject" field.
func (ec *EnvelopeCreate) SetSubject(s string) *EnvelopeCreate {
	ec.mutation.SetSubject(s)
	return ec
}

// SetContent sets the "content" field.
func (ec *EnvelopeCreate) SetContent(s string) *EnvelopeCreate {
	ec.mutation.SetContent(s)
	return ec
}

// SetCreatedAt sets the "created_at" field.
func (ec *EnvelopeCreate) SetCreatedAt(t time.Time) *EnvelopeCreate {
	ec.mutation.SetCreatedAt(t)
	return ec
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (ec *EnvelopeCreate) SetNillableCreatedAt(t *time.Time) *EnvelopeCreate {
	if t != nil {
		ec.SetCreatedAt(*t)
	}
	return ec
}

// Mutation returns the EnvelopeMutation object of the builder.
func (ec *EnvelopeCreate) Mutation() *EnvelopeMutation {
	return ec.mutation
}

// Save creates the Envelope in the database.
func (ec *EnvelopeCreate) Save(ctx context.Context) (*Envelope, error) {
	ec.defaults()
	return withHooks(ctx, ec.sqlSave, ec.mutation, ec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EnvelopeCreate) SaveX(ctx context.Context) *Envelope {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *EnvelopeCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *EnvelopeCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ec *EnvelopeCreate) defaults() {
	if _, ok := ec.mutation.CreatedAt(); !ok {
		v := envelope.DefaultCreatedAt()
		ec.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ec *EnvelopeCreate) check() error {
	if _, ok := ec.mutation.To(); !ok {
		return &ValidationError{Name: "to", err: errors.New(`ent: missing required field "Envelope.to"`)}
	}
	if v, ok := ec.mutation.To(); ok {
		if err := envelope.ToValidator(v); err != nil {
			return &ValidationError{Name: "to", err: fmt.Errorf(`ent: validator failed for field "Envelope.to": %w`, err)}
		}
	}
	if _, ok := ec.mutation.From(); !ok {
		return &ValidationError{Name: "from", err: errors.New(`ent: missing required field "Envelope.from"`)}
	}
	if v, ok := ec.mutation.From(); ok {
		if err := envelope.FromValidator(v); err != nil {
			return &ValidationError{Name: "from", err: fmt.Errorf(`ent: validator failed for field "Envelope.from": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Subject(); !ok {
		return &ValidationError{Name: "subject", err: errors.New(`ent: missing required field "Envelope.subject"`)}
	}
	if v, ok := ec.mutation.Subject(); ok {
		if err := envelope.SubjectValidator(v); err != nil {
			return &ValidationError{Name: "subject", err: fmt.Errorf(`ent: validator failed for field "Envelope.subject": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "Envelope.content"`)}
	}
	if v, ok := ec.mutation.Content(); ok {
		if err := envelope.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "Envelope.content": %w`, err)}
		}
	}
	if _, ok := ec.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Envelope.created_at"`)}
	}
	return nil
}

func (ec *EnvelopeCreate) sqlSave(ctx context.Context) (*Envelope, error) {
	if err := ec.check(); err != nil {
		return nil, err
	}
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ec.mutation.id = &_node.ID
	ec.mutation.done = true
	return _node, nil
}

func (ec *EnvelopeCreate) createSpec() (*Envelope, *sqlgraph.CreateSpec) {
	var (
		_node = &Envelope{config: ec.config}
		_spec = sqlgraph.NewCreateSpec(envelope.Table, sqlgraph.NewFieldSpec(envelope.FieldID, field.TypeInt))
	)
	if value, ok := ec.mutation.To(); ok {
		_spec.SetField(envelope.FieldTo, field.TypeString, value)
		_node.To = value
	}
	if value, ok := ec.mutation.From(); ok {
		_spec.SetField(envelope.FieldFrom, field.TypeString, value)
		_node.From = value
	}
	if value, ok := ec.mutation.Subject(); ok {
		_spec.SetField(envelope.FieldSubject, field.TypeString, value)
		_node.Subject = value
	}
	if value, ok := ec.mutation.Content(); ok {
		_spec.SetField(envelope.FieldContent, field.TypeString, value)
		_node.Content = value
	}
	if value, ok := ec.mutation.CreatedAt(); ok {
		_spec.SetField(envelope.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	return _node, _spec
}

// EnvelopeCreateBulk is the builder for creating many Envelope entities in bulk.
type EnvelopeCreateBulk struct {
	config
	err      error
	builders []*EnvelopeCreate
}

// Save creates the Envelope entities in the database.
func (ecb *EnvelopeCreateBulk) Save(ctx context.Context) ([]*Envelope, error) {
	if ecb.err != nil {
		return nil, ecb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Envelope, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnvelopeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *EnvelopeCreateBulk) SaveX(ctx context.Context) []*Envelope {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *EnvelopeCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *EnvelopeCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}
