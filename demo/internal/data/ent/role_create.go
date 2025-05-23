// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"github/invokerw/gintos/demo/internal/data/ent/role"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// RoleCreate is the builder for creating a Role entity.
type RoleCreate struct {
	config
	mutation *RoleMutation
	hooks    []Hook
}

// SetCreateTime sets the "create_time" field.
func (rc *RoleCreate) SetCreateTime(t time.Time) *RoleCreate {
	rc.mutation.SetCreateTime(t)
	return rc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (rc *RoleCreate) SetNillableCreateTime(t *time.Time) *RoleCreate {
	if t != nil {
		rc.SetCreateTime(*t)
	}
	return rc
}

// SetUpdateTime sets the "update_time" field.
func (rc *RoleCreate) SetUpdateTime(t time.Time) *RoleCreate {
	rc.mutation.SetUpdateTime(t)
	return rc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (rc *RoleCreate) SetNillableUpdateTime(t *time.Time) *RoleCreate {
	if t != nil {
		rc.SetUpdateTime(*t)
	}
	return rc
}

// SetStatus sets the "status" field.
func (rc *RoleCreate) SetStatus(r role.Status) *RoleCreate {
	rc.mutation.SetStatus(r)
	return rc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (rc *RoleCreate) SetNillableStatus(r *role.Status) *RoleCreate {
	if r != nil {
		rc.SetStatus(*r)
	}
	return rc
}

// SetCreateBy sets the "create_by" field.
func (rc *RoleCreate) SetCreateBy(u uint64) *RoleCreate {
	rc.mutation.SetCreateBy(u)
	return rc
}

// SetNillableCreateBy sets the "create_by" field if the given value is not nil.
func (rc *RoleCreate) SetNillableCreateBy(u *uint64) *RoleCreate {
	if u != nil {
		rc.SetCreateBy(*u)
	}
	return rc
}

// SetUpdateBy sets the "update_by" field.
func (rc *RoleCreate) SetUpdateBy(u uint64) *RoleCreate {
	rc.mutation.SetUpdateBy(u)
	return rc
}

// SetNillableUpdateBy sets the "update_by" field if the given value is not nil.
func (rc *RoleCreate) SetNillableUpdateBy(u *uint64) *RoleCreate {
	if u != nil {
		rc.SetUpdateBy(*u)
	}
	return rc
}

// SetRemark sets the "remark" field.
func (rc *RoleCreate) SetRemark(s string) *RoleCreate {
	rc.mutation.SetRemark(s)
	return rc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (rc *RoleCreate) SetNillableRemark(s *string) *RoleCreate {
	if s != nil {
		rc.SetRemark(*s)
	}
	return rc
}

// SetName sets the "name" field.
func (rc *RoleCreate) SetName(s string) *RoleCreate {
	rc.mutation.SetName(s)
	return rc
}

// SetCode sets the "code" field.
func (rc *RoleCreate) SetCode(s string) *RoleCreate {
	rc.mutation.SetCode(s)
	return rc
}

// SetSortID sets the "sort_id" field.
func (rc *RoleCreate) SetSortID(i int32) *RoleCreate {
	rc.mutation.SetSortID(i)
	return rc
}

// SetNillableSortID sets the "sort_id" field if the given value is not nil.
func (rc *RoleCreate) SetNillableSortID(i *int32) *RoleCreate {
	if i != nil {
		rc.SetSortID(*i)
	}
	return rc
}

// SetID sets the "id" field.
func (rc *RoleCreate) SetID(u uint64) *RoleCreate {
	rc.mutation.SetID(u)
	return rc
}

// Mutation returns the RoleMutation object of the builder.
func (rc *RoleCreate) Mutation() *RoleMutation {
	return rc.mutation
}

// Save creates the Role in the database.
func (rc *RoleCreate) Save(ctx context.Context) (*Role, error) {
	rc.defaults()
	return withHooks(ctx, rc.sqlSave, rc.mutation, rc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RoleCreate) SaveX(ctx context.Context) *Role {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *RoleCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *RoleCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *RoleCreate) defaults() {
	if _, ok := rc.mutation.Status(); !ok {
		v := role.DefaultStatus
		rc.mutation.SetStatus(v)
	}
	if _, ok := rc.mutation.Remark(); !ok {
		v := role.DefaultRemark
		rc.mutation.SetRemark(v)
	}
	if _, ok := rc.mutation.SortID(); !ok {
		v := role.DefaultSortID
		rc.mutation.SetSortID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RoleCreate) check() error {
	if v, ok := rc.mutation.Status(); ok {
		if err := role.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Role.status": %w`, err)}
		}
	}
	if _, ok := rc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Role.name"`)}
	}
	if v, ok := rc.mutation.Name(); ok {
		if err := role.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Role.name": %w`, err)}
		}
	}
	if _, ok := rc.mutation.Code(); !ok {
		return &ValidationError{Name: "code", err: errors.New(`ent: missing required field "Role.code"`)}
	}
	if v, ok := rc.mutation.Code(); ok {
		if err := role.CodeValidator(v); err != nil {
			return &ValidationError{Name: "code", err: fmt.Errorf(`ent: validator failed for field "Role.code": %w`, err)}
		}
	}
	if v, ok := rc.mutation.ID(); ok {
		if err := role.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Role.id": %w`, err)}
		}
	}
	return nil
}

func (rc *RoleCreate) sqlSave(ctx context.Context) (*Role, error) {
	if err := rc.check(); err != nil {
		return nil, err
	}
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = uint64(id)
	}
	rc.mutation.id = &_node.ID
	rc.mutation.done = true
	return _node, nil
}

func (rc *RoleCreate) createSpec() (*Role, *sqlgraph.CreateSpec) {
	var (
		_node = &Role{config: rc.config}
		_spec = sqlgraph.NewCreateSpec(role.Table, sqlgraph.NewFieldSpec(role.FieldID, field.TypeUint64))
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := rc.mutation.CreateTime(); ok {
		_spec.SetField(role.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := rc.mutation.UpdateTime(); ok {
		_spec.SetField(role.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := rc.mutation.Status(); ok {
		_spec.SetField(role.FieldStatus, field.TypeEnum, value)
		_node.Status = &value
	}
	if value, ok := rc.mutation.CreateBy(); ok {
		_spec.SetField(role.FieldCreateBy, field.TypeUint64, value)
		_node.CreateBy = &value
	}
	if value, ok := rc.mutation.UpdateBy(); ok {
		_spec.SetField(role.FieldUpdateBy, field.TypeUint64, value)
		_node.UpdateBy = &value
	}
	if value, ok := rc.mutation.Remark(); ok {
		_spec.SetField(role.FieldRemark, field.TypeString, value)
		_node.Remark = &value
	}
	if value, ok := rc.mutation.Name(); ok {
		_spec.SetField(role.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := rc.mutation.Code(); ok {
		_spec.SetField(role.FieldCode, field.TypeString, value)
		_node.Code = value
	}
	if value, ok := rc.mutation.SortID(); ok {
		_spec.SetField(role.FieldSortID, field.TypeInt32, value)
		_node.SortID = &value
	}
	return _node, _spec
}

// RoleCreateBulk is the builder for creating many Role entities in bulk.
type RoleCreateBulk struct {
	config
	err      error
	builders []*RoleCreate
}

// Save creates the Role entities in the database.
func (rcb *RoleCreateBulk) Save(ctx context.Context) ([]*Role, error) {
	if rcb.err != nil {
		return nil, rcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Role, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RoleMutation)
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
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint64(id)
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
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RoleCreateBulk) SaveX(ctx context.Context) []*Role {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *RoleCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *RoleCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}
