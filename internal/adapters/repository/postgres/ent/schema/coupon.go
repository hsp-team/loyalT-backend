package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"loyalit/internal/domain/utils/location"
	"time"
)

// Coupon holds the schema definition for the Coupon entity.
type Coupon struct {
	ent.Schema
}

// Fields of the Coupon.
func (Coupon) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.Time("created_at").
			Default(time.Now().In(location.Location())).
			Immutable(),
		field.String("code").
			Unique().
			NotEmpty(),
		field.Bool("activated").
			Default(false),
	}
}

// Edges of the Coupon.
func (Coupon) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("coupons").
			Unique(),
		edge.From("reward", Reward.Type).
			Ref("rewards").
			Unique().
			Required(),
	}
}

func (Coupon) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Unique(),
	}
}
