package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultOriginFolderName = "origin"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:]) // PERFORMANCE:

	blocksize := 5
	sliceLen := len(hashString) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashString[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashString,
	}
}

type PathTransformFunc func(string) PathKey

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type StoreOpts struct {
	// Origin is the origin for all files on the system, it is the root
	Origin            string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if opts.Origin == "" {
		opts.Origin = defaultOriginFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithOrigin := fmt.Sprintf("%s/%s", s.Origin, pathKey.FullPath())

	_, err := os.Stat(fullPathWithOrigin)
	return !errors.Is(err, fs.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Origin)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	firstPathNameWithOrigin := fmt.Sprintf("%s/%s", s.Origin, pathKey.FirstPathName())

	return os.RemoveAll(firstPathNameWithOrigin)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithOrigin := fmt.Sprintf("%s/%s", s.Origin, pathKey.FullPath())
	return os.Open(fullPathWithOrigin)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithOrigin := fmt.Sprintf("%s/%s", s.Origin, pathKey.Pathname)

	if err := os.MkdirAll(pathNameWithOrigin, os.ModePerm); err != nil {
		return err
	}

	// fullPath := pathKey.FullPath()
	fullPathWithOrigin := fmt.Sprintf("%s/%s", s.Origin, pathKey.FullPath())

	f, err := os.Create(fullPathWithOrigin)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk:\n%s", n, fullPathWithOrigin)

	return nil
}
