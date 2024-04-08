package libp2p_test

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	streammsgv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/streammsg/v1"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p/libp2p"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type bufferStream struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

func newBufferStream() *bufferStream {
	return &bufferStream{
		buf: bytes.NewBuffer(nil),
	}
}

func (bs *bufferStream) Read(b []byte) (n int, err error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	return bs.buf.Read(b)
}

func (bs *bufferStream) Write(b []byte) (n int, err error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	return bs.buf.Write(b)
}

func (bs *bufferStream) Close() error {
	return nil
}

func (bs *bufferStream) Reset() error {
	bs.buf = bytes.NewBuffer(nil)
	return nil
}

func TestStream_ReadMsg_WriteMsg(t *testing.T) {
	// Prepare a StreamMsg
	msg := &streammsgv1.StreamMsg{
		Body: &streammsgv1.StreamMsg_Data{
			Data: []byte("test message"),
		},
	}

	// Prepare a bufferStream
	bs := newBufferStream()

	// Create a new stream
	s := libp2p.NewStream(bs, nil, nil)

	// Prepare a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call WriteMsg
	err := s.WriteMsg(ctx, msg)

	// Assert no error
	assert.NoError(t, err)

	// Prepare a proto.Message to hold the result
	result := &streammsgv1.StreamMsg{}

	// Call ReadMsg
	err = s.ReadMsg(ctx, result)

	// Assert no error
	assert.NoError(t, err)

	// Assert the message was correctly read
	assert.Equal(t, msg.GetData(), result.GetData())
}

func TestStream_ReadMsg_WriteMsg_ContextCancellation(t *testing.T) {
	// Prepare a StreamMsg
	msg := &streammsgv1.StreamMsg{
		Body: &streammsgv1.StreamMsg_Data{
			Data: []byte("test message"),
		},
	}

	// Prepare a bufferStream
	bs := newBufferStream()

	// Create a new stream
	s := libp2p.NewStream(bs, nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// Call WriteMsg
	err := s.WriteMsg(ctx, msg)

	// Assert context cancellation error
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	// Prepare a proto.Message to hold the result
	result := &streammsgv1.StreamMsg{}

	// Call ReadMsg with a cancelled context
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	err = s.ReadMsg(ctx, result)

	// Assert context cancellation error
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestStream_ReadMsg_InvalidMessage(t *testing.T) {
	// Prepare a bufferStream
	bs := newBufferStream()

	// Create a new stream
	s := libp2p.NewStream(bs, nil, nil)

	// Prepare a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Prepare a StreamMsg
	msg := &streammsgv1.StreamMsg{
		Body: &streammsgv1.StreamMsg_Data{
			Data: []byte("test message"),
		},
	}
	err := s.WriteMsg(ctx, msg)
	assert.NoError(t, err)

	// Prepare a different proto.Message to hold the result
	result := &streammsgv1.Header{}

	// Call ReadMsg
	err = s.ReadMsg(ctx, result)

	// Assert error
	assert.Error(t, err)
}

func TestMetadataStream_ReadHeader_WriteHeader(t *testing.T) {
	// Prepare a Header
	header := map[string]*structpb.Value{
		"test": structpb.NewStringValue("test"),
	}

	// Prepare a bufferStream
	bs := newBufferStream()

	// Create a new stream
	s := libp2p.NewMetadataStream(bs)

	// Prepare a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call WriteHeader
	err := s.WriteHeader(ctx, header)

	// Assert no error
	assert.NoError(t, err)

	// Call ReadHeader
	result, err := s.ReadHeader(ctx)

	// Assert no error
	assert.NoError(t, err)

	val1 := header["test"].GetStringValue()
	val2 := result["test"].GetStringValue()

	assert.Equal(t, val1, val2)
}

func TestMetadataStream_ReadHeader_WriteHeader_ContextCancellation(t *testing.T) {
	// Prepare a Header
	header := map[string]*structpb.Value{
		"test": structpb.NewStringValue("test"),
	}

	// Prepare a bufferStream
	bs := newBufferStream()

	// Create a new stream
	s := libp2p.NewMetadataStream(bs)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// Call WriteHeader
	err := s.WriteHeader(ctx, header)

	// Assert context cancellation error
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	// Call ReadHeader with a cancelled context
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	_, err = s.ReadHeader(ctx)

	// Assert context cancellation error
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestMetadataStream_WriteError(t *testing.T) {
	// Prepare a bufferStream
	bs := newBufferStream()

	// Create a new stream
	s := libp2p.NewMetadataStream(bs)

	// Prepare a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Prepare an error
	err := s.WriteError(ctx, status.New(codes.Internal, "test error"))
	assert.NoError(t, err)

	dataStream := libp2p.NewStream(bs, nil, nil)
	err = dataStream.ReadMsg(ctx, &streammsgv1.StreamMsg{})
	assert.Error(t, err)

	assert.Equal(t, codes.Internal, status.Convert(err).Code())
	assert.Equal(t, "test error", status.Convert(err).Message())
}
