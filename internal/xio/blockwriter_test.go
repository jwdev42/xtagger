package xio

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"testing"
)

func TestBlockWriter(t *testing.T) {
	pseudoRandomData := func(size int) []byte {
		data := make([]byte, size)
		for i := 0; i < len(data); i++ {
			data[i] = byte(rand.Intn(255))
		}
		return data
	}
	withOffset := func(src []byte, blockSize int) []byte {
		remainder := len(src) % blockSize
		if remainder == 0 {
			return src
		}
		dst := make([]byte, len(src)+blockSize-remainder)
		copy(dst, src)
		return dst
	}

	singleWrite := func(data []byte, wr io.Writer) error {
		_, err := wr.Write(data)
		return err
	}

	multiWrite := func(data []byte, wr io.Writer) error {
		const blocksizeMulti = 64
		//write 64byte blocks
		for i := 0; i < len(data)/blocksizeMulti; i++ {
			offset := i * blocksizeMulti
			_, err := wr.Write(data[offset : offset+blocksizeMulti])
			if err != nil {
				return err
			}
		}
		//write remaining data
		if len(data)%blocksizeMulti > 0 {
			offset := len(data) / blocksizeMulti * blocksizeMulti
			_, err := wr.Write(data[offset:])
			if err != nil {
				return err
			}
		}
		return nil
	}

	testFunc := func(data []byte, blocksize int, writeFunc func([]byte, io.Writer) error) error {
		buf := new(bytes.Buffer)
		wr := NewBlockWriter(buf, blocksize)
		//Write
		if err := writeFunc(data, wr); err != nil {
			return err
		}
		_, err := wr.Finalize()
		if err != nil {
			return err
		}
		//Compare
		res := buf.Bytes()
		exp := withOffset(data, blocksize)
		if !bytes.Equal(res, exp) {
			return fmt.Errorf("Buffers are not equal (len exp: %d, len res: %d)", len(exp), len(res))
		}
		return nil
	}

	const blocksize = 512
	const maxLen = 4096
	//Execute tests
	for i := 0; i <= maxLen; i++ {
		test := pseudoRandomData(i)
		if err := testFunc(test, blocksize, singleWrite); err != nil {
			t.Errorf("SingleWrite: Error on test length %d: %s", i, err)
		}
		if err := testFunc(test, blocksize, multiWrite); err != nil {
			t.Errorf("MultiWrite: Error on test length %d: %s", i, err)
		}
	}
}
