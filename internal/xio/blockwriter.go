package xio

import (
	"bytes"
	"io"
)

// BlockWriter writes data in fixed-size blocks to block-oriented devices
// such as tape drives.
type BlockWriter struct {
	w      io.Writer     //Wrapped writer
	buf    *bytes.Buffer //Internal buffer of length blockSize
	bs     int           //Block size
	broken bool          //True after a write error to the wrapped writer
}

// NewBlockWriter creates a new BlockWriter with writer as the underlying
// write target. Parameter blockSize specifies the size of each data block
// to be written, must be > 0.
func NewBlockWriter(writer io.Writer, blockSize int) *BlockWriter {
	return &BlockWriter{
		w:   writer,
		buf: new(bytes.Buffer),
		bs:  blockSize,
	}
}

// Replaces the internal buffer with a new one that holds only unread data.
func (r *BlockWriter) compactBuffer() {
	r.buf = bytes.NewBuffer(r.buf.Bytes())
}

// Writes blocks of blocksize to the wrapped Writer as long as buffer
// provides enough data to fill a block. Return parameter n returns the
// amount of bytes written to the wrapped Writer.
func (r *BlockWriter) writeBlocks() (n int, err error) {
	defer r.compactBuffer()
	//Calculate blocks to write
	blocks := r.buf.Len() / r.bs
	//Write blocks of r.bs bytes
	for i := 0; i < blocks; i++ {
		block := r.buf.Next(r.bs)
		written, err := r.w.Write(block)
		if err != nil {
			newBuf := bytes.NewBuffer(block[written:])
			newBuf.Write(r.buf.Bytes())
			r.buf = newBuf
			return i*r.bs + written, err
		}
	}
	return blocks * r.bs, nil
}

// Appends p to BlockWriter's internal buffer buf, then writes
// len(buf)/blocksize blocks of data to the wrapped Writer. The remaining
// data will be kept in the internal buffer until a new Write call appends
// enough data to write a block or Finalize() is called.
//
// Parameter n returns the amount of bytes written to the internal buffer.
// If an error occurs while writing a block, Write() will return a non-nil
// err. Write will panic if called again after an error occured.
func (r *BlockWriter) Write(p []byte) (n int, err error) {
	if r.broken {
		panic("BUG: Writing to a broken BlockWriter is forbidden")
	}
	//Append p to buffer
	r.buf.Write(p)
	//Write blocks
	_, err = r.writeBlocks()
	return len(p), err
}

func (r *BlockWriter) Finalize() (n int, err error) {
	if r.buf.Len() == 0 {
		return 0, nil
	}
	//Calculate padding bytes to append
	remainder := r.buf.Len() % r.bs
	if remainder > 0 {
		//Write \0-Terminators as padding
		r.buf.Write(make([]byte, r.bs-remainder))
	}
	//Write blocks
	return r.writeBlocks()
}
