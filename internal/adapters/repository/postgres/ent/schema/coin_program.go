package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"loyalit/internal/domain/utils/location"
	"time"
)

// CoinProgram holds the schema definition for the CoinProgram entity.
type CoinProgram struct {
	ent.Schema
}

// Fields of the CoinProgram.
func (CoinProgram) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.Time("created_at").
			Default(time.Now().In(location.Location())).
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.String("description"),
		field.Uint("day_limit").
			Min(1).
			Default(1),
		field.String("card_color").
			NotEmpty(),
	}
}

// Edges of the CoinProgram.
func (CoinProgram) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("business", Business.Type).
			Ref("coin_program").
			Unique().
			Required(),
		edge.To("rewards", Reward.Type),
		edge.From("participants", CoinProgramParticipant.Type).
			Ref("coin_program"),
	}
}
