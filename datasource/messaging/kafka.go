package messaging

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/dewidyabagus/monorepo/shared-libs/mstring"
	kafka "github.com/segmentio/kafka-go"
)

type Kafka struct {
	Host             string
	RebalanceTimeout int
	MinBytes         int
	MaxBytes         int
	MaxAttempts      int
	MaxWait          int
}

// Indicator to reopen session
const (
	sessionClosed = "closed network connection"
	ioTimeout     = "i/o timeout"
	brokenPipe    = "write: broken pipe"
	reconnect     = "reconnect"
)
const (
	openConnMaxWait = 500 * time.Millisecond
	maxAttempt      = 3
)

type (
	Reader       = kafka.Reader
	ReaderConfig = kafka.ReaderConfig
	KafkaMessage = kafka.Message
)

type KafkaConn struct {
	*kafka.Conn
	mx              *sync.Mutex
	proto           string        // Protocol
	address         string        // Kafka address for ex: 127.0.0.1:9002 (host:port)
	openConnMaxWait time.Duration // Default 500 millisecond
	maxAttempt      int           // Maximum number of trials
}

func (k *KafkaConn) SetMaxAttempt(n int) {
	k.maxAttempt = n
}

func (k *KafkaConn) SetOpenConnMaxWait(m int) {
	k.openConnMaxWait = time.Duration(m) * time.Millisecond
}

func (k *KafkaConn) Reconnect(ctx context.Context, err error) {
	if err == nil {
		return
	}

	if mstring.ContainsOneOf(err.Error(), sessionClosed, ioTimeout, reconnect, brokenPipe) {
		if k.Conn != nil {
			k.Close()
		}

		if conn, err := NewKafkaConn(ctx, k.proto, k.address); err == nil {
			k.Conn = conn.Conn
		}
	}
}

func (k *KafkaConn) PingContext(ctx context.Context) (err error) {
	k.mx.Lock()
	defer k.mx.Unlock()

	if k.Conn == nil {
		k.Reconnect(ctx, errors.New(reconnect))
		return errors.New("kafka connection not established")
	}

	if deadline, ok := ctx.Deadline(); ok {
		k.SetDeadline(deadline)
	}

	for i := 0; i < k.maxAttempt; i++ {
		select {
		case <-ctx.Done():
			return

		default:
			err = func() error {
				if brokers, errF := k.Brokers(); errF != nil {
					ctxWT, cancel := context.WithTimeout(context.Background(), k.openConnMaxWait)
					defer cancel()

					k.Reconnect(ctxWT, errF) // Reopen connection
					return errF

				} else if len(brokers) == 0 {
					return errors.New("broker list is empty")
				}

				return nil
			}()
			if err == nil {
				return
			}
		}
	}

	return
}

func NewReader(config ReaderConfig) *Reader {
	return kafka.NewReader(config)
}

func NewKafkaReader(cfg *Kafka, groupID, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:          strings.Split(cfg.Host, ","),
		Topic:            topic,
		GroupID:          groupID,
		MinBytes:         cfg.MinBytes,
		MaxBytes:         cfg.MaxBytes,
		MaxWait:          time.Duration(cfg.MaxWait) * time.Second,
		MaxAttempts:      cfg.MaxAttempts,
		RebalanceTimeout: time.Duration(cfg.RebalanceTimeout) * time.Second,
	})
}

func NewKafkaPublisher(cfg *Kafka, topic string) *kafka.Writer {
	//Read and write timeout default 10s
	return &kafka.Writer{
		Addr:         kafka.TCP(strings.Split(cfg.Host, ",")...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		MaxAttempts:  3,
		RequiredAcks: 1,
		Async:        true,
	}
}

func NewKafkaMultipleHostConns(network, hosts string) (conns []*KafkaConn, err error) {
	address := strings.Split(strings.TrimSpace(hosts), ",")
	if len(address) == 0 {
		return nil, errors.New("kafka host is empty")
	}

	for _, host := range address {
		conn, _ := NewKafkaConn(context.Background(), network, host)
		conns = append(conns, conn)
	}

	return
}

func NewKafkaConn(ctx context.Context, network, address string) (*KafkaConn, error) {
	conn, err := kafka.DialContext(ctx, network, address)

	return &KafkaConn{
		Conn:            conn,
		mx:              new(sync.Mutex),
		proto:           network,
		address:         address,
		openConnMaxWait: openConnMaxWait,
		maxAttempt:      maxAttempt,
	}, err
}
