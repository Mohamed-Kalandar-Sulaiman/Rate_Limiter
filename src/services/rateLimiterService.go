package services

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/Proto/generated"

	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/repository"
	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/utils"
)

type RateLimiterServer struct {
	repo       *repository.RateLimiterRepository
	configMap  *utils.ConfigMap
	proto.UnimplementedRateLimitServiceServer
}

func NewRateLimiterServer(repo repository.RateLimiterRepository, configMap *utils.ConfigMap) *RateLimiterServer {
	return &RateLimiterServer{
		repo:      &repo,
		configMap: configMap,
	}
}


func (s *RateLimiterServer) GetApplicationLayerRateLimit(ctx context.Context, req *proto.RateLimitRequest) (*proto.RateLimitResponse, error) {
	key 		:= fmt.Sprintf("%s:%s:%s", req.GetServiceName(), req.GetActionName(), req.GetConfigName())
	config, err := s.configMap.GetConfig(key)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get config for key %s: %v", key, err)
	}

	factory := utils.NewRateLimiterFactory(s.repo)
	limiter, err := factory.CreateRateLimiter(config.RateLimitConfig.Algorithm)

	if err != nil{
		response := &proto.RateLimitResponse{
			IsAllowed   : false,
			Remaining   : 0,
			Limit       : 0,
			ResetTime   : 0,
			ResetAfter  : 0,
			ErrorCode   : 1,
			ErrorMessage: err.Error(),
		}
		return response, status.Error(codes.Internal, "Rate Limit error")
	}

	isAllowed , limit, remaining, reset , reset_after ,err := limiter.RateLimitFunction(		key,
																								config.RateLimitConfig.Unit,
																								config.RateLimitConfig.RequestPerUnit,
																								)
 	
    if err != nil{
		log.Fatal("unknown error occurred")
	}
	response := &proto.RateLimitResponse{
		IsAllowed   : isAllowed,
		Remaining   : int32(remaining),
		Limit       : int32(limit),
		ResetTime   : int32(reset),
		ResetAfter  : int32(reset_after),
		ErrorCode   : 0,
		ErrorMessage: "",
	}

	return response, nil
}

func (s *RateLimiterServer) GetHealth(ctx context.Context, req *proto.Void) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{Status: "OK"}, nil
}
