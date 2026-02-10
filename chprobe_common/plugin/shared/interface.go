package shared

import (
	"context"

	"github.com/ricky97gr/chprobe/chprobe_common/plugin/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type PluginService interface {
	GetPluginInfo(ctx context.Context) (*PluginInfo, error)
	GetRoutes(ctx context.Context) ([]*RouteInfo, error)
	HandleRequest(ctx context.Context, method string, path string, data map[string]string) (*HandleResult, error)
	GetMetrics(ctx context.Context) (map[string]string, error)
	GetMetricsSummary(ctx context.Context) (map[string]string, error)
	HealthCheck(ctx context.Context) (*HealthStatus, error)
}

type PluginInfo struct {
	Id          string
	Name        string
	Version     string
	Author      string
	Description string
}

type RouteInfo struct {
	Path        string
	Method      string
	Description string
	Handler     string
	Metadata    map[string]string
}

type HandleResult struct {
	Success bool
	Data    map[string]string
	Error   string
}

type HealthStatus struct {
	Status      string
	Name        string
	Version     string
	Description string
	Metadata    map[string]string
}

type GRPCPlugin struct {
	plugin.Plugin
	Impl PluginService
}

func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterPluginServiceServer(server, &GRPCServerWrapper{impl: p.Impl})
	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &GRPCClientWrapper{client: proto.NewPluginServiceClient(conn)}, nil
}

type GRPCServerWrapper struct {
	impl PluginService
	proto.UnimplementedPluginServiceServer
}

func (s *GRPCServerWrapper) GetPluginInfo(ctx context.Context, req *proto.GetPluginInfoRequest) (*proto.GetPluginInfoResponse, error) {
	info, err := s.impl.GetPluginInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.GetPluginInfoResponse{
		Id:          info.Id,
		Name:        info.Name,
		Version:     info.Version,
		Author:      info.Author,
		Description: info.Description,
	}, nil
}

func (s *GRPCServerWrapper) GetRoutes(ctx context.Context, req *proto.GetRoutesRequest) (*proto.GetRoutesResponse, error) {
	routes, err := s.impl.GetRoutes(ctx)
	if err != nil {
		return nil, err
	}
	protoRoutes := make([]*proto.RouteInfo, len(routes))
	for i, r := range routes {
		protoRoutes[i] = &proto.RouteInfo{
			Path:        r.Path,
			Method:      r.Method,
			Description: r.Description,
			Handler:     r.Handler,
			Metadata:    r.Metadata,
		}
	}
	return &proto.GetRoutesResponse{Routes: protoRoutes}, nil
}

func (s *GRPCServerWrapper) HandleRequest(ctx context.Context, req *proto.HandleRequestRequest) (*proto.HandleRequestResponse, error) {
	result, err := s.impl.HandleRequest(ctx, req.Method, req.Path, req.Data)
	if err != nil {
		return nil, err
	}
	return &proto.HandleRequestResponse{
		Success: result.Success,
		Data:    result.Data,
		Error:   result.Error,
	}, nil
}

func (s *GRPCServerWrapper) GetMetrics(ctx context.Context, req *proto.GetMetricsRequest) (*proto.GetMetricsResponse, error) {
	metrics, err := s.impl.GetMetrics(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.GetMetricsResponse{Metrics: metrics}, nil
}

func (s *GRPCServerWrapper) GetMetricsSummary(ctx context.Context, req *proto.GetMetricsSummaryRequest) (*proto.GetMetricsSummaryResponse, error) {
	summary, err := s.impl.GetMetricsSummary(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.GetMetricsSummaryResponse{Summary: summary}, nil
}

func (s *GRPCServerWrapper) HealthCheck(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	status, err := s.impl.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.HealthCheckResponse{
		Status:      status.Status,
		Name:        status.Name,
		Version:     status.Version,
		Description: status.Description,
		Metadata:    status.Metadata,
	}, nil
}

type GRPCClientWrapper struct {
	client proto.PluginServiceClient
}

func (c *GRPCClientWrapper) GetPluginInfo(ctx context.Context) (*PluginInfo, error) {
	resp, err := c.client.GetPluginInfo(ctx, &proto.GetPluginInfoRequest{})
	if err != nil {
		return nil, err
	}
	return &PluginInfo{
		Id:          resp.Id,
		Name:        resp.Name,
		Version:     resp.Version,
		Author:      resp.Author,
		Description: resp.Description,
	}, nil
}

func (c *GRPCClientWrapper) GetRoutes(ctx context.Context) ([]*RouteInfo, error) {
	resp, err := c.client.GetRoutes(ctx, &proto.GetRoutesRequest{})
	if err != nil {
		return nil, err
	}
	routes := make([]*RouteInfo, len(resp.Routes))
	for i, r := range resp.Routes {
		routes[i] = &RouteInfo{
			Path:        r.Path,
			Method:      r.Method,
			Description: r.Description,
			Handler:     r.Handler,
			Metadata:    r.Metadata,
		}
	}
	return routes, nil
}

func (c *GRPCClientWrapper) HandleRequest(ctx context.Context, method string, path string, data map[string]string) (*HandleResult, error) {
	resp, err := c.client.HandleRequest(ctx, &proto.HandleRequestRequest{
		Method: method,
		Path:   path,
		Data:   data,
	})
	if err != nil {
		return nil, err
	}
	return &HandleResult{
		Success: resp.Success,
		Data:    resp.Data,
		Error:   resp.Error,
	}, nil
}

func (c *GRPCClientWrapper) GetMetrics(ctx context.Context) (map[string]string, error) {
	resp, err := c.client.GetMetrics(ctx, &proto.GetMetricsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}

func (c *GRPCClientWrapper) GetMetricsSummary(ctx context.Context) (map[string]string, error) {
	resp, err := c.client.GetMetricsSummary(ctx, &proto.GetMetricsSummaryRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Summary, nil
}

func (c *GRPCClientWrapper) HealthCheck(ctx context.Context) (*HealthStatus, error) {
	resp, err := c.client.HealthCheck(ctx, &proto.HealthCheckRequest{})
	if err != nil {
		return nil, err
	}
	return &HealthStatus{
		Status:      resp.Status,
		Name:        resp.Name,
		Version:     resp.Version,
		Description: resp.Description,
		Metadata:    resp.Metadata,
	}, nil
}
