# Developer Guide: payslip-service

This guide helps you get started with running the service, generating Swagger documentation, creating mocks, and running tests.

---

## 📁 Prerequisites

Make sure you have the following tools installed:

* Go (>= 1.18)
* [swag](https://github.com/swaggo/swag):

  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```
* [mockery](https://github.com/vektra/mockery):

  ```bash
  go install github.com/vektra/mockery/v2/...@latest
  ```
* Git

---

## ⚙️ Commands via `make`

The `Makefile` defines key automation tasks:

### ✅ Install Dependencies

```bash
make dep
```

> Downloads and syncs all Go module dependencies.

---

### 🚀 Run API Server

```bash
make run-api
```

> Starts the HTTP server using `go run main.go serve-http`.

---

### 📄 Generate Swagger Docs

```bash
make swag-init
```

> Scans all Go files for `@Summary`, `@Router`, etc., and generates Swagger documentation in the `docs/` folder.

To preview:

```bash
make run-api
# Then open http://localhost:<your-port>/swagger/index.html
```

Make sure to import the Swagger docs in your server:

```go
import _ "your/module/path/docs"
```

---

### 🧪 Run Tests

```bash
make run-test
```

> Runs all unit tests with:

* Verbose output
* Coverage enabled
* No test result caching
* Stops at first failure

---

### 🔁 Regenerate Mocks

```bash
make remock
```

> Uses `mockery` to generate mocks for interfaces.

Ensure your interfaces are annotated:

```go
//go:generate mockery --name=PayslipInteractor
```

Or run `mockery` manually with flags.

Example:

```go
//go:generate mockery --name=PayslipInteractor
type PayslipInteractor interface {
  Payroll(ctx context.Context, req *request.PayrollRequest) (*entity.PayrollResponse, error)
}
```

---

## 🔍 API Routes Overview

The service exposes these main endpoints under `/v1`:

| Endpoint                | Method | Purpose                  |
| ----------------------- | ------ | ------------------------ |
| `/v1/payroll/run`       | POST   | Run payroll              |
| `/v1/payroll/summary`   | POST   | Generate payroll summary |
| `/v1/payslip`           | POST   | Generate payslip         |
| `/v1/attendance`        | POST   | Record attendance        |
| `/v1/attendance-period` | POST   | Create attendance period |
| `/v1/overtime`          | POST   | Submit overtime          |
| `/v1/reimbursement`     | POST   | Submit reimbursement     |

---

## 📚 Validating Swagger

Ensure your handler functions are annotated:

```go
// RunPayroll godoc
// @Summary      Run payroll
// @Description  Admin runs payroll for a given attendance period
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        data  body request.PayrollRequest true "RunPayroll payload"
// @Success      200  {object} entity.PayrollResponse
// @Failure      400  {string} string "Bad request"
// @Failure      500  {string} string "Internal server error"
// @Router       /v1/payroll/run [post]
```

Then regenerate docs:

```bash
make swag-init
```

And preview:

```http
http://localhost:<your-port>/docs/index.html
```

---

Happy coding! 🚀
