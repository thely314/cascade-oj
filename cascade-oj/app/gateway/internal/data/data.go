package data

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewJudgeRepo, NewContestRepo, NewMiscRepo, NewGrpcUserClient)

// Data
type Data struct {
	cache          *Cache
	mq_channel     *amqp.Channel
	grpcUserClient pb.UserClient
}

// Cache store maps
type Cache struct {
	// TODO cache fields
	token string
}

func NewData(c *conf.Data, grpcUserClient pb.UserClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	cache := &Cache{}

	// register gateway endpoint to be called by other services
	var endpoint string
	if c.EndpointPreset != "" {
		endpoint = c.EndpointPreset
	} else {
		// get available gateway ipv4 address
		ip, err := getIPv4(byte(c.Ipv4Prefix))
		if err != nil {
			return nil, nil, err
		}
		endpoint = ip + ":" + strconv.FormatInt(c.Port, 10)
	}
	registeredGateway, err := grpcUserClient.Register(context.Background(), &pb.RegisterRequest{
		Endpoint: endpoint,
		ApiKey:   c.ApiKey,
	})
	if err != nil {
		return nil, nil, err
	}
	// token is used to identify which service the request comes from
	cache.token = registeredGateway.Token
	// update cache if needed

	// connect to mq
	mq_connection, err := amqp.Dial(c.Mq)
	if err != nil {
		log.Errorf("failed to connect to mq: %v", err)
	}
	mq_ch, err := mq_connection.Channel()
	if err != nil {
		log.Errorf("failed to open a channel: %v", err)
	}
	return &Data{
		cache:          cache,
		mq_channel:     mq_ch,
		grpcUserClient: grpcUserClient,
	}, cleanup, nil
}

func NewGrpcUserClient(c *conf.Client) pb.UserClient {
	connection, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Grpc.Addr),
	)
	if err != nil {
		panic(err)
	}
	grpcUserClient := pb.NewUserClient(connection)
	return grpcUserClient
}

// get the current client's ipv4 address
func getIPv4(ipv4Prefix byte) (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", errors.New("cannot get system IPv4 address")
	}
	for _, addr := range addrs {
		possibleNet, ok := addr.(*net.IPNet)
		if ok && !possibleNet.IP.IsLoopback() {
			if possibleNet.IP.To4() != nil && possibleNet.IP[0] == ipv4Prefix {
				return possibleNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("get %s, but IPv4 start with %d not found", addrs, ipv4Prefix)
}
