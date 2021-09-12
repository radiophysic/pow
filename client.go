package main

import (
    "bufio"
    "encoding/gob"
    "fmt"
    "log"
    "strconv"

    "github.com/google/uuid"
    "github.com/pkg/errors"
    "tcp/pkg/lib"
    "tcp/service"
)

func client(ip string) error {
    rw, err := lib.Open(ip + lib.Port)
    if err != nil {
        fmt.Println("The client cannot link to change the address:" + ip + lib.Port)
        return err
    }

    var challengeRequestID, respRequestID, ChallengeID string
    if challengeRequestID, err = challengeRequest(rw); err != nil {
        return err
    }

    if ChallengeID, respRequestID, err = challengeResponse(rw); err != nil {
        return err
    }

    if challengeRequestID != respRequestID {
        return errors.New("challengeRequestID != respRequestID")
    }

    // TBD: check existed blocks

    var hash []byte
    var nonces []uint32
    if hash, nonces, err = lib.Work(); err != nil {
        return errors.Wrap(err, "Proof-of-work failed.")
    }

    pow := &service.VerifyRequest{
        RequestID:   uuid.New().String(),
        ChallengeID: ChallengeID,
        Payload: struct {
            Hash   []byte
            Nonces []uint32
        }{Hash: hash, Nonces: nonces},
    }

    if _, err := verifyRequest(rw, pow); err != nil {
        return err
    }

    message, lastRequestId, err := verifyResponse(rw)
    if err != nil {
        return err
    }

    if pow.RequestID != lastRequestId {
        return errors.New("pow.RequestID != lastRequestId")
    }

    fmt.Printf("\n\nWork has been proved\n Quote: %s\n\n", message)

    return nil
}

func challengeRequest(rw *bufio.ReadWriter) (id string, err error) {
    challengeReq := service.ChallengeRequest{RequestID: uuid.New().String()}
    log.Printf("Ask server for challenge: RequestID: %s\n", challengeReq.RequestID)
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

func challengeResponse(rw *bufio.ReadWriter) (ChallengeID string, RequestID string, err error) {
    resp := &service.ChallengeResponse{}
    decoder := gob.NewDecoder(rw)
    if err := decoder.Decode(resp); err != nil {
        return "", "", errors.Wrap(err, "ChallengeResponse failed to parse.")
    }
    log.Printf("Got ChallengeID: %s\n", resp.ChallengeID)
    return resp.ChallengeID, resp.RequestID, nil
}

func verifyRequest(rw *bufio.ReadWriter, pow *service.VerifyRequest) (id string, err error) {
    log.Printf("Ask server for verification: RequestID: %s\n", pow.RequestID)
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

func verifyResponse(rw *bufio.ReadWriter) (Message string, RequestID string, err error) {
    resp := &service.VerifyResponse{}
    decoder := gob.NewDecoder(rw)
    if err := decoder.Decode(resp); err != nil {
        return "", "", errors.Wrap(err, "VerifyResponse failed to parse.")
    }
    return resp.Message, resp.RequestID, nil
}

func main() {
    err := client("localhost")
    if err != nil {
        fmt.Println("Error:", errors.WithStack(err))
    }
}
