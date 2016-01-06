// Copyright 2015 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textprot

import "bufio"
import "fmt"
import "strings"

import "../common"

const VERBOSE = false

type TextProt struct{}

func (t TextProt) Set(rw *bufio.ReadWriter, key []byte, value []byte) error {
	strKey := string(key)
	if VERBOSE {
		fmt.Printf("Setting key %s to value of length %v\n", strKey, len(value))
	}

	if _, err := fmt.Fprintf(rw, "set %s 0 0 %v\r\n%s\r\n", strKey, len(value), string(value)); err != nil {
		return err
	}

	rw.Flush()

	response, err := rw.ReadString('\n')
	if err != nil {
		return err
	}

	if VERBOSE {
		fmt.Println(response)
		fmt.Printf("Set key %s\n", strKey)
	}

	return nil
}

func (t TextProt) Get(rw *bufio.ReadWriter, key []byte) error {
	strKey := string(key)
	if VERBOSE {
		fmt.Printf("Getting key %s\n", strKey)
	}

	if _, err := fmt.Fprintf(rw, "get %s\r\n", strKey); err != nil {
		return err
	}

	rw.Flush()

	// read the header line
	response, err := rw.ReadString('\n')
	if err != nil {
		return err
	}
	if VERBOSE {
		fmt.Println(response)
	}

	if strings.TrimSpace(response) == "END" {
		if VERBOSE {
			fmt.Println("Empty response / cache miss")
		}
		return nil
	}

	// then read the value
	response, err = rw.ReadString('\n')
	if err != nil {
		return err
	}
	if VERBOSE {
		fmt.Println(response)
	}

	// then read the END
	response, err = rw.ReadString('\n')
	if err != nil {
		return err
	}
	if VERBOSE {
		fmt.Println(response)
		fmt.Printf("Got key %s\n", key)
	}
	return nil
}

func (t TextProt) BatchGet(rw *bufio.ReadWriter, keys [][]byte) error {
	if VERBOSE {
		fmt.Printf("Getting keys %v\n", keys)
	}

	cmd := []byte("get")
	space := byte(' ')
	end := []byte("\r\n")

	for _, key := range keys {
		cmd = append(cmd, space)
		cmd = append(cmd, key...)
	}

	cmd = append(cmd, end...)

	if _, err := fmt.Fprint(rw, string(cmd)); err != nil {
		return err
	}

	rw.Flush()

	for {
		// read the header line
		response, err := rw.ReadString('\n')
		if err != nil {
			return err
		}
		if VERBOSE {
			fmt.Println(response)
		}

		if strings.TrimSpace(response) == "END" {
			if VERBOSE {
				fmt.Println("End of batch response")
			}
			return nil
		}

		// then read the value
		response, err = rw.ReadString('\n')
		if err != nil {
			return err
		}
		if VERBOSE {
			fmt.Println(response)
		}
	}
}

func (t TextProt) GAT(rw *bufio.ReadWriter, key []byte) error {
	// Yes, the abstraction is a little bit leaky, but the code
	// in other places benefits from the consistency.
	panic("GAT in text protocol")
}

func (t TextProt) Delete(rw *bufio.ReadWriter, key []byte) error {
	strKey := string(key)
	if VERBOSE {
		fmt.Printf("Deleting key %s\n", strKey)
	}

	if _, err := fmt.Fprintf(rw, "delete %s\r\n", strKey); err != nil {
		return err
	}

	rw.Flush()

	response, err := rw.ReadString('\n')
	if err != nil {
		return err
	}
	if VERBOSE {
		fmt.Println(response)
		fmt.Printf("Deleted key %s\r\n", strKey)
	}
	return nil
}

func (t TextProt) Touch(rw *bufio.ReadWriter, key []byte) error {
	strKey := string(key)
	if VERBOSE {
		fmt.Printf("Touching key %s\n", strKey)
	}

	if _, err := fmt.Fprintf(rw, "touch %s %v\r\n", strKey, common.Exp()); err != nil {
		return err
	}

	rw.Flush()

	response, err := rw.ReadString('\n')
	if err != nil {
		return err
	}
	if VERBOSE {
		fmt.Println(response)
		fmt.Printf("Touched key %s\n", strKey)
	}
	return nil
}
