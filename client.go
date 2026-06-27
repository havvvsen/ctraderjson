package ctrader

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type Message struct {
	ClientMsgId string `json:"clientMsgId"`
	PayloadType int    `json:"payloadType"`
	Payload     string `json:"payload"`
}


type Client struct {
	ApplicationClientId     string
	ApplicationClientSecret string
	AccountId               int64
	AccessToken             string
	RefreshToken            string
	Logger                  *slog.Logger
	Live                    bool
	Conn                    net.Conn
	HandlerFunc             func(proto.Message)
	Deadline                time.Duration
}

func (c *Client) Start() error {
	var host string

	if c.Live {
		host = "live.ctraderapi.com:5036"
	} else {
		host = "demo.ctraderapi.com:5036"
	}

	conn, err := net.DialTimeout("tcp", host, c.Deadline)
	if err != nil {
		return err
	}

	c.Conn = conn
	c.Logger.Info("Connected to ctrader successfully")

	reader := bufio.NewReader(c.Conn)

	for {
		msg, err := reader.ReadBytes('}')

		if err != nil {
			if err == io.EOF {
				c.Logger.Error("Error EOF")

				return nil

			}
			return err
		}
		c.Logger.Info(fmt.Sprintf("Message: %s\n", string(msg)))
	}

}

func (c *Client) Stop() error {
	if err := c.Conn.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) keepAlive() {
	ticker := time.Tick(time.Second * 10)
	protoHeartBeatEvent :=Message {
		ClientMsgId: fmt.Sprintf("id-%d", rand.Intn(492)),
		PayloadType: ,
	}

	for _ = range ticker {
		c.Send(context.Background())

	}

}

func (c *Client) Send(ctx context.Context, msg Message) error {
	messageRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.Conn.Write(messageRaw)

	return nil
