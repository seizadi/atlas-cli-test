package main

import (
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/seizadi/atlas-cli-test/pkg/pb"
	"github.com/spf13/viper"
)

func RegisterGatewayEndpoints() gateway.Option {
	return gateway.WithEndpointRegistration(viper.GetString("server.version"),
		pb.RegisterAccountsHandlerFromEndpoint,
		pb.RegisterUsersHandlerFromEndpoint,
		pb.RegisterGroupsHandlerFromEndpoint,
	)
}
