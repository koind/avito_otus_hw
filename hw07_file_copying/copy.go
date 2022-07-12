package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fileFromStat, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if !fileFromStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if fileFromStat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	_, _ = fromFile.Seek(offset, io.SeekStart)

	if limit == 0 || fileFromStat.Size()-offset < limit {
		limit = fileFromStat.Size() - offset
	}

	bar := pb.New64(limit)
	bar.SetUnits(pb.U_BYTES)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	_, err = io.CopyN(toFile, barReader, limit)
	if err != nil {
		return err
	}

	return nil
}
