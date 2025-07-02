package user

// var (
// 	mockAdminRepo            *mocks_domain.AdminRepository
// 	mockAttendancePeriodRepo *mocks_domain.AttendancePeriodRepository
// 	mockAttendanceRepo       *mocks_domain.AttendanceRepository
// 	mockAuditLogRepo         *mocks_domain.AuditLogRepository
// 	mockEmployeeRepo         *mocks_domain.EmployeeRepository
// 	mockOvertimeRepo         *mocks_domain.OvertimeRepository
// 	mockPayrollRepo          *mocks_domain.PayrollRepository
// 	mockPayslipRepo          *mocks_domain.PayslipRepository
// 	mockReimbursementRepo    *mocks_domain.ReimbursementRepository
// 	mockSalaryRepo           *mocks_domain.SalaryRepository
// 	mockTxManager            *mocks_interfaces.TxManager
// )

// func initTest() interactor.PayslipInteractor {
// 	mockAdminRepo = &mocks_domain.AdminRepository{}
// 	mockAttendancePeriodRepo = &mocks_domain.AttendancePeriodRepository{}
// 	mockAttendanceRepo = &mocks_domain.AttendanceRepository{}
// 	mockAuditLogRepo = &mocks_domain.AuditLogRepository{}
// 	mockEmployeeRepo = &mocks_domain.EmployeeRepository{}
// 	mockOvertimeRepo = &mocks_domain.OvertimeRepository{}
// 	mockPayrollRepo = &mocks_domain.PayrollRepository{}
// 	mockPayslipRepo = &mocks_domain.PayslipRepository{}
// 	mockReimbursementRepo = &mocks_domain.ReimbursementRepository{}
// 	mockSalaryRepo = &mocks_domain.SalaryRepository{}
// 	mockTxManager = &mocks_interfaces.TxManager{}

// 	return New(&Opts{
// 		AdminRepository:            mockAdminRepo,
// 		AttendancePeriodRepository: mockAttendancePeriodRepo,
// 		AttendanceRepository:       mockAttendanceRepo,
// 		AuditLogRepository:         mockAuditLogRepo,
// 		EmployeeRepository:         mockEmployeeRepo,
// 		OvertimeRepository:         mockOvertimeRepo,
// 		PayrollRepository:          mockPayrollRepo,
// 		PayslipRepository:          mockPayslipRepo,
// 		ReimbursementRepository:    mockReimbursementRepo,
// 		SalaryRepository:           mockSalaryRepo,
// 		TxManager:                  mockTxManager,
// 	})
// }

// func TestDelivery_New(t *testing.T) {
// 	t.Run("Testdelivery_New", func(t *testing.T) {
// 		opt := &Opts{
// 			AdminRepository: mockAdminRepo,
// 		}
// 		tx := New(opt)
// 		if tx == nil {
// 			t.Errorf("New() = %+v, want %+v", tx, tx)
// 		}
// 	})
// }
