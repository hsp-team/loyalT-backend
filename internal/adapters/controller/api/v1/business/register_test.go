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

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	businessServiceMock := mocks.NewMockbusinessService(ctrl)
	handler := business.NewHandler(nil, validator.New(), businessServiceMock, true)

	businessID := uuid.New()

	cases := []struct {
		name string
		body interface{} // тело запроса

		mockSetup func() // настройка мока

		expCode int         // ожидаемый HTTP код
		expResp interface{} // ожидаемый ответ

		errResp *echo.HTTPError // ожидаемая ошибка
		respErr bool            // флаг, указывающий на ошибку
	}{
		{
			name: "Success",
			body: dto.BusinessRegisterRequest{
				Name:     "Тест",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(&dto.BusinessRegisterResponse{
						ID:    businessID,
						Name:  "Тест",
						Email: "test@example.com",
					}, nil)
			},
			expCode: http.StatusCreated,
			expResp: &dto.BusinessRegisterResponse{
				ID:    businessID,
				Name:  "Тест",
				Email: "test@example.com",
			},
		},
		{
			name: "Empty name",
			body: dto.BusinessRegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {},
			errResp: &echo.HTTPError{
				Message: "[Name]: '' | Needs to implement 'required'",
				Code:    echo.ErrBadRequest.Code,
			},
			respErr: true,
		},
		{
			name: "Invalid password",
			body: dto.BusinessRegisterRequest{
				Name:     "Тест",
				Email:    "failed-email",
				Password: "password123",
			},
			mockSetup: func() {}, // валидация происходит до вызова сервиса
			errResp: &echo.HTTPError{
				Message: "[Email]: 'failed-email' | Needs to implement 'email'",
				Code:    echo.ErrBadRequest.Code,
			},
			respErr: true,
		},
		{
			name: "Email already exists",
			body: dto.BusinessRegisterRequest{
				Name:     "Тест",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil, errorz.ErrEmailAlreadyExists)
			},
			errResp: &echo.HTTPError{
				Message: "email already exists",
				Code:    echo.ErrConflict.Code,
			},
			respErr: true,
		},
		{
			name: "Internal server error",
			body: dto.BusinessRegisterRequest{
				Name:     "Тест",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				businessServiceMock.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
			},
			errResp: &echo.HTTPError{
				Message: "internal error",
				Code:    echo.ErrInternalServerError.Code,
			},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			tCase.mockSetup()

			// Создаем JSON из тела запроса
			bodyBytes, err := json.Marshal(tCase.body)
			require.NoError(t, err)

			// Создаем тестовый запрос и рекордер
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Создаем Echo и контекст правильным способом
			e := echo.New()
			c := e.NewContext(req, rec)

			// Вызываем handler напрямую
			err = handler.Register(c)

			// Проверяем результат
			if tCase.respErr {
				require.Error(t, err)
				var he *echo.HTTPError
				ok := errors.As(err, &he)
				require.True(t, ok)

				require.Equal(t, tCase.errResp, he)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.expCode, rec.Code)

				var got *dto.BusinessRegisterResponse
				err = json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				require.Equal(t, tCase.expResp, got)
			}
		})
	}
}
