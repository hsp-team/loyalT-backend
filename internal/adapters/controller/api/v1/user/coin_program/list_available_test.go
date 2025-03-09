package coin_program_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/user/coin_program"
	"loyalit/internal/adapters/controller/api/v1/user/coin_program/mocks"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rewardServiceMock := mocks.NewMockrewardService(ctrl)
	handler := coin_program.NewHandler(validator.New(), nil, rewardServiceMock)

	userID := uuid.New()
	coinProgramParticipantID1 := uuid.New()
	rewardID1 := uuid.New()

	coinProgramParticipantID2 := uuid.New()
	rewardID2 := uuid.New()

	cases := []struct {
		name   string
		body   interface{}
		userID string

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name: "Success",
			body: &dto.CoinProgramParticipantListRequest{
				Limit:  5,
				Offset: 0,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					UserListAvailable(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]dto.CoinProgramWithRewardsReturn{
						{
							CoinProgram: dto.CoinProgramParticipantReturn{
								ID:        coinProgramParticipantID1,
								Name:      "Тест1",
								Balance:   10,
								CardColor: "#fff",
							},
							Rewards: []dto.RewardReturn{
								{
									ID:       rewardID1,
									Name:     "Награда1",
									Cost:     2,
									ImageURL: "Картинка1",
								},
							},
						},
						{
							CoinProgram: dto.CoinProgramParticipantReturn{
								ID:        coinProgramParticipantID2,
								Name:      "Тест2",
								Balance:   5,
								CardColor: "#aaa",
							},
							Rewards: []dto.RewardReturn{
								{
									ID:       rewardID2,
									Name:     "Награда2",
									Cost:     2,
									ImageURL: "Картинка2",
								},
							},
						},
					}, nil)
			},

			expCode: http.StatusOK,
			expResp: []dto.CoinProgramWithRewardsReturn{
				{
					CoinProgram: dto.CoinProgramParticipantReturn{
						ID:        coinProgramParticipantID1,
						Name:      "Тест1",
						Balance:   10,
						CardColor: "#fff",
					},
					Rewards: []dto.RewardReturn{
						{
							ID:       rewardID1,
							Name:     "Награда1",
							Cost:     2,
							ImageURL: "Картинка1",
						},
					},
				},
				{
					CoinProgram: dto.CoinProgramParticipantReturn{
						ID:        coinProgramParticipantID2,
						Name:      "Тест2",
						Balance:   5,
						CardColor: "#aaa",
					},
					Rewards: []dto.RewardReturn{
						{
							ID:       rewardID2,
							Name:     "Награда2",
							Cost:     2,
							ImageURL: "Картинка2",
						},
					},
				},
			},
		},
		{
			name: "Invalid limit and offset",
			body: &dto.CoinProgramParticipantListRequest{
				Limit:  -1,
				Offset: -1,
			},

			mockSetup: func() {},

			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Limit]: '-1' | Needs to implement 'gte' and [Offset]: '-1' | Needs to implement 'gte'",
			},
			respErr: true,
		},
		{
			name: "Invalid user id",
			body: &dto.CoinProgramParticipantListRequest{
				Limit:  5,
				Offset: 0,
			},

			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name: "Internal error",
			body: &dto.CoinProgramParticipantListRequest{
				Limit:  5,
				Offset: 0,
			},
			userID: userID.String(),

			mockSetup: func() {
				rewardServiceMock.EXPECT().
					UserListAvailable(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "internal error",
			},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			tCase.mockSetup()

			bodyBytes, err := json.Marshal(tCase.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/list_available", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("user_id", tCase.userID)

			err = handler.ListAvailable(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got []dto.CoinProgramWithRewardsReturn
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
