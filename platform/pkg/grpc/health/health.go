package health

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type HealthChecker struct {
	mu       sync.RWMutex
	services map[string]grpchealth.HealthCheckResponse_ServingStatus
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		services: make(map[string]grpchealth.HealthCheckResponse_ServingStatus),
	}
}

func (h *HealthChecker) Check(ctx context.Context, req *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	serviceName := req.GetService()

	if serviceName == "" {
		return &grpchealth.HealthCheckResponse{
			Status: grpchealth.HealthCheckResponse_SERVING,
		}, nil
	}

	status, exists := h.services[serviceName]
	if !exists {
		return &grpchealth.HealthCheckResponse{
			Status: grpchealth.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}

	return &grpchealth.HealthCheckResponse{
		Status: status,
	}, nil
}

func (h *HealthChecker) Watch(req *grpchealth.HealthCheckRequest, stream grpchealth.Health_WatchServer) error {
	serviceName := req.GetService()

	h.mu.RLock()
	currentStatus, exists := h.services[serviceName]
	h.mu.RUnlock()

	if !exists && serviceName != "" {
		currentStatus = grpchealth.HealthCheckResponse_SERVICE_UNKNOWN
	} else if !exists {
		currentStatus = grpchealth.HealthCheckResponse_SERVING
	}

	resp := &grpchealth.HealthCheckResponse{
		Status: currentStatus,
	}

	if err := stream.Send(resp); err != nil {
		return err
	}

	ch := make(chan grpchealth.HealthCheckResponse_ServingStatus, 1)
	ch <- currentStatus

	for {
		select {
		case <-stream.Context().Done():
			return status.Errorf(codes.Canceled, "stream has ended")
		case newStatus := <-ch:
			if newStatus != currentStatus {
				currentStatus = newStatus
				resp := &grpchealth.HealthCheckResponse{
					Status: currentStatus,
				}
				if err := stream.Send(resp); err != nil {
					return err
				}
			}
		}
	}
}

func (h *HealthChecker) SetServing(serviceName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.services[serviceName] = grpchealth.HealthCheckResponse_SERVING
}

func (h *HealthChecker) SetNotServing(serviceName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.services[serviceName] = grpchealth.HealthCheckResponse_NOT_SERVING
}

func (h *HealthChecker) SetUnknown(serviceName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.services[serviceName] = grpchealth.HealthCheckResponse_SERVICE_UNKNOWN
}

func RegisterService(server *grpc.Server) {
	healthChecker := NewHealthChecker()
	grpchealth.RegisterHealthServer(server, healthChecker)
	healthChecker.SetServing("")
}
