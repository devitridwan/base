package user

// func TestModule_Payslip(t *testing.T) {
// 	m := initTest()

// 	payrollID := uuid.New()
// 	employeeID := uuid.New()
// 	payslipID := uuid.New()

// 	req := &request.PayslipRequest{
// 		Username:  "user1",
// 		Password:  "pass1",
// 		PayrollID: payrollID,
// 	}

// 	tests := []struct {
// 		name    string
// 		req     *request.PayslipRequest
// 		wantErr bool
// 		setup   func()
// 	}{
// 		{
// 			name:    "error getting employee",
// 			req:     req,
// 			wantErr: true,
// 			setup: func() {
// 				mockEmployeeRepo.On("GetByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(nil, errors.New("employee not found")).Once()
// 			},
// 		},
// 		{
// 			name:    "error getting payroll",
// 			req:     req,
// 			wantErr: true,
// 			setup: func() {
// 				mockEmployeeRepo.On("GetByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Employee{EmployeeID: employeeID}, nil).Once()
// 				mockPayrollRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
// 					Return(nil, errors.New("payroll error")).Once()
// 			},
// 		},
// 		{
// 			name:    "error getting reimbursement",
// 			req:     req,
// 			wantErr: true,
// 			setup: func() {
// 				mockEmployeeRepo.On("GetByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Employee{EmployeeID: employeeID}, nil).Once()
// 				mockPayrollRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Payroll{PayrollID: payrollID}, nil).Once()
// 				mockReimbursementRepo.On("ListByPeriod", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(nil, errors.New("reimb error")).Once()
// 			},
// 		},
// 		{
// 			name:    "error getting payslip",
// 			req:     req,
// 			wantErr: true,
// 			setup: func() {
// 				mockEmployeeRepo.On("GetByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Employee{EmployeeID: employeeID}, nil).Once()
// 				mockPayrollRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Payroll{PayrollID: payrollID}, nil).Once()
// 				mockReimbursementRepo.On("ListByPeriod", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return([]models.Reimbursement{}, nil).Once()
// 				mockPayslipRepo.On("GetByEmployeeAndPayroll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(nil, errors.New("payslip error")).Once()
// 			},
// 		},
// 		{
// 			name:    "success",
// 			req:     req,
// 			wantErr: false,
// 			setup: func() {
// 				mockEmployeeRepo.On("GetByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Employee{EmployeeID: employeeID}, nil).Once()
// 				mockPayrollRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.Payroll{PayrollID: payrollID}, nil).Once()
// 				mockReimbursementRepo.On("ListByPeriod", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return([]models.Reimbursement{
// 						{Description: "Transport", Amount: 100000},
// 						{Description: "Internet", Amount: 150000},
// 					}, nil).Once()
// 				mockPayslipRepo.On("GetByEmployeeAndPayroll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(&models.User{
// 						PayslipID:           payslipID,
// 						EmployeeID:          employeeID,
// 						PayrollID:           payrollID,
// 						BaseSalary:          5000000,
// 						ProratedSalary:      4800000,
// 						OvertimeHours:       2,
// 						OvertimeAmount:      100000,
// 						ReimbursementAmount: 250000,
// 						TotalTakeHomePay:    5150000,
// 						DaysAttended:        20,
// 					}, nil).Once()
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
// 			got, err := m.User(context.Background(), tt.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("User() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !tt.wantErr && got == nil {
// 				t.Error("Expected non-nil result on success")
// 			}
// 		})
// 	}
// }
