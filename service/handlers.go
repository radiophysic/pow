package service

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
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
		fmt.Printf("ChallengeHandler: decoder.Decode: %s", err)
		return
	}

	fmt.Println("Client asked for challenge, id:", request.RequestID)

	response := &ChallengeResponse{
		RequestID:   request.RequestID,
		ChallengeID: uuid.New().String(),
	}

	encoder := gob.NewEncoder(rw)
	if err := encoder.Encode(&response); err != nil {
		fmt.Printf("ChallengeHandler: encoder.Encode: %s", err)
		return
	}

	err := rw.Flush()
	if err != nil {
		fmt.Println("Flush write failure")
		return
	}
}

func (s Service) VerifyHandler(rw *bufio.ReadWriter) {
	var request VerifyRequest
	decoder := gob.NewDecoder(rw)
	if err := decoder.Decode(&request); err != nil {
		fmt.Printf("VerifyHandler: decoder.Decode: %s", err)
		return
	}
	fmt.Println("Client asked for verification, id:", request.RequestID)

	response := &VerifyResponse{
		RequestID: request.RequestID,
	}

	err := lib.Verify(request.Payload.Hash, request.Payload.Nonces)
	if err != nil {
		fmt.Printf("Verification failed: %s", err)
		response.Error.Code = 0
		response.Error.ErrorMessage = err.Error()
	}

	quote, err := getRandomQuote(s.cfg.DatasetFile)
	if err != nil {
		fmt.Printf("Service internal error: %s\n", err)
		response.Error.Code = 0
		response.Error.ErrorMessage = err.Error()
	}
	response.Message = quote

	encoder := gob.NewEncoder(rw)
	if err = encoder.Encode(&response); err != nil {
		fmt.Printf("VerifyHandler: encoder.Encode: %s", err)
		return
	}

	err = rw.Flush()
	if err != nil {
		fmt.Println("Flush write failure")
		return
	}
}

func getRandomQuote(f string) (string, error) {
	var q string
	var line int
	var file *os.File

	if _, err := os.Stat(f); err == nil {
		file, err = os.Open(f)
		if err != nil {
			fmt.Printf("getRandomQuote: flie open: %s", err.Error())
			return "", err
		}

		reader := bufio.NewReader(file)
		linesCount, err := utils.LineCounter(reader)
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

		q, _, err = utils.ReadLine(reader, line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Got EOF " + q)
			}
			return "", err
		}

		fmt.Printf("\nQuote picked: \n%s\n", q)
		defer func() {
			if err = file.Close(); err != nil {
				fmt.Println("Closing file error" + err.Error())
			}
		}()
	} else if os.IsNotExist(err) {
		fmt.Println("getRandomQuote: flie not found")
		return "", err
    } else {
		fmt.Printf("getRandomQuote: flie open: %s", err.Error())
        return "", err
    }

	return q, nil
}
