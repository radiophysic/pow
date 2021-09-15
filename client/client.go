package client

import (
	"bufio"
	"encoding/gob"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"tcp/conf"
	"tcp/pkg/lib"
	"tcp/service"
)

type Client struct {
	cfg conf.Config
	log *log.Logger
}

func NewClient(config conf.Config, logger *log.Logger) Client {
	return Client{cfg: config, log: logger}
}

func (c Client) GetQuote() error {
	serverAddr := c.cfg.ServerAddr + c.cfg.ServerPort
	rw, err := lib.Open(serverAddr)
	if err != nil {
		c.log.Println("The client cannot link to change the address:" + serverAddr)
		return err
	}

	var challengeRequestID, respRequestID, challengeID string
	if challengeRequestID, err = c.challengeRequest(rw); err != nil {
		return err
	}

	if challengeID, respRequestID, err = c.challengeResponse(rw); err != nil {
		return err
	}

	if challengeRequestID != respRequestID {
		return errors.New("challengeRequestID != respRequestID")
	}

	var hash []byte
	var nonces []uint32
	if hash, nonces, err = lib.Work(); err != nil {
		return errors.Wrap(err, "Proof-of-work failed.")
	}

	pow := &service.VerifyRequest{
		RequestID:   uuid.New().String(),
		ChallengeID: challengeID,
		Payload: struct {
			Hash   []byte
			Nonces []uint32
		}{Hash: hash, Nonces: nonces},
	}

	if _, err = c.verifyRequest(rw, pow); err != nil {
		return err
	}

	message, lastRequestID, err := c.verifyResponse(rw)
	if err != nil {
		return err
	}

	if pow.RequestID != lastRequestID {
		return errors.New("pow.RequestID != lastRequestID")
	}

	c.log.Printf("\n\nMy work has been proved\nQuote: [%s]\n\n", message)

	return nil
}

func (c Client) challengeRequest(rw *bufio.ReadWriter) (id string, err error) {
	challengeReq := service.ChallengeRequest{RequestID: uuid.New().String()}
	c.log.Printf("Ask server for challenge: RequestID: %s\n", challengeReq.RequestID)
	enc := gob.NewEncoder(rw)
	n, err := rw.WriteString("challenge\n")
	if err != nil {
		return challengeReq.RequestID, errors.Wrap(err, "Could not write ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(challengeReq)
	if err != nil {
		return challengeReq.RequestID, errors.Wrapf(err, "Encode failed for: %#v", challengeReq)
	}
	err = rw.Flush()
	if err != nil {
		return challengeReq.RequestID, errors.Wrap(err, "Flush failed.")
	}

	return challengeReq.RequestID, nil
}

func (c Client) challengeResponse(rw *bufio.ReadWriter) (challengeID string, requestID string, err error) {
	resp := &service.ChallengeResponse{}
	decoder := gob.NewDecoder(rw)
	if err = decoder.Decode(resp); err != nil {
		return "", "", errors.Wrap(err, "ChallengeResponse failed to parse.")
	}
	c.log.Printf("Got ChallengeID: %s\n", resp.ChallengeID)
	return resp.ChallengeID, resp.RequestID, nil
}

func (c Client) verifyRequest(rw *bufio.ReadWriter, pow *service.VerifyRequest) (id string, err error) {
	c.log.Printf("Ask server for verification: RequestID: %s\n", pow.RequestID)
	enc := gob.NewEncoder(rw)

	n, err := rw.WriteString("verify\n")
	if err != nil {
		return pow.RequestID, errors.Wrap(err, "Could not write ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(&pow)
	if err != nil {
		return pow.RequestID, errors.Wrapf(err, "Encode failed for: %#v", &pow)
	}
	err = rw.Flush()
	if err != nil {
		return pow.RequestID, errors.Wrap(err, "Flush failed.")
	}

	return pow.RequestID, nil
}

func (c Client) verifyResponse(rw *bufio.ReadWriter) (message string, requestID string, err error) {
	resp := &service.VerifyResponse{}
	decoder := gob.NewDecoder(rw)
	if err = decoder.Decode(resp); err != nil {
		return "", "", errors.Wrap(err, "VerifyResponse failed to parse.")
	}
	return resp.Message, resp.RequestID, nil
}
