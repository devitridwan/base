package controller

// func TestHandler_Overtime(t *testing.T) {
// 	initServer()
// 	overtimeReq := request.OvertimeRequest{
// 		Username: "johndoe",
// 		Password: "securepass",
// 		Date:     time.Now(),
// 		Hours:    2.5,
// 	}
// 	jsonBody, _ := json.Marshal(overtimeReq)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Overtime", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Overtime", mock.Anything, mock.Anything).Return(&entity.OvertimeResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/overtime", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }

// func TestHandler_Reimbursement(t *testing.T) {
// 	initServer()
// 	reimbReq := request.ReimbursementRequest{
// 		Username:    "johndoe",
// 		Password:    "securepass",
// 		Date:        time.Now(),
// 		Amount:      100000,
// 		Description: "Travel reimbursement",
// 	}
// 	jsonBody, _ := json.Marshal(reimbReq)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Reimbursement", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Reimbursement", mock.Anything, mock.Anything).Return(&entity.ReimbursementResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/reimbursement", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }

// func TestHandler_RunPayroll(t *testing.T) {
// 	initServer()
// 	req := request.PayrollRequest{
// 		Username: "admin",
// 		Password: "pass",
// 		PeriodID: uuid.New(),
// 	}
// 	jsonBody, _ := json.Marshal(req)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Payroll", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Payroll", mock.Anything, mock.Anything).Return(&entity.PayrollResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/payroll/run", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }

// func TestHandler_GenerateSummary(t *testing.T) {
// 	initServer()
// 	req := request.SummaryRequest{
// 		Username:  "admin",
// 		Password:  "pass",
// 		PayrollID: uuid.New(),
// 	}
// 	jsonBody, _ := json.Marshal(req)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Summary", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Summary", mock.Anything, mock.Anything).Return(&entity.PayrollSummaryResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/payroll/summary", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }

// func TestHandler_Payslip(t *testing.T) {
// 	initServer()
// 	req := request.PayslipRequest{
// 		Username:  "user",
// 		Password:  "pass",
// 		PayrollID: uuid.New(),
// 	}
// 	jsonBody, _ := json.Marshal(req)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("User", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("User", mock.Anything, mock.Anything).Return(&entity.PayslipResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/payslip", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }

// func TestHandler_CreateAttendance(t *testing.T) {
// 	initServer()
// 	req := request.AttendanceRequest{
// 		Username: "admin",
// 		Password: "pass",
// 		Date:     time.Now(),
// 	}
// 	jsonBody, _ := json.Marshal(req)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Attendance", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("Attendance", mock.Anything, mock.Anything).Return(&entity.AttendanceResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/attendance", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }

// func TestHandler_CreateAttendancePeriod(t *testing.T) {
// 	initServer()
// 	req := request.CreateAttendancePeriodRequest{
// 		Username:  "admin",
// 		Password:  "pass",
// 		StartDate: time.Now().AddDate(0, 0, -7),
// 		EndDate:   time.Now(),
// 	}
// 	jsonBody, _ := json.Marshal(req)
// 	tests := []struct {
// 		name       string
// 		body       string
// 		mockFn     func()
// 		wantStatus int
// 	}{
// 		{
// 			name:       "Invalid JSON",
// 			body:       "{invalid",
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:       "Validation Error",
// 			body:       `{}`,
// 			mockFn:     func() {},
// 			wantStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Interactor Error",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("CreatePeriod", mock.Anything, mock.Anything).Return(nil, errors.New("internal error")).Once()
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name: "Success",
// 			body: string(jsonBody),
// 			mockFn: func() {
// 				mockPayslipInteractor.On("CreatePeriod", mock.Anything, mock.Anything).Return(&entity.CreateAttendancePeriodResponse{}, nil).Once()
// 			},
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockFn()
// 			req := httptest.NewRequest("POST", "/v1/attendance-period", strings.NewReader(tt.body))
// 			req.Header.Set("Content-Type", "application/json")
// 			resp := httptest.NewRecorder()
// 			server.Httprouter.ServeHTTP(resp, req)
// 			require.Equal(t, tt.wantStatus, resp.Code)
// 		})
// 	}
// }
