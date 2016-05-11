package producer

import (
	"bytes"
	"io"
)

const (
	// line size must be smaller than bufferSize or front of the line will be lost
	DefaultBufferSize = 1<<10
)

type UnitReader struct {
	reader		io.Reader
	delimiter	string
	buf			[]byte
	r, w		int        // buffer的读写位置
	eof			bool       // 是否读到文件尾
}

func NewUnitReader(reader io.Reader, delimiter string, bufSize int) *UnitReader {
	if bufSize == 0 {
		bufSize = DefaultBufferSize
	}
	return &UnitReader{
		reader: reader,
		delimiter: delimiter,
		buf: make([]byte, bufSize),
	}
}

// 根据分隔符read一个数据单元，有可能读到文件尾了，这时候
// 返回异常
func (c *UnitReader) ReadOne() ([]byte, error) {
	for true {
		var offset int
		// 不用每次都读数据
		if c.r == c.w {
			if _, err := c.fill(); err != nil {
				return c.buf[c.r:c.w], err
			}
		}
		offset = bytes.Index(c.buf[c.r:c.w], []byte(c.delimiter))
		// 找不到的情况下
		if offset == -1 {
			// 重新读一次数据
			if c.w - c.r != cap(c.buf) {
				if _, err := c.fill(); err != nil {
					return c.buf[c.r:c.w], err
				}
				continue
			}
			// buf已经读满了，直接返回buf中的全部数据作为一个数据单元
			msg := make([]byte, len(c.buf))
			copy(msg, c.buf)
			c.r = c.w
			return msg, nil
		}
		// 读到数据了
		msg := make([]byte, offset)
		copy(msg, c.buf[c.r:c.r + offset])
		c.r += offset + len(c.delimiter)
		return msg, nil
	}
	return nil, nil
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

