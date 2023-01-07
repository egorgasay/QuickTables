package handlers

import (
	"bytes"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"log"
	"net/http/httptest"
	"quicktables/internal/globals"
	"quicktables/internal/service"
	service_mocks "quicktables/internal/service/mocks"
	"strings"
	"testing"
)

func TestHandler_Login(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockIService)

	tests := []struct {
		name               string
		url                string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
		username           string
		password           string
	}{
		{
			name: "Bad Login",
			url:  "http://localhost:8080/login",
			mockBehavior: func(r *service_mocks.MockIService) {
				r.EXPECT().CheckPassword("qwe1q", "1234").
					Return(false).AnyTimes()
			},
			expectedStatusCode: 200,
			expectedBody:       "Wrong password or username",
			username:           "qwe1q",
			password:           "1234",
		},
		{
			name: "Good Login",
			url:  "http://localhost:8080/login",
			mockBehavior: func(r *service_mocks.MockIService) {
				r.EXPECT().CheckPassword("admin", "admin").
					Return(true).AnyTimes()
			},
			expectedStatusCode: 302,
			expectedBody:       "",
			username:           "admin",
			password:           "admin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repos := service_mocks.NewMockIService(c)
			test.mockBehavior(repos)

			services := &service.Service{DB: repos}
			handler := Handler{services}

			req := httptest.NewRequest("POST", test.url,
				bytes.NewBufferString(""))
			err := req.ParseForm()
			if err != nil {
				t.Fail()
				return
			}

			w := httptest.NewRecorder()
			r := gin.Default()

			r.LoadHTMLGlob("../../templates/html/*")
			r.Static("/static", "static")
			r.NoRoute(handler.NotFoundHandler)

			store := sessions.Store(cookie.NewStore(globals.Secret))
			r.Use(sessions.Sessions("session", store))

			r.POST("/login", handler.LoginHandler)

			req.PostForm.Set("username", test.username)
			req.PostForm.Set("password", test.password)

			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			if !strings.Contains(w.Body.String(), test.expectedBody) {
				t.Fail()
			}
		})
	}
}

func TestHandler_Reg(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockIService)

	tests := []struct {
		name               string
		url                string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
		username           string
		password           string
		passwordConfirm    string
	}{
		{
			name: "Bad Reg #1",
			url:  "http://localhost:8080/reg",
			mockBehavior: func(r *service_mocks.MockIService) {
				r.EXPECT().CheckPassword("admin", "admin").
					Return(true).AnyTimes()
			},
			expectedStatusCode: 200,
			expectedBody:       "Passwords don&#39;t match",
			username:           "qwe1q",
			password:           "1234",
			passwordConfirm:    "1333qdwdqwdq",
		},
		{
			name: "Bad Reg #2",
			url:  "http://localhost:8080/reg",
			mockBehavior: func(r *service_mocks.MockIService) {
				r.EXPECT().CreateUser("admin", "admin").
					Return(errors.New("UNIQUE constraint failed: Users.Name")).
					AnyTimes()
			},
			expectedStatusCode: 200,
			expectedBody:       "",
			username:           "admin",
			password:           "admin",
			passwordConfirm:    "admin",
		},
		{
			name: "Good Reg",
			url:  "http://localhost:8080/reg",
			mockBehavior: func(r *service_mocks.MockIService) {
				r.EXPECT().CreateUser("admin2", "admin2").
					Return(nil).AnyTimes()
			},
			expectedStatusCode: 308,
			expectedBody:       "",
			username:           "admin2",
			password:           "admin2",
			passwordConfirm:    "admin2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repos := service_mocks.NewMockIService(c)
			test.mockBehavior(repos)

			services := &service.Service{DB: repos}
			handler := Handler{services}

			req := httptest.NewRequest("POST", test.url,
				bytes.NewBufferString(""))
			err := req.ParseForm()
			if err != nil {
				t.Fail()
				return
			}

			w := httptest.NewRecorder()
			r := gin.Default()

			r.LoadHTMLGlob("../../templates/html/*")
			r.Static("/static", "static")
			r.NoRoute(handler.NotFoundHandler)

			store := sessions.Store(cookie.NewStore(globals.Secret))
			r.Use(sessions.Sessions("session", store))

			r.POST("/reg", handler.RegisterHandler)

			req.PostForm.Set("username", test.username)
			req.PostForm.Set("password", test.password)
			req.PostForm.Set("password2", test.passwordConfirm)

			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			if !strings.Contains(w.Body.String(), test.expectedBody) {
				log.Println(w.Body.String())
				t.Fail()
			}
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockIService)

	tests := []struct {
		name               string
		url                string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:               "Good Logout",
			url:                "http://localhost:8080/logout",
			mockBehavior:       func(r *service_mocks.MockIService) {},
			expectedStatusCode: 307,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repos := service_mocks.NewMockIService(c)
			test.mockBehavior(repos)

			services := &service.Service{DB: repos}
			handler := Handler{services}

			req := httptest.NewRequest("POST", test.url,
				bytes.NewBufferString(""))
			err := req.ParseForm()
			if err != nil {
				t.Fail()
				return
			}

			w := httptest.NewRecorder()
			r := gin.Default()

			r.LoadHTMLGlob("../../templates/html/*")
			r.Static("/static", "static")
			r.NoRoute(handler.NotFoundHandler)

			store := sessions.Store(cookie.NewStore(globals.Secret))
			r.Use(sessions.Sessions("session", store))

			r.POST("/logout", handler.LogoutHandler)

			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
