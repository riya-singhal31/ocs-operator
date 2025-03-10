package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	pb "github.com/red-hat-storage/ocs-operator/services/provider/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type OCSProviderClient struct {
	Client     pb.OCSProviderClient
	clientConn *grpc.ClientConn
	timeout    time.Duration
}

// NewProviderClient creates a client to talk to the external OCS storage provider server
func NewProviderClient(ctx context.Context, serverAddr string, timeout time.Duration) (*OCSProviderClient, error) {
	apiCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
	}

	conn, err := grpc.DialContext(apiCtx, serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return &OCSProviderClient{
		Client:     pb.NewOCSProviderClient(conn),
		clientConn: conn,
		timeout:    timeout}, nil
}

// Close closes the gRPC connection of the external OCS storage provider client
func (cc *OCSProviderClient) Close() {
	if cc.clientConn != nil {
		_ = cc.clientConn.Close()
		cc.clientConn = nil
	}
	cc.Client = nil
}

// OnboardConsumer to validate the consumer and create StorageConsumer
// resource on the StorageProvider cluster
func (cc *OCSProviderClient) OnboardConsumer(ctx context.Context, ticket, name, capacity string) (*pb.OnboardConsumerResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.OnboardConsumerRequest{
		OnboardingTicket: ticket,
		ConsumerName:     name,
		Capacity:         capacity,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.OnboardConsumer(apiCtx, req)
}

// GetStorageConfig generates the json config for connecting to storage provider cluster
func (cc *OCSProviderClient) GetStorageConfig(ctx context.Context, consumerUUID string) (*pb.StorageConfigResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.StorageConfigRequest{
		StorageConsumerUUID: consumerUUID,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.GetStorageConfig(apiCtx, req)
}

// OffboardConsumer deletes the StorageConsumer CR on the storage provider cluster
func (cc *OCSProviderClient) OffboardConsumer(ctx context.Context, consumerUUID string) (*pb.OffboardConsumerResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.OffboardConsumerRequest{
		StorageConsumerUUID: consumerUUID,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.OffboardConsumer(apiCtx, req)
}

// UpdateCapacity increases or decreases the storage block pool size
func (cc *OCSProviderClient) UpdateCapacity(ctx context.Context, consumerUUID, capacity string) (*pb.UpdateCapacityResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.UpdateCapacityRequest{
		StorageConsumerUUID: consumerUUID,
		Capacity:            capacity,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.UpdateCapacity(apiCtx, req)
}

func (cc *OCSProviderClient) AcknowledgeOnboarding(ctx context.Context, consumerUUID string) (*pb.AcknowledgeOnboardingResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.AcknowledgeOnboardingRequest{
		StorageConsumerUUID: consumerUUID,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.AcknowledgeOnboarding(apiCtx, req)
}

func (cc *OCSProviderClient) FulfillStorageClassClaim(ctx context.Context, consumerUUID, storageClassClaimName, encryptionMethod string, storageType pb.FulfillStorageClassClaimRequest_StorageType) (*pb.FulfillStorageClassClaimResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.FulfillStorageClassClaimRequest{
		StorageConsumerUUID:   consumerUUID,
		StorageClassClaimName: storageClassClaimName,
		EncryptionMethod:      encryptionMethod,
		StorageType:           storageType,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.FulfillStorageClassClaim(apiCtx, req)
}

func (cc *OCSProviderClient) RevokeStorageClassClaim(ctx context.Context, consumerUUID, storageClassClaimName string) (*pb.RevokeStorageClassClaimResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.RevokeStorageClassClaimRequest{
		StorageConsumerUUID:   consumerUUID,
		StorageClassClaimName: storageClassClaimName,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.RevokeStorageClassClaim(apiCtx, req)
}

func (cc *OCSProviderClient) GetStorageClassClaimConfig(ctx context.Context, consumerUUID, storageClassClaimName string) (*pb.StorageClassClaimConfigResponse, error) {
	if cc.Client == nil || cc.clientConn == nil {
		return nil, fmt.Errorf("provider client is closed")
	}

	req := &pb.StorageClassClaimConfigRequest{
		StorageConsumerUUID:   consumerUUID,
		StorageClassClaimName: storageClassClaimName,
	}

	apiCtx, cancel := context.WithTimeout(ctx, cc.timeout)
	defer cancel()

	return cc.Client.GetStorageClassClaimConfig(apiCtx, req)
}
