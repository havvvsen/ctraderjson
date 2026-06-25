package ctrader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"

	"google.golang.org/protobuf/proto"
)

// type CtraderMessage interface{}

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
	Conn                    *net.Conn
	HandlerFunc             func(proto.Message)
}

func (c *Client) Start() error {
	var host string

	if c.Live {
		host = "live.ctraderapi.com:5036"
	} else {
		host = "demo.ctraderapi.com:5036"
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()

	c.Conn = &conn
	c.Logger.Info("Connected to ctrader successfully")

	reader := bufio.NewReader(*c.Conn)

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
	return nil
}

func (c *Client) keepAlive() {

}

func (c *Client) Send(ctx context.Context, conn net.Conn, msg Message) error {
	messageRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	conn.Write(messageRaw)

	return nil
}
