package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"loyalit/internal/domain/utils/location"
	"time"
)

// Reward holds the schema definition for the Reward entity.
type Reward struct {
	ent.Schema
}

// Fields of the Reward.
func (Reward) Fields() []ent.Field {
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
		field.Uint("cost").
			Min(1),
		field.String("image_url").
			NotEmpty(),
	}
}

// Edges of the Reward.
func (Reward) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("coin_program", CoinProgram.Type).
			Ref("rewards").
			Unique(),
		edge.To("rewards", Coupon.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
