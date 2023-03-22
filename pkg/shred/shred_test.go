package shred

import (
	"crypto/rand"
	"encoding/binary"
	"os"
	"testing"
)

func Test_EmptyFile(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/emptyfile"
	func() {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
	}()

	// Shred the file
	if err := Shred(fileName); err != nil {
		t.Fatal(err)
	}

	// Check that the file has been deleted
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, but it still exists")
	}
}

func Test_ShortText(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/short_text_file.txt"
	func() {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Write some data to the file
		_, err = file.WriteString("Short text")
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Shred the file
	if err := Shred(fileName); err != nil {
		t.Fatal(err)
	}

	// Check that the file has been deleted
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, but it still exists")
	}
}

func Test_BinaryFile(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/binaryfile.hex"
	func() {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Write some binary data to the file
		err = binary.Write(file, binary.LittleEndian, []byte{0x00, 0x01, 0x02, 0x03})
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Shred the file
	if err := Shred(fileName); err != nil {
		t.Fatal(err)
	}

	// Check that the file has been deleted
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, but it still exists")
	}
}

func Test_BigBinaryFile(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/binaryfile.hex"
	func() {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Write 101 MB to the file.
		buf := make([]byte, 1024*1024)
		for i := 0; i < 101; i++ {
			_, err := rand.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			err = binary.Write(file, binary.LittleEndian, buf)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// Shred the file
	if err := Shred(fileName); err != nil {
		t.Fatal(err)
	}

	// Check that the file has been deleted
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, but it still exists")
	}
}

func Test_OpenFileFails(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/file.hex"
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatal(err)
	}
	// Note that this defer will be called only when this function returns, therefore the file remains open when calling Shred().
	defer file.Close()

	// Write some binary data to the file
	err = binary.Write(file, binary.LittleEndian, []byte{0x00, 0x01, 0x02, 0x03})
	if err != nil {
		t.Fatal(err)
	}

	// Shred the file
	if err := Shred(fileName); err != nil {
		return
	}

	t.Fatalf("Shred should have failed")
}

func Test_NonExistentFileFails(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/tmp.txt"

	// Shred the file
	if err := Shred(fileName); err != nil {
		return
	}

	t.Fatalf("Shred should have failed")
}

func Test_ReadOnlyFileFails(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	fileName := tmpDir + "/read_only.txt"
	func() {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Write some data to the file
		_, err = file.WriteString("Short text")
		if err != nil {
			t.Fatal(err)
		}
	}()

	err := os.Chmod(fileName, 0444)
	if err != nil {
		t.Fatal(err)
	}

	// Shred the file
	if err := Shred(fileName); err != nil {
		return
	}

	t.Fatalf("Shred should have failed")
}

func Test_DirectoryFails(t *testing.T) {
	// Create a test directory
	tmpDir := t.TempDir()
	fileName := tmpDir + "/testdir"
	err := os.Mkdir(fileName, 0777)
	if err != nil {
		t.Fatal(err)
	}

	// Shred the file
	if err := Shred(fileName); err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
