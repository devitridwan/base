-- Table: admins
-- Stores information about administrative users
CREATE TABLE admins (
    admin_id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID, -- Could reference admin_id or a system user
    updated_by UUID, -- Could reference admin_id
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ
);

-- Index for efficient username lookups and filtering by deletion status
CREATE INDEX idx_admins_username ON admins (username);
CREATE INDEX idx_admins_is_deleted ON admins (is_deleted);
CREATE INDEX idx_admins_created_at ON admins (created_at);

-- Table: employees
-- Stores information about employees
CREATE TABLE employees (
    employee_id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES admins(admin_id), -- Assuming created by an admin
    updated_by UUID REFERENCES admins(admin_id), -- Assuming updated by an admin
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ
);

-- Index for efficient username lookups and filtering by deletion status
CREATE INDEX idx_employees_username ON employees (username);
CREATE INDEX idx_employees_is_deleted ON employees (is_deleted);
CREATE INDEX idx_employees_created_at ON employees (created_at);


-- Table: attendance_periods
-- Defines the start and end dates for attendance tracking periods
CREATE TABLE attendance_periods (
    attendance_period_id UUID PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES admins(admin_id), -- Typically created by an admin
    updated_by UUID REFERENCES admins(admin_id), -- Typically updated by an admin
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ
);

-- Indexes for date range queries and deletion status
CREATE INDEX idx_attendance_periods_start_date ON attendance_periods (start_date);
CREATE INDEX idx_attendance_periods_end_date ON attendance_periods (end_date);
CREATE INDEX idx_attendance_periods_is_deleted ON attendance_periods (is_deleted);


-- Table: attendance
-- Records daily attendance for employees
CREATE TABLE attendance (
    attendance_id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(employee_id),
    date DATE NOT NULL,
    attendance_period_id UUID NOT NULL REFERENCES attendance_periods(attendance_period_id),
    created_by_employee_id UUID REFERENCES employees(employee_id), -- Employee who recorded their attendance (e.g., self-service)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID, -- Could be an admin or the employee themselves
    updated_by UUID, -- Could be an admin or the employee themselves
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ,
    ip_address VARCHAR(45) -- To store IPv4 or IPv6 addresses
);

