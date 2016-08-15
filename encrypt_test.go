package aclient_test

import (
        "bytes"
        "testing"
        "github.com/tortuoise/aclient"
)

func TestRSA(t *testing.T) {

        tests := []string{ "send reinforcement, we're going to advance", "The magic words are squeamish ossifrage"}
        newrsa,err := aclient.NewRSA()
        if err != nil {
                t.Errorf("Error setting up: %q", err)
        }

        for _, test := range tests {
                newrsa.Plain = []byte(test)
                newrsa.Cipher,err = newrsa.Encrypt()
                if err != nil {
                        t.Errorf("Error encrypting: %q", err)
                }
                newrsa.Plain,err = newrsa.Decrypt()
                if err != nil {
                        t.Errorf("Error decrypting: %q", err)
                }
                if bytes.Compare(newrsa.Plain,[]byte(test)) != 0 {
                        t.Error("Expecting: ", []byte(test),"Got: ", newrsa.Plain)
                }
        }
}

func TestRand(t *testing.T) {

        tests := []string{ "send reinforcement, we're going to advance", "The magic words are squeamish ossifrage"}
        newrsa,err := aclient.NewRSA()
        if err != nil {
                t.Errorf("Error setting up: %q", err)
        }

        for _, test := range tests {
                newrsa.Plain = []byte(test)
                newrsa.Cipher,err = newrsa.Encrypt()
                if err != nil {
                        t.Errorf("Error encrypting: %q", err)
                }
                oldCipher := make([]byte, len(newrsa.Cipher))
                copy(oldCipher, newrsa.Cipher)
                newrsa.Randomize()
                newrsa.Cipher,err = newrsa.Encrypt()
                if err != nil {
                        t.Errorf("Error encrypting: %q", err)
                }
                if bytes.Compare(newrsa.Cipher, oldCipher) == 0 {
                        t.Error("Expecting: ", oldCipher,"Got: ", newrsa.Cipher)
                }
        }

}
