package interceptors

import (
	"context"
	"time"

	"bookService/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	statusCode := codes.OK.String()
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code().String()
		}
	}

	methodName := info.FullMethod
	metrics.GRPCRequestsTotal.WithLabelValues(methodName, statusCode).Inc()
	metrics.GRPCDuration.WithLabelValues(methodName).Observe(time.Since(start).Seconds())

	return resp, err
}
