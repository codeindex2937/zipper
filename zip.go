package zipper

import (
	"io"

	"github.com/kdungs/zip"
)

type File = zip.File

type ZipReader struct {
	r *zip.Reader
}

func NewReader(instream io.ReaderAt, size int64) (*ZipReader, error) {
	r, err := zip.NewReader(instream, size)
	if err != nil {
		return nil, err
	}

	return &ZipReader{
		r: r,
	}, nil
}

func (r ZipReader) ReadFiles(receivers map[string]io.Writer) error {
	for _, f := range r.r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		w, ok := receivers[f.Name]
		if !ok {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(w, rc); err != nil {
			return err
		}
		rc.Close()
	}

	return nil
}

func (r ZipReader) ListEntries() []*File {
	return r.r.File
}

func (r ZipReader) ListFiles() []*File {
	files := []*zip.File{}
	for _, f := range r.r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		files = append(files, f)
	}

	return files
}
