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

func ChallengeHandler(rw *bufio.ReadWriter) {
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

func VerifyHandler(rw *bufio.ReadWriter) {
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

	quote, err := getRandomQuote()
	if err != nil {
		fmt.Printf("Service internal error: %s", err)
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

func getRandomQuote() (string, error) {
	var q string
	var line int
	var file *os.File
	fname := "assets/dataset.txt"
	if _, err := os.Stat(fname); err == nil {
		file, err = os.Open(fname)
		if err != nil {
			return "", err
		}
		defer func() {
			if err = file.Close(); err != nil {
				fmt.Println("Closing file error" + err.Error())
			}
		}()

		rand.Seed(time.Now().UnixNano())
		line = rand.Intn(400)
		if line == 0 {
			line = 1
		}
		reader := bufio.NewReader(file)
		q, _, err = ReadLine(reader, line)
		if err != nil {
			return "", err
		}

		fmt.Printf("\nQuote picked: \n%s\n", q)

	} else if os.IsNotExist(err) {
		return "", err
    } else {
        return "", err
    }

	return q, nil
}

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}