-- Composite unique index for frequent queries by employee and date within a period to prevent duplicate attendance records
CREATE UNIQUE INDEX uidx_attendance_employee_date_period ON attendance (employee_id, date, attendance_period_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_attendance_employee_id ON attendance (employee_id);
CREATE INDEX idx_attendance_date ON attendance (date);
CREATE INDEX idx_attendance_attendance_period_id ON attendance (attendance_period_id);
CREATE INDEX idx_attendance_is_deleted ON attendance (is_deleted);


-- Table: overtime
-- Records overtime hours for employees
CREATE TABLE overtime (
    overtime_id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(employee_id),
    date DATE NOT NULL,
    hours DOUBLE PRECISION NOT NULL CHECK (hours >= 0),
    attendance_period_id UUID NOT NULL REFERENCES attendance_periods(attendance_period_id),
    created_by_employee_id UUID REFERENCES employees(employee_id), -- Employee who submitted overtime
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID, -- Could be an admin or the employee themselves
    updated_by UUID, -- Could be an admin or the employee themselves
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ,
    ip_address VARCHAR(45)
);

-- Indexes for efficient queries
CREATE INDEX idx_overtime_employee_id ON overtime (employee_id);
CREATE INDEX idx_overtime_date ON overtime (date);
CREATE INDEX idx_overtime_attendance_period_id ON overtime (attendance_period_id);
CREATE INDEX idx_overtime_is_deleted ON overtime (is_deleted);


-- Table: reimbursements
-- Records reimbursements for employees
CREATE TABLE reimbursements (
    reimbursement_id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(employee_id),
    date DATE NOT NULL,
    amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
    description TEXT,
    attendance_period_id UUID NOT NULL REFERENCES attendance_periods(attendance_period_id),
    created_by_employee_id UUID REFERENCES employees(employee_id), -- Employee who submitted reimbursement
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID, -- Could be an admin or the employee themselves
    updated_by UUID, -- Could be an admin or the employee themselves
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ,
    ip_address VARCHAR(45)
);

-- Indexes for efficient queries
CREATE INDEX idx_reimbursements_employee_id ON reimbursements (employee_id);
CREATE INDEX idx_reimbursements_date ON reimbursements (date);
CREATE INDEX idx_reimbursements_attendance_period_id ON reimbursements (attendance_period_id);
CREATE INDEX idx_reimbursements_is_deleted ON reimbursements (is_deleted);


-- Table: payrolls
-- Records the generation of payrolls for specific attendance periods
CREATE TABLE payrolls (
    payroll_id UUID PRIMARY KEY,
    attendance_period_id UUID NOT NULL REFERENCES attendance_periods(attendance_period_id),
    request_id UUID NOT NULL, -- To track specific payroll generation requests
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES admins(admin_id), -- Payrolls are likely initiated by an admin
    updated_by UUID REFERENCES admins(admin_id), -- Payrolls are likely updated by an admin
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ
);

-- Unique constraint for attendance period to payroll, assuming one active payroll per period
CREATE UNIQUE INDEX uidx_payrolls_attendance_period_id ON payrolls (attendance_period_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_payrolls_request_id ON payrolls (request_id);
CREATE INDEX idx_payrolls_is_deleted ON payrolls (is_deleted);


-- Table: payslips
-- Stores generated payslips for employees for each payroll
CREATE TABLE payslips (
    payslip_id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(employee_id),
    payroll_id UUID NOT NULL REFERENCES payrolls(payroll_id),
    total_working_days INTEGER NOT NULL CHECK (total_working_days >= 0),
    days_attended INTEGER NOT NULL CHECK (days_attended >= 0),
    base_salary NUMERIC(12, 2) NOT NULL CHECK (base_salary >= 0),
    prorated_salary NUMERIC(12, 2) NOT NULL CHECK (prorated_salary >= 0),
    overtime_hours DOUBLE PRECISION NOT NULL CHECK (overtime_hours >= 0),
    overtime_amount NUMERIC(12, 2) NOT NULL CHECK (overtime_amount >= 0),
    reimbursement_amount NUMERIC(12, 2) NOT NULL CHECK (reimbursement_amount >= 0),
    total_take_home_pay NUMERIC(12, 2) NOT NULL CHECK (total_take_home_pay >= 0),
    generated_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES admins(admin_id), -- Payslips are generated by the system/admin
    updated_by UUID REFERENCES admins(admin_id), -- Payslips are updated by the system/admin
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ
);

-- Composite unique index to ensure one active payslip per employee per payroll
CREATE UNIQUE INDEX uidx_payslips_employee_payroll ON payslips (employee_id, payroll_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_payslips_employee_id ON payslips (employee_id);
CREATE INDEX idx_payslips_payroll_id ON payslips (payroll_id);
CREATE INDEX idx_payslips_generated_at ON payslips (generated_at);
CREATE INDEX idx_payslips_is_deleted ON payslips (is_deleted);


-- Table: salaries
-- Stores the salary information for employees
CREATE TABLE salaries (
    salary_id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(employee_id),
    monthly_salary NUMERIC(12, 2) NOT NULL CHECK (monthly_salary >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES admins(admin_id), -- Salaries are set by admin
    updated_by UUID REFERENCES admins(admin_id), -- Salaries are updated by admin
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMPTZ
);

-- Unique index to ensure one active salary record per employee at any time
-- If you need historical salaries, you would add a 'valid_from'/'valid_to' date and adjust this index.
CREATE UNIQUE INDEX uidx_salaries_employee_id_active ON salaries (employee_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_salaries_is_deleted ON salaries (is_deleted);


-- Table: audit_logs
-- Records all significant actions performed in the system for auditing purposes
CREATE TABLE audit_logs (
    audit_log_id UUID PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id UUID NOT NULL, -- The ID of the record that was audited
    action VARCHAR(50) NOT NULL, -- e.g., 'CREATE', 'UPDATE', 'DELETE'
    performed_by_admin_id UUID REFERENCES admins(admin_id),
    performed_by_employee_id UUID REFERENCES employees(employee_id),
    performed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR(45),
    request_id UUID, -- For correlating actions within a single request/transaction
    changes JSONB -- Store the old/new values as a JSON blob
);

-- Indexes for efficient querying of audit logs
CREATE INDEX idx_audit_logs_table_name ON audit_logs (table_name);
CREATE INDEX idx_audit_logs_record_id ON audit_logs (record_id);
CREATE INDEX idx_audit_logs_action ON audit_logs (action);
CREATE INDEX idx_audit_logs_performed_by_admin_id ON audit_logs (performed_by_admin_id);
CREATE INDEX idx_audit_logs_performed_by_employee_id ON audit_logs (performed_by_employee_id);
CREATE INDEX idx_audit_logs_performed_at ON audit_logs (performed_at DESC); -- For most recent logs
CREATE INDEX idx_audit_logs_request_id ON audit_logs (request_id);

-- Run this once if the uuid-ossp extension is not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO admins (
    admin_id,
    username,
    password_hash,
    full_name,
    created_at,
    updated_at,
    created_by,
    updated_by,
    is_deleted,
    deleted_time
) VALUES (
    'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', -- Static UUID for easy reference by employees
    'admin.main',
    'hashed_password_for_admin', -- In a real application, this would be a securely hashed password
    'Main Administrator',
    NOW(),
    NOW(),
    NULL, -- No creator for the very first admin, or use a system UUID
    NULL,
    FALSE,
    NULL
);

INSERT INTO employees (
    employee_id,
    username,
    password_hash,
    full_name,
    created_at,
    updated_at,
    created_by,
    updated_by,
    is_deleted,
    deleted_time
) VALUES
-- Employees 1-100
('e0000001-0000-4000-8000-000000000001', 'employee.001', 'hashed_password_001', 'Employee One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000002-0000-4000-8000-000000000002', 'employee.002', 'hashed_password_002', 'Employee Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000003-0000-4000-8000-000000000003', 'employee.003', 'hashed_password_003', 'Employee Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000004-0000-4000-8000-000000000004', 'employee.004', 'hashed_password_004', 'Employee Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000005-0000-4000-8000-000000000005', 'employee.005', 'hashed_password_005', 'Employee Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000006-0000-4000-8000-000000000006', 'employee.006', 'hashed_password_006', 'Employee Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000007-0000-4000-8000-000000000007', 'employee.007', 'hashed_password_007', 'Employee Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000008-0000-4000-8000-000000000008', 'employee.008', 'hashed_password_008', 'Employee Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000009-0000-4000-8000-000000000009', 'employee.009', 'hashed_password_009', 'Employee Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000010-0000-4000-8000-000000000010', 'employee.010', 'hashed_password_010', 'Employee Ten', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000011-0000-4000-8000-000000000011', 'employee.011', 'hashed_password_011', 'Employee Eleven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000012-0000-4000-8000-000000000012', 'employee.012', 'hashed_password_012', 'Employee Twelve', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000013-0000-4000-8000-000000000013', 'employee.013', 'hashed_password_013', 'Employee Thirteen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000014-0000-4000-8000-000000000014', 'employee.014', 'hashed_password_014', 'Employee Fourteen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000015-0000-4000-8000-000000000015', 'employee.015', 'hashed_password_015', 'Employee Fifteen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000016-0000-4000-8000-000000000016', 'employee.016', 'hashed_password_016', 'Employee Sixteen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000017-0000-4000-8000-000000000017', 'employee.017', 'hashed_password_017', 'Employee Seventeen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000018-0000-4000-8000-000000000018', 'employee.018', 'hashed_password_018', 'Employee Eighteen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000019-0000-4000-8000-000000000019', 'employee.019', 'hashed_password_019', 'Employee Nineteen', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000020-0000-4000-8000-000000000020', 'employee.020', 'hashed_password_020', 'Employee Twenty', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000021-0000-4000-8000-000000000021', 'employee.021', 'hashed_password_021', 'Employee Twenty-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000022-0000-4000-8000-000000000022', 'employee.022', 'hashed_password_022', 'Employee Twenty-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000023-0000-4000-8000-000000000023', 'employee.023', 'hashed_password_023', 'Employee Twenty-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000024-0000-4000-8000-000000000024', 'employee.024', 'hashed_password_024', 'Employee Twenty-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000025-0000-4000-8000-000000000025', 'employee.025', 'hashed_password_025', 'Employee Twenty-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000026-0000-4000-8000-000000000026', 'employee.026', 'hashed_password_026', 'Employee Twenty-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000027-0000-4000-8000-000000000027', 'employee.027', 'hashed_password_027', 'Employee Twenty-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000028-0000-4000-8000-000000000028', 'employee.028', 'hashed_password_028', 'Employee Twenty-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000029-0000-4000-8000-000000000029', 'employee.029', 'hashed_password_029', 'Employee Twenty-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000030-0000-4000-8000-000000000030', 'employee.030', 'hashed_password_030', 'Employee Thirty', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000031-0000-4000-8000-000000000031', 'employee.031', 'hashed_password_031', 'Employee Thirty-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000032-0000-4000-8000-000000000032', 'employee.032', 'hashed_password_032', 'Employee Thirty-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000033-0000-4000-8000-000000000033', 'employee.033', 'hashed_password_033', 'Employee Thirty-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000034-0000-4000-8000-000000000034', 'employee.034', 'hashed_password_034', 'Employee Thirty-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000035-0000-4000-8000-000000000035', 'employee.035', 'hashed_password_035', 'Employee Thirty-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000036-0000-4000-8000-000000000036', 'employee.036', 'hashed_password_036', 'Employee Thirty-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000037-0000-4000-8000-000000000037', 'employee.037', 'hashed_password_037', 'Employee Thirty-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000038-0000-4000-8000-000000000038', 'employee.038', 'hashed_password_038', 'Employee Thirty-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000039-0000-4000-8000-000000000039', 'employee.039', 'hashed_password_039', 'Employee Thirty-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000040-0000-4000-8000-000000000040', 'employee.040', 'hashed_password_040', 'Employee Forty', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000041-0000-4000-8000-000000000041', 'employee.041', 'hashed_password_041', 'Employee Forty-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000042-0000-4000-8000-000000000042', 'employee.042', 'hashed_password_042', 'Employee Forty-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000043-0000-4000-8000-000000000043', 'employee.043', 'hashed_password_043', 'Employee Forty-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000044-0000-4000-8000-000000000044', 'employee.044', 'hashed_password_044', 'Employee Forty-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000045-0000-4000-8000-000000000045', 'employee.045', 'hashed_password_045', 'Employee Forty-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000046-0000-4000-8000-000000000046', 'employee.046', 'hashed_password_046', 'Employee Forty-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000047-0000-4000-8000-000000000047', 'employee.047', 'hashed_password_047', 'Employee Forty-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000048-0000-4000-8000-000000000048', 'employee.048', 'hashed_password_048', 'Employee Forty-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000049-0000-4000-8000-000000000049', 'employee.049', 'hashed_password_049', 'Employee Forty-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000050-0000-4000-8000-000000000050', 'employee.050', 'hashed_password_050', 'Employee Fifty', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000051-0000-4000-8000-000000000051', 'employee.051', 'hashed_password_051', 'Employee Fifty-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000052-0000-4000-8000-000000000052', 'employee.052', 'hashed_password_052', 'Employee Fifty-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000053-0000-4000-8000-000000000053', 'employee.053', 'hashed_password_053', 'Employee Fifty-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000054-0000-4000-8000-000000000054', 'employee.054', 'hashed_password_054', 'Employee Fifty-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000055-0000-4000-8000-000000000055', 'employee.055', 'hashed_password_055', 'Employee Fifty-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000056-0000-4000-8000-000000000056', 'employee.056', 'hashed_password_056', 'Employee Fifty-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000057-0000-4000-8000-000000000057', 'employee.057', 'hashed_password_057', 'Employee Fifty-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000058-0000-4000-8000-000000000058', 'employee.058', 'hashed_password_058', 'Employee Fifty-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000059-0000-4000-8000-000000000059', 'employee.059', 'hashed_password_059', 'Employee Fifty-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000060-0000-4000-8000-000000000060', 'employee.060', 'hashed_password_060', 'Employee Sixty', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000061-0000-4000-8000-000000000061', 'employee.061', 'hashed_password_061', 'Employee Sixty-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000062-0000-4000-8000-000000000062', 'employee.062', 'hashed_password_062', 'Employee Sixty-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000063-0000-4000-8000-000000000063', 'employee.063', 'hashed_password_063', 'Employee Sixty-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000064-0000-4000-8000-000000000064', 'employee.064', 'hashed_password_064', 'Employee Sixty-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000065-0000-4000-8000-000000000065', 'employee.065', 'hashed_password_065', 'Employee Sixty-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000066-0000-4000-8000-000000000066', 'employee.066', 'hashed_password_066', 'Employee Sixty-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000067-0000-4000-8000-000000000067', 'employee.067', 'hashed_password_067', 'Employee Sixty-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000068-0000-4000-8000-000000000068', 'employee.068', 'hashed_password_068', 'Employee Sixty-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000069-0000-4000-8000-000000000069', 'employee.069', 'hashed_password_069', 'Employee Sixty-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000070-0000-4000-8000-000000000070', 'employee.070', 'hashed_password_070', 'Employee Seventy', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000071-0000-4000-8000-000000000071', 'employee.071', 'hashed_password_071', 'Employee Seventy-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000072-0000-4000-8000-000000000072', 'employee.072', 'hashed_password_072', 'Employee Seventy-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000073-0000-4000-8000-000000000073', 'employee.073', 'hashed_password_073', 'Employee Seventy-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000074-0000-4000-8000-000000000074', 'employee.074', 'hashed_password_074', 'Employee Seventy-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000075-0000-4000-8000-000000000075', 'employee.075', 'hashed_password_075', 'Employee Seventy-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000076-0000-4000-8000-000000000076', 'employee.076', 'hashed_password_076', 'Employee Seventy-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000077-0000-4000-8000-000000000077', 'employee.077', 'hashed_password_077', 'Employee Seventy-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000078-0000-4000-8000-000000000078', 'employee.078', 'hashed_password_078', 'Employee Seventy-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000079-0000-4000-8000-000000000079', 'employee.079', 'hashed_password_079', 'Employee Seventy-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000080-0000-4000-8000-000000000080', 'employee.080', 'hashed_password_080', 'Employee Eighty', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000081-0000-4000-8000-000000000081', 'employee.081', 'hashed_password_081', 'Employee Eighty-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000082-0000-4000-8000-000000000082', 'employee.082', 'hashed_password_082', 'Employee Eighty-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000083-0000-4000-8000-000000000083', 'employee.083', 'hashed_password_083', 'Employee Eighty-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000084-0000-4000-8000-000000000084', 'employee.084', 'hashed_password_084', 'Employee Eighty-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000085-0000-4000-8000-000000000085', 'employee.085', 'hashed_password_085', 'Employee Eighty-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000086-0000-4000-8000-000000000086', 'employee.086', 'hashed_password_086', 'Employee Eighty-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000087-0000-4000-8000-000000000087', 'employee.087', 'hashed_password_087', 'Employee Eighty-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000088-0000-4000-8000-000000000088', 'employee.088', 'hashed_password_088', 'Employee Eighty-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000089-0000-4000-8000-000000000089', 'employee.089', 'hashed_password_089', 'Employee Eighty-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000090-0000-4000-8000-000000000090', 'employee.090', 'hashed_password_090', 'Employee Ninety', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000091-0000-4000-8000-000000000091', 'employee.091', 'hashed_password_091', 'Employee Ninety-One', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000092-0000-4000-8000-000000000092', 'employee.092', 'hashed_password_092', 'Employee Ninety-Two', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000093-0000-4000-8000-000000000093', 'employee.093', 'hashed_password_093', 'Employee Ninety-Three', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000094-0000-4000-8000-000000000094', 'employee.094', 'hashed_password_094', 'Employee Ninety-Four', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000095-0000-4000-8000-000000000095', 'employee.095', 'hashed_password_095', 'Employee Ninety-Five', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000096-0000-4000-8000-000000000096', 'employee.096', 'hashed_password_096', 'Employee Ninety-Six', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000097-0000-4000-8000-000000000097', 'employee.097', 'hashed_password_097', 'Employee Ninety-Seven', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000098-0000-4000-8000-000000000098', 'employee.098', 'hashed_password_098', 'Employee Ninety-Eight', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000099-0000-4000-8000-000000000099', 'employee.099', 'hashed_password_099', 'Employee Ninety-Nine', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL),
('e0000100-0000-4000-8000-000000000100', 'employee.100', 'hashed_password_100', 'Employee One Hundred', NOW(), NOW(), 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', 'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', FALSE, NULL);

-- Ensure the uuid-ossp extension is enabled if not already done:
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO salaries (
    salary_id,
    employee_id,
    monthly_salary,
    created_at,
    updated_at,
    created_by,
    updated_by,
    is_deleted,
    deleted_time
)
SELECT
    uuid_generate_v4(), -- Generate a new UUID for each salary record
    e.employee_id,
    -- Generate a random monthly salary between 4000.00 and 6000.00
    FLOOR(4000 + (RANDOM() * 2000))::NUMERIC(12, 2),
    NOW(), -- Set creation timestamp to current time
    NOW(), -- Set update timestamp to current time
    'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', -- Replace with the actual AdminID used previously
    'a1a1a1a1-a1a1-4a1a-8a1a-a1a1a1a1a1a1', -- Replace with the actual AdminID used previously
    FALSE, -- Mark as not deleted
    NULL   -- No deletion time
FROM
    employees e
WHERE NOT EXISTS (
    SELECT 1
    FROM salaries s
    WHERE s.employee_id = e.employee_id AND s.is_deleted = FALSE
);