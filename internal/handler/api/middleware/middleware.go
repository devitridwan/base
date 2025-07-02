package middleware

// type Handle func(*http.Request) *response.JSONResponse

// type MiddlewareInterface interface {
// 	Auth(handle Handle) Handle
// }

// func (m *middleware) Auth(handle Handle) Handle {
// 	return func(r *http.Request) *response.JSONResponse {
// 		ctx := r.Context()

// 		userId, _ := strconv.ParseInt(r.Header.Get("x-user-id"), 10, 64)
// 		if userId == 0 {
// 			return response.CustomErrorResponse(err error)
// 		}
// 		ctx = context.WithValue(ctx, constants.ContextUserId, userId)

// 		branchId, _ := strconv.ParseInt(r.Header.Get("x-branch-id"), 10, 64)
// 		if branchId == 0 {
// 			return custresp.CustomErrorResponse(&custerr.ErrChain{
// 				Message: "invalid x-branch-id",
// 				Type:    custresp.ErrBadRequest,
// 			})
// 		}
// 		ctx = context.WithValue(ctx, constants.ContextBranchId, branchId)

// 		return handle(r.WithContext(ctx))
// 	}
// }
