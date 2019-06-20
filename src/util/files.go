package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string, orquestrador *sync.WaitGroup) (err error) {
	in, err := os.Open(src)
	if err != nil {
		orquestrador.Done()
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		orquestrador.Done()
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		orquestrador.Done()
		return
	}

	err = out.Sync()
	if err != nil {
		orquestrador.Done()
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		orquestrador.Done()
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		orquestrador.Done()
		return
	}
	orquestrador.Done()
	return
}

//CopyDirAsync do this sh** async Bro!!!
func CopyDirAsync(src string, dst string, orquestrador *sync.WaitGroup) (err error) {
	orquestrador.Add(1)
	go copyDir(src, dst, orquestrador)
	orquestrador.Wait()
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string, orquestrador *sync.WaitGroup) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		orquestrador.Done()
		return err
	}
	if !si.IsDir() {
		orquestrador.Done()
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		orquestrador.Done()
		return
	}
	if err == nil {
		orquestrador.Done()
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		orquestrador.Done()
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		orquestrador.Done()
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			orquestrador.Add(1)
			go copyDir(srcPath, dstPath, orquestrador)
			// if err != nil {
			// 	orquestrador.Done()
			// 	return
			// }
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			orquestrador.Add(1)
			go CopyFile(srcPath, dstPath, orquestrador)
			// if err != nil {
			// 	orquestrador.Done()
			// 	return
			// }
		}
	}
	orquestrador.Done()
	return
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}
