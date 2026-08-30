package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	ddpConnectMsg               = `{"msg":"connect","version":"1","support":["1"]}`
	wsEndpoint                  = "/websocket"
	wsTimeout                   = 30 * time.Second
	maximumIntermediaryMessages = 500
)

type webSocketClient struct {
	baseURL  string
	username string
	password string
}

type ddpMsg struct {
	Error  *ddpError       `json:"error,omitempty"`
	ID     string          `json:"id,omitempty"`
	Msg    string          `json:"msg"`
	Method string          `json:"method,omitempty"`
	Params []any           `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

type ddpError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type wsAppDTO struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	ActiveWorkloads wsWorkloadDTO `json:"active_workloads"`
}

type wsWorkloadDTO struct {
	UsedPorts []wsPortDTO `json:"used_ports"`
}

type wsPortDTO struct {
	Protocol      string          `json:"protocol"`
	HostPorts     []wsHostPortDTO `json:"host_ports"`
	ContainerPort int             `json:"container_port"`
}

type wsHostPortDTO struct {
	HostIP   string `json:"host_ip"`
	HostPort int    `json:"host_port"`
}

func newWebSocketClient(baseURL, username, password string) *webSocketClient {
	return &webSocketClient{
		baseURL:  baseURL,
		username: username,
		password: password,
	}
}

func (c *webSocketClient) GetAvailableApps() ([]AvailableAppDTO, error) {
	cacheKey := fmt.Sprintf("%s:%s:ws:app.query", c.baseURL, c.username)
	res, err := getFromCache(cacheKey, func() (*[]AvailableAppDTO, error) {
		res, err := c.fetchApps()
		if err != nil {
			return nil, err
		}

		return &res, nil
	})

	if err != nil || res == nil {
		return nil, err
	}

	return *res, nil
}

func (c *webSocketClient) fetchApps() ([]AvailableAppDTO, error) {
	wsURL, err := buildWSURL(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("truenas websocket: invalid base URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("truenas websocket: dial failed: %w", err)
	}

	//nolint:errcheck
	defer conn.Close()

	deadline := time.Now().Add(wsTimeout)
	if err = conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("truenas websocket: set read deadline: %w", err)
	}

	if err = conn.SetWriteDeadline(deadline); err != nil {
		return nil, fmt.Errorf("truenas websocket: set write deadline: %w", err)
	}

	if err = conn.WriteMessage(websocket.TextMessage, []byte(ddpConnectMsg)); err != nil {
		return nil, fmt.Errorf("truenas websocket: send connect: %w", err)
	}

	if err = c.waitForMsg(conn, "connected"); err != nil {
		return nil, fmt.Errorf("truenas websocket: handshake failed: %w", err)
	}

	authPayload, err := json.Marshal(ddpMsg{
		ID:     "1",
		Msg:    "method",
		Method: "auth.login",
		Params: []any{c.username, c.password},
	})
	if err != nil {
		return nil, fmt.Errorf("truenas websocket: marshal auth payload: %w", err)
	}

	if err = conn.WriteMessage(websocket.TextMessage, authPayload); err != nil {
		return nil, fmt.Errorf("truenas websocket: send auth.login: %w", err)
	}

	authResult, err := c.waitForResult(conn, "1")
	if err != nil {
		return nil, fmt.Errorf("truenas websocket: auth failed: %w", err)
	}

	var ok bool
	if err = json.Unmarshal(authResult, &ok); err != nil || !ok {
		return nil, errors.New("truenas websocket: authentication rejected")
	}

	queryPayload, err := json.Marshal(ddpMsg{
		ID:     "2",
		Msg:    "method",
		Method: "app.query",
		Params: []any{},
	})
	if err != nil {
		return nil, fmt.Errorf("truenas websocket: marshal app.query payload: %w", err)
	}

	if err = conn.WriteMessage(websocket.TextMessage, queryPayload); err != nil {
		return nil, fmt.Errorf("truenas websocket: send app.query: %w", err)
	}

	appsResult, err := c.waitForResult(conn, "2")
	if err != nil {
		return nil, fmt.Errorf("truenas websocket: app.query failed: %w", err)
	}

	var raw []wsAppDTO
	if err = json.Unmarshal(appsResult, &raw); err != nil {
		return nil, fmt.Errorf("truenas websocket: parse app.query result: %w", err)
	}

	return convertWSApps(raw), nil
}

func convertWSApps(raw []wsAppDTO) []AvailableAppDTO {
	apps := make([]AvailableAppDTO, len(raw))
	for index, item := range raw {
		ports := make([]WorkloadPortDTO, len(item.ActiveWorkloads.UsedPorts))

		for portIndex, port := range item.ActiveWorkloads.UsedPorts {
			hostPorts := make([]HostPortDTO, len(port.HostPorts))
			for hostPortIndex, hostPort := range port.HostPorts {
				hostPorts[hostPortIndex] = HostPortDTO(hostPort)
			}

			ports[portIndex] = WorkloadPortDTO{
				Protocol:      port.Protocol,
				ContainerPort: port.ContainerPort,
				HostPorts:     hostPorts,
			}
		}

		apps[index] = AvailableAppDTO{
			ID:              item.ID,
			Name:            item.Name,
			ActiveWorkloads: WorkloadDTO{UsedPorts: ports},
		}
	}

	return apps
}

func (c *webSocketClient) waitForMsg(conn *websocket.Conn, msgType string) error {
	remainingAttempts := maximumIntermediaryMessages
	for remainingAttempts > 0 {
		remainingAttempts--

		_, raw, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var msg ddpMsg
		if err = json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if msg.Msg == msgType {
			return nil
		}
	}

	return errors.New("truenas websocket: timeout")
}

func (c *webSocketClient) waitForResult(conn *websocket.Conn, id string) (json.RawMessage, error) {
	remainingAttempts := maximumIntermediaryMessages
	for remainingAttempts > 0 {
		remainingAttempts--

		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil, err
		}

		var msg ddpMsg
		if err = json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if msg.Msg == "result" && msg.ID == id {
			if msg.Error != nil {
				return nil, fmt.Errorf("%s: %s", msg.Error.Error, msg.Error.Message)
			}

			return msg.Result, nil
		}
	}

	return nil, errors.New("truenas websocket: timeout")
}

func buildWSURL(baseURL string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if parsed.Scheme == "https" {
		parsed.Scheme = "wss"
	} else if parsed.Scheme != "wss" {
		parsed.Scheme = "ws"
	}

	parsed.User = nil
	parsed.Path = wsEndpoint
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}
