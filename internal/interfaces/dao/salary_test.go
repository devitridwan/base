package dao

// func TestSalaryRepo_GetByEmployeeID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	sqlxDB := sqlx.NewDb(db, "sqlmock")
// 	repo := NewSalaryRepo(&OptSalary{
// 		DB: &database.Store{
// 			Master: sqlxDB,
// 			Slave:  sqlxDB,
// 		},
// 	})

// 	query := regexp.QuoteMeta(getSalaryByEmployeeIDQuery)

// 	employeeID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

// 	tests := []struct {
// 		name    string
// 		mockFn  func()
// 		want    *models.Salary
// 		wantErr bool
// 	}{
// 		{
// 			name: "success get salary",
// 			mockFn: func() {
// 				rows := sqlmock.NewRows([]string{"salary_id", "employee_id", "monthly_salary"}).
// 					AddRow(uuid.New(), employeeID, 10000000.0)
// 				mock.ExpectQuery(query).
// 					WithArgs(employeeID).
// 					WillReturnRows(rows)
// 			},
// 			want: &models.Salary{
// 				EmployeeID: employeeID,
// 				Monthly:    10000000.0,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "error no rows",
// 			mockFn: func() {
// 				mock.ExpectQuery(query).
// 					WithArgs(employeeID).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "unexpected db error",
// 			mockFn: func() {
// 				mock.ExpectQuery(query).
// 					WithArgs(employeeID).
// 					WillReturnError(errors.New("db failure"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt.mockFn()
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := repo.GetByEmployeeID(context.Background(), nil, employeeID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetByEmployeeID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.want != nil && got != nil {
// 				assert.Equal(t, tt.want.EmployeeID, got.EmployeeID)
// 				assert.Equal(t, tt.want.Monthly, got.Monthly)
// 			}
// 		})
// 	}
// }
