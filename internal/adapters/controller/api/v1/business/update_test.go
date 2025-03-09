package business_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"loyalit/internal/adapters/controller/api/v1/business"
	"loyalit/internal/adapters/controller/api/v1/business/mocks"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	businessServiceMock := mocks.NewMockbusinessService(ctrl)
	handler := business.NewHandler(nil, validator.New(), businessServiceMock, true)

	businessID := uuid.New()

	cases := []struct {
		name        string
		businessID  string
		requestBody dto.BusinessUpdateRequest

		mockSetup func()

		expCode int
		expResp interface{}

		errResp *echo.HTTPError
		respErr bool
	}{
		{
			name:       "Success",
			businessID: businessID.String(),
			requestBody: dto.BusinessUpdateRequest{
				Name:        "Updated Name",
				Description: "Updated Description",
			},

			mockSetup: func() {
				businessServiceMock.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.BusinessReturn{
						ID:          businessID,
						Name:        "Updated Name",
						Email:       "updated@example.com",
						Description: "Updated Description",
					}, nil)
			},
			expCode: http.StatusOK,
			expResp: &dto.BusinessReturn{
				ID:          businessID,
				Name:        "Updated Name",
				Email:       "updated@example.com",
				Description: "Updated Description",
			},
		},
		{
			name: "Empty Name",
			requestBody: dto.BusinessUpdateRequest{
				Description: "Updated Description",
			},

			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "[Name]: '' | Needs to implement 'required'",
			},
			respErr: true,
		},
		{
			name: "Invalid business id",
			requestBody: dto.BusinessUpdateRequest{
				Name:        "Name",
				Description: "Updated Description",
			},

			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid UUID length: 0",
			},
			respErr: true,
		},
		{
			name:       "Business not found",
			businessID: businessID.String(),
			requestBody: dto.BusinessUpdateRequest{
				Name:        "Name",
				Description: "Updated Description",
			},

			mockSetup: func() {
				businessServiceMock.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrBusinessNotFound)
			},

			errResp: &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "business not found",
			},
			respErr: true,
		},
		{
			name:       "Internal server error",
			businessID: businessID.String(),
			requestBody: dto.BusinessUpdateRequest{
				Name:        "Name",
				Description: "Updated Description",
			},

			mockSetup: func() {
				businessServiceMock.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any()).
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

			reqBody, err := json.Marshal(tCase.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("business_id", tCase.businessID)

			err = handler.Update(c)

			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)
				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.BusinessReturn
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
