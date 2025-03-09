package reward

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogramparticipant"
	"loyalit/internal/adapters/repository/postgres/ent/reward"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) UserListAvailable(ctx context.Context, req *dto.CoinProgramParticipantListAvailableRequest, userID uuid.UUID) ([]dto.CoinProgramWithRewardsReturn, error) {
	// Сначала получаем все программы пользователя с их балансами
	participants, err := s.db.CoinProgramParticipant.Query().
		Where(
			coinprogramparticipant.HasUserWith(
				user.ID(userID),
			),
		).
		WithCoinProgram().
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	// Создаем результирующий слайс
	result := make([]dto.CoinProgramWithRewardsReturn, 0, len(participants))

	// Для каждой программы получаем доступные награды
	for _, p := range participants {
		// Получаем награды для конкретной программы
		rewards, err := s.db.Reward.Query().
			Where(
				reward.And(
					reward.HasCoinProgramWith(
						coinprogram.ID(p.Edges.CoinProgram.ID),
					),
					reward.CostLTE(p.Balance), // Сразу фильтруем по балансу
				),
			).
			Order(ent.Asc(reward.FieldCost)).
			Limit(req.Limit).
			Offset(req.Offset).
			All(ctx)

		if err != nil {
			continue // Пропускаем программу при ошибке
		}

		// Если есть доступные награды, добавляем программу в результат
		if len(rewards) > 0 {
			rewardsDTO := make([]dto.RewardReturn, len(rewards))
			for i, r := range rewards {
				rewardsDTO[i] = dto.RewardReturn{
					ID:          r.ID,
					Name:        r.Name,
					Description: r.Description,
					Cost:        r.Cost,
					ImageURL:    r.ImageURL,
				}
			}

			result = append(result, dto.CoinProgramWithRewardsReturn{
				CoinProgram: dto.CoinProgramParticipantReturn{
					ID:          p.ID,
					Name:        p.Edges.CoinProgram.Name,
					Description: p.Edges.CoinProgram.Description,
					Balance:     p.Balance,
					CardColor:   p.Edges.CoinProgram.CardColor,
				},
				Rewards: rewardsDTO,
			})
		}
	}

	return result, nil
}
