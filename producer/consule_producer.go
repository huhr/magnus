// 从io.Reader读取数据文件到channel中
package producer

import (
	"io"
)

const (
	// line size must be smaller than bufferSize or front of the line will be lost
	bufferSize = 10^6
)

// 从io.Reader中获取数据
type ConsoleProducer struct {
	reader   io.Reader
	buf      []byte     // buffer
	r, w     int        // buffer的读写位置
	eof      bool
}

func NewConsoleProducer(reader io.Reader) *ConsoleProducer {
	return &ConsoleProducer{
		reader: reader,
		buf: make([]byte, bufferSize),
	}
}

// 讲一行数据读入channel
func (c *ConsoleProducer) Produce(channel chan []byte) error {
	var offset int
	for true {
		for c.r < c.w {
			// search '\n' from buffer[r:w]
			for i, b := range c.buf[c.r:c.w] {
				if b == '\n' {
					offset = i + 1
					break
				// get no line, need filling
				} else if i == c.w - c.r - 1 {
					c.r = c.w
					offset = 0
					break
				}
			}
			if offset != 0 {
				// get line
				msg := make([]byte, offset)
				copy(msg, c.buf[c.r:c.r+offset])
				channel <- msg
				c.r += offset
				// 重置offset
				offset = 0
			} else {
				break
			}
		}
		// read EOF, stop
		if c.eof {
			break
		}
		// 先读取一次
		n, err := c.fill()
		if err != nil {
			return err
		}
		for n == 0 {
			_, err := c.fill()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ConsoleProducer) fill() (n int, err error) {
	// 丢掉已经读取的数据
	if c.r > 0 {
		copy(c.buf, c.buf[c.r:c.w])
		c.w -= c.r
		c.r = 0
	}

	// io.Reader is not block
	n, err = c.reader.Read(c.buf[c.w:])
	c.w += n
	if err != nil {
		// read EOF, 尝试再读几次
		if err == io.EOF {
			c.eof = true
			return n, nil
		} else {
			return n, err
		}
	}
	if n > 0 {
		return n, nil
	}
	return
}

