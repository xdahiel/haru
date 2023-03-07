package fileReader

import (
	"bufio"
	"log"
	"os"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(path string) *Scanner {
	f, err := os.Open(path)
	log.Println("open file: ", path)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	return &Scanner{
		r: r,
	}
}

func (s *Scanner) Next() (string, bool) {
	line, isPrefix, err := s.r.ReadLine()
	if err != nil {
		return "", false
	}

	res := string(line)
	for isPrefix {
		line, isPrefix, _ = s.r.ReadLine()
		res += string(line)
	}

	return res, true
}
