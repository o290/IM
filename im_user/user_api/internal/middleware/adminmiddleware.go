package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"server/common/response"
)

type AdminMiddleware struct {
}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("Role")
		fmt.Println(role)
		if role != "1" {
			response.Response(r, w, nil, errors.New("角色鉴权失败"))
			return
		}
		next(w, r)
	}
}
