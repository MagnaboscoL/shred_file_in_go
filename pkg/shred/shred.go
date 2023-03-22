package shred

import (
	"crypto/rand"
	"fmt"
	"os"
)

const MaxChunkSize int64 = 10 * 1024 * 1024

func writeChunk(file *os.File, buffer []byte, size int64, offset *int64) error {
	// Fill the slice with random data
	n, err := rand.Read(buffer)
	if err != nil {
		return err
	}
	if int64(n) != size {
		return fmt.Errorf("failed to create chunk data")
	}

	// Write the random data to the file
	n, err = file.WriteAt(buffer, *offset)
	if err != nil {
		return err
	}
	if int64(n) != size {
		return fmt.Errorf("failed to write chunk")
	}
	*offset += size
	return nil
}

func overwriteInChunks(file *os.File, size int64, buffer []byte) error {
	bufferSize := int64(len(buffer))
	chunks := size / bufferSize
	offset := int64(0)
	for i := int64(0); i < chunks; i++ {
		err := writeChunk(file, buffer, bufferSize, &offset)
		if err != nil {
			return err
		}
	}
	if offset < size {
		remainder := size - offset
		err := writeChunk(file, buffer[:remainder], remainder, &offset)
		if err != nil {
			return err
		}
	}

	return nil
}

func overwriteNTimes(fileName string, size int64, n int) error {
	// Open the file for writing with the truncate flag
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	// Generate a byte slice whose size is limited to avoid too much RAM usage.
	// TODO Anyway the impact of the RAM usage should be verified on different systems to check if this approach is useful
	// or the straightforward approach with one big slice would be enough and favoured due to simplicity.
	chunkSize := size
	if size > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	// Allocate the buffer once and use it multiple times.
	data := make([]byte, chunkSize)

	// Write random data to the file n times
	for i := 0; i < n; i++ {
		err = overwriteInChunks(file, size, data)
		if err != nil {
			return err
		}
		// Sync the file to disk to ensure data is written
		if err := file.Sync(); err != nil {
			return err
		}
	}
	return nil
}

func Shred(fileName string) error {
	stats, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if stats.IsDir() {
		return fmt.Errorf("Shred works only with files and '%s' is a directory", fileName)
	}

	if stats.Size() > 0 {
		// Write the file multiple times
		if err := overwriteNTimes(fileName, stats.Size(), 3); err != nil {
			return err
		}
	}

	// TODO Rename the file multiple times using atomic operations

	// Delete the file
	if err := os.Remove(fileName); err != nil {
		return err
	}

	return nil
}
