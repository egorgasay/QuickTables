package handlers

//func TestHandler_AddDB(t *testing.T) {
//	type mockBehavior func(r *service_mocks.MockIService)
//
//	tests := []struct {
//		name               string
//		url                string
//		mockBehavior       mockBehavior
//		expectedStatusCode int
//		expectedBody       string
//		dbName             string
//		connStr            string
//		dbVendorName       string
//	}{
//		{
//			name: "Bad AddDB",
//			url:  "http://localhost:8080/addDB",
//			mockBehavior: func(r *service_mocks.MockIService) {
//				r.EXPECT().CheckPassword("admin", "admin").
//					Return(true).AnyTimes()
//			},
//			expectedStatusCode: 302,
//			expectedBody:       "Error",
//			dbName:             "qwe1q",
//			connStr:            "1234",
//			dbVendorName:       "1234",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//			repos := service_mocks.NewMockIService(c)
//			test.mockBehavior(repos)
//
//			services := &service.Service{DB: repos}
//			handler := Handler{services}
//			req := httptest.NewRequest("POST", "http://localhost:8080/login",
//				bytes.NewBufferString(""))
//			err := req.ParseForm()
//			if err != nil {
//				t.Fail()
//				return
//			}
//			w := httptest.NewRecorder()
//			r := gin.Default()
//
//			r.LoadHTMLGlob("../../templates/html/*")
//			r.Static("/static", "static")
//			r.NoRoute(handler.NotFoundHandler)
//
//			store := sessions.Store(cookie.NewStore(globals.Secret))
//			r.Use(sessions.Sessions("session", store))
//
//			r.POST("/addDB", handler.AddDBPostHandler)
//
//			r.POST("/login", handler.LoginHandler)
//
//			req.PostForm.Set("username", "admin")
//			req.PostForm.Set("password", "admin")
//			r.ServeHTTP(w, req)
//
//			req = httptest.NewRequest("POST", test.url,
//				bytes.NewBufferString(""))
//
//			err = req.ParseForm()
//			if err != nil {
//				t.Fail()
//				return
//			}
//
//			req.PostForm.Set("dbName", test.dbName)
//			req.PostForm.Set("con_str", test.connStr)
//			req.PostForm.Set("bdVendorName", test.dbVendorName)
//
//			r.ServeHTTP(w, req)
//
//			// Assert
//			log.Println(w.Body.String())
//			assert.Equal(t, test.expectedStatusCode, w.Code)
//			if !strings.Contains(w.Body.String(), test.expectedBody) {
//				t.Fail()
//			}
//		})
//	}
//}
