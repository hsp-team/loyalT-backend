package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"loyalit/internal/domain/utils/location"
	"time"
)

// CoinProgramParticipant holds the schema definition for the CoinProgramParticipant entity.
type CoinProgramParticipant struct {
	ent.Schema
}

// Fields of the CoinProgramParticipant.
func (CoinProgramParticipant) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.Time("created_at").
			Default(time.Now().In(location.Location())).
			Immutable(),
		field.Uint("balance").
			Default(0),
	}
}

// Edges of the CoinProgramParticipant.
func (CoinProgramParticipant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("coin_program", CoinProgram.Type).
			Unique(),
		edge.From("user", User.Type).
			Ref("coin_programs").
			Unique(),
	}
}
