package service_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"tcp/conf"
	"tcp/service"
)

func MockService() service.Service{
	cfg :=  conf.Config{
		ServerAddr:  "localhost",
		ServerPort:  ":7777",
		DatasetFile: "../assets/dataset.txt",
	}
	return service.NewService(cfg, log.Default())
}

func TestChallengeHandler(t *testing.T) {
	s := MockService()

	req := service.ChallengeRequest{RequestID: "6056f1c3-b316-4bf4-bc42-244410340ded"}
	resp := service.ChallengeResponse{}

	var clientReq, serverResp bytes.Buffer

	// emulate client's request
	encoder := gob.NewEncoder(&clientReq)
	err := encoder.Encode(&req)

	testID :=0
	t.Logf("\tTest %d:\tencoder.Encode no error.", testID)
	{
		require.NoError(t, err)
	}

	// emulate network
	r := bufio.NewReader(bytes.NewReader(clientReq.Bytes()))
	w := bufio.NewWriter(&serverResp)
	rw := bufio.NewReadWriter(r, w)

	// call actual endpoint
	s.ChallengeHandler(rw)

	decoder := gob.NewDecoder(&serverResp)
	err = decoder.Decode(&resp)

	testID++
	t.Logf("\tTest %d:\tdecoder.Decode no error.", testID)
	{
		require.NoError(t, err)
	}

	testID++
	t.Logf("\tTest %d:\treq.RequestID equal to resp.RequestID.", testID)
	{
		require.Equal(t, req.RequestID, resp.RequestID)
	}

	testID++
	t.Logf("\tTest %d:\tresp.ChallengeID is NOT empty.", testID)
	{
		require.NotEmpty(t, resp.ChallengeID)
	}

	testID++
	t.Logf("\tTest %d:\tresp.ChallengeID is a proper UUID string.", testID)
	{
		_, err = uuid.Parse(resp.ChallengeID)
		require.NoError(t, err)
	}
}

func TestVerifyHandler_BrokenPayload(t *testing.T) {
	s := MockService()
	
	req := service.VerifyRequest{
		RequestID:   "6ca1efb1-7134-406e-9304-0c3e774fd8ba",
		ChallengeID: "c9c0747d-587a-4f65-97fa-332ebddc6ed3",
		Payload:     service.ProofOfWork{
			Hash:   nil,
			Nonces: nil,
		},
	}
	
	resp := service.VerifyResponse{}

	var clientReq, serverResp bytes.Buffer

	// emulate client's request
	encoder := gob.NewEncoder(&clientReq)
	err := encoder.Encode(&req)

	testID :=0
	t.Logf("\tTest %d:\tencoder.Encode no error.", testID)
	{
		require.NoError(t, err)
	}

	// emulate network
	r := bufio.NewReader(bytes.NewReader(clientReq.Bytes()))
	w := bufio.NewWriter(&serverResp)
	rw := bufio.NewReadWriter(r, w)

	// call actual endpoint
	s.VerifyHandler(rw)

	decoder := gob.NewDecoder(&serverResp)
	err = decoder.Decode(&resp)
	testID++
	t.Logf("\tTest %d:\tdecoder.Decode no error.", testID)
	{
		require.NoError(t, err)
	}

	testID++
	t.Logf("\tTest %d:\tverify payload is malformed.", testID)
	{
		require.Equal(t, 2, resp.Error.Code)
	}
}

func TestVerifyHandler_PayloadWrong(t *testing.T) {
	s := MockService()

	var k0 uint64 = 0xf4956dc403730b01
	var k1 uint64 = 0xe6d45de39c2a5a3e
	nonces := []uint32{
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
		0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
	}

	hash := make([]byte, 16)
	binary.LittleEndian.PutUint64(hash, k0)
	binary.LittleEndian.PutUint64(hash[8:], k1)

	req := service.VerifyRequest{
		RequestID:   "6ca1efb1-7134-406e-9304-0c3e774fd8ba",
		ChallengeID: "c9c0747d-587a-4f65-97fa-332ebddc6ed3",
		Payload:     service.ProofOfWork{
			Hash:   hash,
			Nonces: nonces,
		},
	}

	resp := service.VerifyResponse{}

	var clientReq, serverResp bytes.Buffer

	// emulate client's request
	encoder := gob.NewEncoder(&clientReq)
	err := encoder.Encode(&req)

	testID :=0
	t.Logf("\tTest %d:\tencoder.Encode no error.", testID)
	{
		require.NoError(t, err)
	}

	// emulate network
	r := bufio.NewReader(bytes.NewReader(clientReq.Bytes()))
	w := bufio.NewWriter(&serverResp)
	rw := bufio.NewReadWriter(r, w)

	// call actual endpoint
	s.VerifyHandler(rw)

	decoder := gob.NewDecoder(&serverResp)
	err = decoder.Decode(&resp)
	testID++
	t.Logf("\tTest %d:\tdecoder.Decode no error.", testID)
	{
		require.NoError(t, err)
	}

	testID++
	t.Logf("\tTest %d:\tVerification failed.", testID)
	{
		require.Equal(t, 3, resp.Error.Code)
	}
}

func TestVerifyHandler_Successful(t *testing.T) {
	s := MockService()

	var k0 uint64 = 0xf4956dc403730b01
	var k1 uint64 = 0xe6d45de39c2a5a3e
	nonces := []uint32{
		0x6d31e, 0x72b0e, 0x7aaaf, 0x134522, 0x18cdb9,
		0x1ffaef, 0x28b919, 0x43d8fa, 0x7fc4fb, 0x968240,
		0xa28796, 0xad8119, 0xb6b419, 0xbbddd6, 0xbd2765,
		0xcb572a, 0xe090d9, 0xeea5a5, 0xf2898f, 0xfa27c0,
	}

	hash := make([]byte, 16)
	binary.LittleEndian.PutUint64(hash, k0)
	binary.LittleEndian.PutUint64(hash[8:], k1)

	req := service.VerifyRequest{
		RequestID:   "6ca1efb1-7134-406e-9304-0c3e774fd8ba",
		ChallengeID: "c9c0747d-587a-4f65-97fa-332ebddc6ed3",
		Payload:     service.ProofOfWork{
			Hash:   hash,
			Nonces: nonces,
		},
	}

	resp := service.VerifyResponse{}

	var clientReq, serverResp bytes.Buffer

	// emulate client's request
	encoder := gob.NewEncoder(&clientReq)
	err := encoder.Encode(&req)

	testID :=0
	t.Logf("\tTest %d:\tencoder.Encode no error.", testID)
	{
		require.NoError(t, err)
	}

	// emulate network
	r := bufio.NewReader(bytes.NewReader(clientReq.Bytes()))
	w := bufio.NewWriter(&serverResp)
	rw := bufio.NewReadWriter(r, w)

	// call actual endpoint
	s.VerifyHandler(rw)

	decoder := gob.NewDecoder(&serverResp)
	err = decoder.Decode(&resp)
	testID++
	t.Logf("\tTest %d:\tdecoder.Decode no error.", testID)
	{
		require.NoError(t, err)
	}

	testID++
	t.Logf("\tTest %d:\tVerification successful.", testID)
	{
		require.Equal(t, 0, resp.Error.Code)
		require.NotEmpty(t, resp.Message)
		require.Equal(t, req.RequestID, resp.RequestID)
	}
}
