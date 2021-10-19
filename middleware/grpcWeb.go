package middleware

import (
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type GrpcWebMiddleware struct {
	*grpcweb.WrappedGrpcServer
}

func (m *GrpcWebMiddleware) Handler() http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		fmt.Printf("Request: %v\n", r)
		m.ServeHTTP(w, r)
		if m.IsAcceptableGrpcCorsRequest(r) || m.IsGrpcWebRequest(r) {
			m.ServeHTTP(w, r)
			fmt.Printf("Is GRPC Web Request: %v", "True")
			return
		} else {
			fmt.Printf("Is GRPC Web Request: %v", "False")
		}

	})
	server := http2.Server{}
	return h2c.NewHandler(handler, &server)
}

func NewGrpcWebMiddleware(grpcWeb *grpcweb.WrappedGrpcServer) *GrpcWebMiddleware {
	return &GrpcWebMiddleware{grpcWeb}
}
