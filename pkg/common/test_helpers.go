package common

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

type ChunkFunc func(chunk *sources.Chunk) error

var MatchError = errors.New("chunk doesn't match")

func HandleTestChannel(chunksCh chan *sources.Chunk, cf ChunkFunc) error {
	for {
		select {
		case gotChunk := <-chunksCh:
			err := cf(gotChunk)
			if err != nil {
				if errors.Is(err, MatchError) {
					continue
				}
				return err
			}
			return nil
		case <-time.After(10 * time.Second):
			return fmt.Errorf("no new chunks recieved after 10 seconds")
		}
	}
}

func WriteTestFile(filename string, content []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return f.Close()
}
