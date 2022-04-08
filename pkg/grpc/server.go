package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/pb"
	"github.com/c9s/bbgo/pkg/types"
)

type Server struct {
	Config  *bbgo.Config
	Environ *bbgo.Environment
	Trader  *bbgo.Trader

	pb.UnimplementedMarketDataServiceServer
}

func (s *Server) Subscribe(request *pb.SubscribeRequest, server pb.MarketDataService_SubscribeServer) error {
	panic("implement me")
	return nil
}

func (s *Server) QueryKLines(ctx context.Context, request *pb.QueryKLinesRequest) (*pb.QueryKLinesResponse, error) {
	exchangeName, err := types.ValidExchangeName(request.Exchange)
	if err != nil {
		return nil, err
	}

	for _, session := range s.Environ.Sessions() {
		if session.ExchangeName == exchangeName {
			response := &pb.QueryKLinesResponse{
				Klines: nil,
				Error:  nil,
			}

			options := types.KLineQueryOptions{
				Limit: int(request.Limit),
			}

			klines, err := session.Exchange.QueryKLines(ctx, request.Symbol, types.Interval(request.Interval), options)
			if err != nil {
				return nil, err
			}

			for _, kline := range klines {
				response.Klines = append(response.Klines, &pb.KLine{
					Exchange:    kline.Exchange.String(),
					Symbol:      kline.Symbol,
					Timestamp:   kline.StartTime.Unix(),
					Open:        kline.Open.Float64(),
					High:        kline.High.Float64(),
					Low:         kline.Low.Float64(),
					Close:       kline.Close.Float64(),
					Volume:      kline.Volume.Float64(),
					QuoteVolume: kline.QuoteVolume.Float64(),
				})
			}

			return response, nil
		}
	}

	return nil, nil
}

func (s *Server) ListenAndServe(bind string) error {
	conn, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrapf(err, "failed to bind network at %s", bind)
	}

	var grpcServer = grpc.NewServer()
	pb.RegisterMarketDataServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(conn); err != nil {
		return errors.Wrap(err, "failed to serve grpc connections")
	}

	return nil
}