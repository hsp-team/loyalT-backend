package qr

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/utils"
)

func (s *Service) GetUserQR(ctx context.Context, userID uuid.UUID) (*dto.QRGetResponse, error) {
	user, err := s.db.User.Get(ctx, userID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrUserNotFound
	case err != nil:
		return nil, err
	}

	if user.QrData == "" {
		code, errGenerateCode := utils.GenerateCode(s.codeLength)
		if errGenerateCode != nil {
			return nil, err
		}

		user, err = user.Update().SetQrData(code).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &dto.QRGetResponse{
		Data: user.QrData,
	}, nil
}
