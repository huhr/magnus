package tools

import (
	"bytes"
	"errors"
	"io"
)

const (
	DefaultBufferSize = 1 << 10
)

type UnitReader struct {
	reader    io.Reader
	delimiter string
	buf       []byte
	retry     int
	r, w      int // buffer的读写位置
}

func NewUnitReader(reader io.Reader, delimiter string, bufSize int) *UnitReader {
	if bufSize == 0 {
		bufSize = DefaultBufferSize
	}
	return &UnitReader{
		reader:    reader,
		delimiter: delimiter,
		buf:       make([]byte, bufSize),
	}
}

func (c *UnitReader) ResetReader(reader io.Reader) {
	c.reader = reader
}

// 根据分隔符read一个数据单元
func (c *UnitReader) ReadOne() (msg []byte, err error) {
	for true {
		offset := bytes.Index(c.buf[c.r:c.w], []byte(c.delimiter))
		if offset == -1 {
			// buf已经读满了，直接返回buf中的全部数据作为一个数据单元
			if c.w-c.r == cap(c.buf) {
				msg := make([]byte, len(c.buf))
				copy(msg, c.buf)
				c.r = c.w
				return msg, errors.New("too large msg")
			}
			// 再读一次，没读到数据的话不需要continue了，返回给上层
			if n, err := c.fill(); n == 0 && err != nil {
				return make([]byte, 0), err
			}

			// 读到了数据，而且buffer还没有满，再查询一次
			continue
		}
		msg := make([]byte, offset)
		copy(msg, c.buf[c.r:c.r+offset])
		c.r += offset + len(c.delimiter)
		return msg, nil
	}
	// 这里不会被触发
	return make([]byte, 0), nil
}

// 填充数据
func (c *UnitReader) fill() (n int, err error) {
	// 将未读书移到缓存首部
	if c.r > 0 {
		copy(c.buf, c.buf[c.r:c.w])
		c.w -= c.r
		c.r = 0
	}
	n, err = c.reader.Read(c.buf[c.w:])
	c.w += n
	return
}
