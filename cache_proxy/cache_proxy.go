package cache_proxy

import (
	"context"
	"io"

	"google.golang.org/grpc"

	repb "github.com/buildbuddy-io/buildbuddy/proto/remote_execution"
	bspb "google.golang.org/genproto/googleapis/bytestream"
)

const (
	// Keep under the limit of ~4MB (1024 * 1024 * 4).
	readBufSizeBytes = (1024 * 1024 * 4) - 100
)

type CacheProxy struct {
	acClient  repb.ActionCacheClient
	bsClient  bspb.ByteStreamClient
	casClient repb.ContentAddressableStorageClient
	cpbClient repb.CapabilitiesClient
}

func NewCacheProxy(conn *grpc.ClientConn) (*CacheProxy, error) {
	return &CacheProxy{
		acClient:  repb.NewActionCacheClient(conn),
		bsClient:  bspb.NewByteStreamClient(conn),
		casClient: repb.NewContentAddressableStorageClient(conn),
		cpbClient: repb.NewCapabilitiesClient(conn),
	}, nil
}

func (p *CacheProxy) GetCapabilities(ctx context.Context, req *repb.GetCapabilitiesRequest) (*repb.ServerCapabilities, error) {
	return p.cpbClient.GetCapabilities(ctx, req)
}

func (p *CacheProxy) GetActionResult(ctx context.Context, req *repb.GetActionResultRequest) (*repb.ActionResult, error) {
	return p.acClient.GetActionResult(ctx, req)
}

func (p *CacheProxy) UpdateActionResult(ctx context.Context, req *repb.UpdateActionResultRequest) (*repb.ActionResult, error) {
	return p.acClient.UpdateActionResult(ctx, req)
}

func (p *CacheProxy) BatchUpdateBlobs(ctx context.Context, req *repb.BatchUpdateBlobsRequest) (*repb.BatchUpdateBlobsResponse, error) {
	return p.casClient.BatchUpdateBlobs(ctx, req)
}

func (p *CacheProxy) BatchReadBlobs(ctx context.Context, req *repb.BatchReadBlobsRequest) (*repb.BatchReadBlobsResponse, error) {
	return p.casClient.BatchReadBlobs(ctx, req)
}

func (p *CacheProxy) GetTree(req *repb.GetTreeRequest, stream repb.ContentAddressableStorage_GetTreeServer) error {
	clientStream, err := p.casClient.GetTree(stream.Context(), req)
	if err != nil {
		return err
	}
	for {
		msg, err := clientStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (p *CacheProxy) FindMissingBlobs(ctx context.Context, req *repb.FindMissingBlobsRequest) (*repb.FindMissingBlobsResponse, error) {
	return p.casClient.FindMissingBlobs(ctx, req)
}

func (p *CacheProxy) Read(req *bspb.ReadRequest, stream bspb.ByteStream_ReadServer) error {
	clientStream, err := p.bsClient.Read(stream.Context(), req)
	if err != nil {
		return err
	}
	for {
		msg, err := clientStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if sendErr := stream.Send(msg); sendErr != nil {
			return sendErr
		}
	}
	return nil
}

func (p *CacheProxy) Write(stream bspb.ByteStream_WriteServer) error {
	clientStream, err := p.bsClient.Write(stream.Context())
	if err != nil {
		return err
	}
	for {
		rsp, err := stream.Recv()
		if err != nil {
			return err
		}
		if err := clientStream.Send(rsp); err != nil {
			return err
		}
		if rsp.GetFinishWrite() {
			lastRsp, err := clientStream.CloseAndRecv()
			if err != nil {
				return err
			}
			return stream.SendAndClose(lastRsp)
		}

	}
	return nil
}

func (p *CacheProxy) QueryWriteStatus(ctx context.Context, req *bspb.QueryWriteStatusRequest) (*bspb.QueryWriteStatusResponse, error) {
	// For now, just tell the client that the entire write failed and let
	// them retry it.
	return &bspb.QueryWriteStatusResponse{
		CommittedSize: 0,
		Complete:      false,
	}, nil
}
