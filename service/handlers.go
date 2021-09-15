package service

import (
	"bufio"
	"encoding/gob"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"tcp/pkg/lib"
	"tcp/pkg/utils"
)

type ChallengeRequest struct {
	RequestID string
}

type ChallengeResponse struct {
	RequestID   string
	ChallengeID string
}

type VerifyRequest struct {
	RequestID   string
	ChallengeID string
	Payload     ProofOfWork
}

type VerifyResponse struct {
	RequestID string
	Message   string
	Error     VerificationErr
}

type ProofOfWork struct {
	Hash   []byte
	Nonces []uint32
}

type VerificationErr struct {
	Code         int
	ErrorMessage string
}

func (s Service) ChallengeHandler(rw *bufio.ReadWriter) {
	var request ChallengeRequest

	decoder := gob.NewDecoder(rw)
	if err := decoder.Decode(&request); err != nil {
		s.log.Printf("ChallengeHandler: decoder.Decode: %s", err)
		return
	}

	s.log.Println("Client asked for challenge, id:", request.RequestID)

	response := &ChallengeResponse{
		RequestID:   request.RequestID,
		ChallengeID: uuid.New().String(),
	}

	encoder := gob.NewEncoder(rw)
	if err := encoder.Encode(&response); err != nil {
		s.log.Printf("ChallengeHandler: encoder.Encode: %s", err)
		return
	}

	err := rw.Flush()
	if err != nil {
		s.log.Println("Flush write failure")
		return
	}
}

func (s Service) VerifyHandler(rw *bufio.ReadWriter) {
	var request VerifyRequest
	var response VerifyResponse

	decoder := gob.NewDecoder(rw)
	if err := decoder.Decode(&request); err != nil {
		response.Error.Code = 1
		response.Error.ErrorMessage = err.Error()
		s.writeVerifyResp(rw, &VerifyResponse{})
		return
	}

	response.RequestID = request.RequestID
	s.log.Printf("Client asked for verification, RequestID:%s\nPayload: %v",
		request.RequestID, request.Payload)

	if request.Payload.Hash == nil || request.Payload.Nonces == nil {
		response.Error.Code = 2
		response.Error.ErrorMessage = "payload is malformed"
		s.writeVerifyResp(rw, &response)
		return
	}

	err := lib.Verify(request.Payload.Hash, request.Payload.Nonces)
	if err != nil {
		s.log.Printf("Verification failed: %s", err)
		response.Error.Code = 3
		response.Error.ErrorMessage = err.Error()
		s.writeVerifyResp(rw, &response)
		return
	}

	quote, err := s.getRandomQuote()
	if err != nil {
		s.log.Printf("Service internal error: %s\n", err)
		response.Error.Code = 4
		response.Error.ErrorMessage = err.Error()
		s.writeVerifyResp(rw, &response)
		return
	}
	response.Message = quote
	s.writeVerifyResp(rw, &response)
}

func (s Service) writeVerifyResp(rw *bufio.ReadWriter, response *VerifyResponse) {
	encoder := gob.NewEncoder(rw)
	if err := encoder.Encode(&response); err != nil {
		s.log.Printf("VerifyHandler: encoder.Encode: %s", err)
		return
	}

	err := rw.Flush()
	if err != nil {
		s.log.Println("Flush write failure")
		return
	}
}

func (s Service) getRandomQuote() (string, error) {
	var (
		err error
		q string
		line, linesCount int
		file *os.File
	)

	if _, err = os.Stat(s.cfg.DatasetFile); err != nil {
		s.log.Println("getRandomQuote: file not found")
		return "", err
	}

	if file, err = os.Open(s.cfg.DatasetFile); err != nil {
		s.log.Printf("getRandomQuote: file open: %s", err.Error())
		return "", err
	}

	reader := bufio.NewReader(file)
	linesCount, err = utils.LineCounter(reader)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	line = rand.Intn(linesCount)
	if line == 0 {
		line = 1
	}

	// rewind position to the beginning of the file because EOF reached due to LineCounter call
	_,_ = file.Seek(0, 0)

	if q, _, err = utils.ReadLine(reader, line); err != nil {
		if errors.Is(err, io.EOF) {
			s.log.Println("Got EOF " + q)
		}
		return "", err
	}

	s.log.Printf("\nQuote picked: \n%s\n", q)
	defer func() {
		if err = file.Close(); err != nil {
			s.log.Println("Closing file error" + err.Error())
		}
	}()

	return q, nil
}
