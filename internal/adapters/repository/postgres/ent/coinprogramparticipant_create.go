// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogramparticipant"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CoinProgramParticipantCreate is the builder for creating a CoinProgramParticipant entity.
type CoinProgramParticipantCreate struct {
	config
	mutation *CoinProgramParticipantMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (cppc *CoinProgramParticipantCreate) SetCreatedAt(t time.Time) *CoinProgramParticipantCreate {
	cppc.mutation.SetCreatedAt(t)
	return cppc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cppc *CoinProgramParticipantCreate) SetNillableCreatedAt(t *time.Time) *CoinProgramParticipantCreate {
	if t != nil {
		cppc.SetCreatedAt(*t)
	}
	return cppc
}

// SetBalance sets the "balance" field.
func (cppc *CoinProgramParticipantCreate) SetBalance(u uint) *CoinProgramParticipantCreate {
	cppc.mutation.SetBalance(u)
	return cppc
}

// SetNillableBalance sets the "balance" field if the given value is not nil.
func (cppc *CoinProgramParticipantCreate) SetNillableBalance(u *uint) *CoinProgramParticipantCreate {
	if u != nil {
		cppc.SetBalance(*u)
	}
	return cppc
}

// SetID sets the "id" field.
func (cppc *CoinProgramParticipantCreate) SetID(u uuid.UUID) *CoinProgramParticipantCreate {
	cppc.mutation.SetID(u)
	return cppc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (cppc *CoinProgramParticipantCreate) SetNillableID(u *uuid.UUID) *CoinProgramParticipantCreate {
	if u != nil {
		cppc.SetID(*u)
	}
	return cppc
}

// SetCoinProgramID sets the "coin_program" edge to the CoinProgram entity by ID.
func (cppc *CoinProgramParticipantCreate) SetCoinProgramID(id uuid.UUID) *CoinProgramParticipantCreate {
	cppc.mutation.SetCoinProgramID(id)
	return cppc
}

// SetNillableCoinProgramID sets the "coin_program" edge to the CoinProgram entity by ID if the given value is not nil.
func (cppc *CoinProgramParticipantCreate) SetNillableCoinProgramID(id *uuid.UUID) *CoinProgramParticipantCreate {
	if id != nil {
		cppc = cppc.SetCoinProgramID(*id)
	}
	return cppc
}

// SetCoinProgram sets the "coin_program" edge to the CoinProgram entity.
func (cppc *CoinProgramParticipantCreate) SetCoinProgram(c *CoinProgram) *CoinProgramParticipantCreate {
	return cppc.SetCoinProgramID(c.ID)
}

// SetUserID sets the "user" edge to the User entity by ID.
func (cppc *CoinProgramParticipantCreate) SetUserID(id uuid.UUID) *CoinProgramParticipantCreate {
	cppc.mutation.SetUserID(id)
	return cppc
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (cppc *CoinProgramParticipantCreate) SetNillableUserID(id *uuid.UUID) *CoinProgramParticipantCreate {
	if id != nil {
		cppc = cppc.SetUserID(*id)
	}
	return cppc
}

// SetUser sets the "user" edge to the User entity.
func (cppc *CoinProgramParticipantCreate) SetUser(u *User) *CoinProgramParticipantCreate {
	return cppc.SetUserID(u.ID)
}

// Mutation returns the CoinProgramParticipantMutation object of the builder.
func (cppc *CoinProgramParticipantCreate) Mutation() *CoinProgramParticipantMutation {
	return cppc.mutation
}

// Save creates the CoinProgramParticipant in the database.
func (cppc *CoinProgramParticipantCreate) Save(ctx context.Context) (*CoinProgramParticipant, error) {
	cppc.defaults()
	return withHooks(ctx, cppc.sqlSave, cppc.mutation, cppc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cppc *CoinProgramParticipantCreate) SaveX(ctx context.Context) *CoinProgramParticipant {
	v, err := cppc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cppc *CoinProgramParticipantCreate) Exec(ctx context.Context) error {
	_, err := cppc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cppc *CoinProgramParticipantCreate) ExecX(ctx context.Context) {
	if err := cppc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cppc *CoinProgramParticipantCreate) defaults() {
	if _, ok := cppc.mutation.CreatedAt(); !ok {
		v := coinprogramparticipant.DefaultCreatedAt
		cppc.mutation.SetCreatedAt(v)
	}
	if _, ok := cppc.mutation.Balance(); !ok {
		v := coinprogramparticipant.DefaultBalance
		cppc.mutation.SetBalance(v)
	}
	if _, ok := cppc.mutation.ID(); !ok {
		v := coinprogramparticipant.DefaultID()
		cppc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cppc *CoinProgramParticipantCreate) check() error {
	if _, ok := cppc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "CoinProgramParticipant.created_at"`)}
	}
	if _, ok := cppc.mutation.Balance(); !ok {
		return &ValidationError{Name: "balance", err: errors.New(`ent: missing required field "CoinProgramParticipant.balance"`)}
	}
	return nil
}

func (cppc *CoinProgramParticipantCreate) sqlSave(ctx context.Context) (*CoinProgramParticipant, error) {
	if err := cppc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cppc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cppc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cppc.mutation.id = &_node.ID
	cppc.mutation.done = true
	return _node, nil
}

func (cppc *CoinProgramParticipantCreate) createSpec() (*CoinProgramParticipant, *sqlgraph.CreateSpec) {
	var (
		_node = &CoinProgramParticipant{config: cppc.config}
		_spec = sqlgraph.NewCreateSpec(coinprogramparticipant.Table, sqlgraph.NewFieldSpec(coinprogramparticipant.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = cppc.conflict
	if id, ok := cppc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cppc.mutation.CreatedAt(); ok {
		_spec.SetField(coinprogramparticipant.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := cppc.mutation.Balance(); ok {
		_spec.SetField(coinprogramparticipant.FieldBalance, field.TypeUint, value)
		_node.Balance = value
	}
	if nodes := cppc.mutation.CoinProgramIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   coinprogramparticipant.CoinProgramTable,
			Columns: []string{coinprogramparticipant.CoinProgramColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(coinprogram.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.coin_program_participant_coin_program = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cppc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   coinprogramparticipant.UserTable,
			Columns: []string{coinprogramparticipant.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_coin_programs = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.CoinProgramParticipant.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.CoinProgramParticipantUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (cppc *CoinProgramParticipantCreate) OnConflict(opts ...sql.ConflictOption) *CoinProgramParticipantUpsertOne {
	cppc.conflict = opts
	return &CoinProgramParticipantUpsertOne{
		create: cppc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.CoinProgramParticipant.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cppc *CoinProgramParticipantCreate) OnConflictColumns(columns ...string) *CoinProgramParticipantUpsertOne {
	cppc.conflict = append(cppc.conflict, sql.ConflictColumns(columns...))
	return &CoinProgramParticipantUpsertOne{
		create: cppc,
	}
}

type (
	// CoinProgramParticipantUpsertOne is the builder for "upsert"-ing
	//  one CoinProgramParticipant node.
	CoinProgramParticipantUpsertOne struct {
		create *CoinProgramParticipantCreate
	}

	// CoinProgramParticipantUpsert is the "OnConflict" setter.
	CoinProgramParticipantUpsert struct {
		*sql.UpdateSet
	}
)

// SetBalance sets the "balance" field.
func (u *CoinProgramParticipantUpsert) SetBalance(v uint) *CoinProgramParticipantUpsert {
	u.Set(coinprogramparticipant.FieldBalance, v)
	return u
}

// UpdateBalance sets the "balance" field to the value that was provided on create.
func (u *CoinProgramParticipantUpsert) UpdateBalance() *CoinProgramParticipantUpsert {
	u.SetExcluded(coinprogramparticipant.FieldBalance)
	return u
}

// AddBalance adds v to the "balance" field.
func (u *CoinProgramParticipantUpsert) AddBalance(v uint) *CoinProgramParticipantUpsert {
	u.Add(coinprogramparticipant.FieldBalance, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.CoinProgramParticipant.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(coinprogramparticipant.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *CoinProgramParticipantUpsertOne) UpdateNewValues() *CoinProgramParticipantUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(coinprogramparticipant.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(coinprogramparticipant.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.CoinProgramParticipant.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *CoinProgramParticipantUpsertOne) Ignore() *CoinProgramParticipantUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *CoinProgramParticipantUpsertOne) DoNothing() *CoinProgramParticipantUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the CoinProgramParticipantCreate.OnConflict
// documentation for more info.
func (u *CoinProgramParticipantUpsertOne) Update(set func(*CoinProgramParticipantUpsert)) *CoinProgramParticipantUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&CoinProgramParticipantUpsert{UpdateSet: update})
	}))
	return u
}

// SetBalance sets the "balance" field.
func (u *CoinProgramParticipantUpsertOne) SetBalance(v uint) *CoinProgramParticipantUpsertOne {
	return u.Update(func(s *CoinProgramParticipantUpsert) {
		s.SetBalance(v)
	})
}

// AddBalance adds v to the "balance" field.
func (u *CoinProgramParticipantUpsertOne) AddBalance(v uint) *CoinProgramParticipantUpsertOne {
	return u.Update(func(s *CoinProgramParticipantUpsert) {
		s.AddBalance(v)
	})
}

// UpdateBalance sets the "balance" field to the value that was provided on create.
func (u *CoinProgramParticipantUpsertOne) UpdateBalance() *CoinProgramParticipantUpsertOne {
	return u.Update(func(s *CoinProgramParticipantUpsert) {
		s.UpdateBalance()
	})
}

// Exec executes the query.
func (u *CoinProgramParticipantUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for CoinProgramParticipantCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *CoinProgramParticipantUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *CoinProgramParticipantUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: CoinProgramParticipantUpsertOne.ID is not supported by MySQL driver. Use CoinProgramParticipantUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *CoinProgramParticipantUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// CoinProgramParticipantCreateBulk is the builder for creating many CoinProgramParticipant entities in bulk.
type CoinProgramParticipantCreateBulk struct {
	config
	err      error
	builders []*CoinProgramParticipantCreate
	conflict []sql.ConflictOption
}

// Save creates the CoinProgramParticipant entities in the database.
func (cppcb *CoinProgramParticipantCreateBulk) Save(ctx context.Context) ([]*CoinProgramParticipant, error) {
	if cppcb.err != nil {
		return nil, cppcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(cppcb.builders))
	nodes := make([]*CoinProgramParticipant, len(cppcb.builders))
	mutators := make([]Mutator, len(cppcb.builders))
	for i := range cppcb.builders {
		func(i int, root context.Context) {
			builder := cppcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CoinProgramParticipantMutation)
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
					_, err = mutators[i+1].Mutate(root, cppcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = cppcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cppcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
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
		if _, err := mutators[0].Mutate(ctx, cppcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cppcb *CoinProgramParticipantCreateBulk) SaveX(ctx context.Context) []*CoinProgramParticipant {
	v, err := cppcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cppcb *CoinProgramParticipantCreateBulk) Exec(ctx context.Context) error {
	_, err := cppcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cppcb *CoinProgramParticipantCreateBulk) ExecX(ctx context.Context) {
	if err := cppcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.CoinProgramParticipant.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.CoinProgramParticipantUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (cppcb *CoinProgramParticipantCreateBulk) OnConflict(opts ...sql.ConflictOption) *CoinProgramParticipantUpsertBulk {
	cppcb.conflict = opts
	return &CoinProgramParticipantUpsertBulk{
		create: cppcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.CoinProgramParticipant.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cppcb *CoinProgramParticipantCreateBulk) OnConflictColumns(columns ...string) *CoinProgramParticipantUpsertBulk {
	cppcb.conflict = append(cppcb.conflict, sql.ConflictColumns(columns...))
	return &CoinProgramParticipantUpsertBulk{
		create: cppcb,
	}
}

// CoinProgramParticipantUpsertBulk is the builder for "upsert"-ing
// a bulk of CoinProgramParticipant nodes.
type CoinProgramParticipantUpsertBulk struct {
	create *CoinProgramParticipantCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.CoinProgramParticipant.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(coinprogramparticipant.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *CoinProgramParticipantUpsertBulk) UpdateNewValues() *CoinProgramParticipantUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(coinprogramparticipant.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(coinprogramparticipant.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.CoinProgramParticipant.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *CoinProgramParticipantUpsertBulk) Ignore() *CoinProgramParticipantUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *CoinProgramParticipantUpsertBulk) DoNothing() *CoinProgramParticipantUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the CoinProgramParticipantCreateBulk.OnConflict
// documentation for more info.
func (u *CoinProgramParticipantUpsertBulk) Update(set func(*CoinProgramParticipantUpsert)) *CoinProgramParticipantUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&CoinProgramParticipantUpsert{UpdateSet: update})
	}))
	return u
}

// SetBalance sets the "balance" field.
func (u *CoinProgramParticipantUpsertBulk) SetBalance(v uint) *CoinProgramParticipantUpsertBulk {
	return u.Update(func(s *CoinProgramParticipantUpsert) {
		s.SetBalance(v)
	})
}

// AddBalance adds v to the "balance" field.
func (u *CoinProgramParticipantUpsertBulk) AddBalance(v uint) *CoinProgramParticipantUpsertBulk {
	return u.Update(func(s *CoinProgramParticipantUpsert) {
		s.AddBalance(v)
	})
}

// UpdateBalance sets the "balance" field to the value that was provided on create.
func (u *CoinProgramParticipantUpsertBulk) UpdateBalance() *CoinProgramParticipantUpsertBulk {
	return u.Update(func(s *CoinProgramParticipantUpsert) {
		s.UpdateBalance()
	})
}

// Exec executes the query.
func (u *CoinProgramParticipantUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the CoinProgramParticipantCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for CoinProgramParticipantCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *CoinProgramParticipantUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
