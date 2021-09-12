package lib

import (
    "crypto/rand"

    "github.com/AidosKuneen/cuckoo"
)

const maxInt = 1<<32 - 1

func Work() ([]byte, []uint32, error) {
    c := cuckoo.NewCuckoo()
    hash := make([]byte, 16)
    nounces := make([]uint32, 0)
    var ok bool
    for i := 0; i < maxInt; i++ {
        if _, err := rand.Read(hash); err != nil {
            return nil, nil, err
        }
        nounces, ok = c.PoW(hash)
        if ok {
            break
        }
    }

    return hash, nounces, nil
}

func Verify(hash []byte, nonces []uint32) error {
    return cuckoo.Verify(hash, nonces)
}
