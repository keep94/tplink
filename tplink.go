// Package tplink controls TP-Link plugs. 
package tplink

import (
  "fmt"
  "net"
)

const kPort = 9999

const (
  kOn = `{"system":{"set_relay_state":{"state":1}}}`
  kOff = `{"system":{"set_relay_state":{"state":0}}}`
)

// TPLink represents a single TP-Link plug
type TPLink struct {
  // Required: IP address of plug
  IP string
  // Optional: Port number. Default is 9999.
  Port int
}

// On turns on this plug
func (t *TPLink) On() error {
  return t.send(kOn)
}

// Off turns off this plug
func (t *TPLink) Off() error {
  return t.send(kOff)
}

func (t *TPLink) send(cmd string) error {
  port := t.Port
  if port == 0 {
    port = kPort
  }
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.IP, port))
  if err != nil {
    return err
  }
  defer conn.Close()
  payload := encrypt(cmd)
  if _, err := conn.Write(payload); err != nil {
    return err
  }
  recvBytes := make([]byte, 2048)
  if _, err := conn.Read(recvBytes); err != nil {
    return err
  }
  return nil
}

func encrypt(cmd string) []byte {
  plain := ([]byte)(cmd)
  // Leading 4 zeros
  result := make([]byte, len(plain) + 4)
  payload := result[4:]
  key := byte(171)
  for i := range plain {
    a := key ^ plain[i]
    key = a
    payload[i] = a
  }
  return result
}
