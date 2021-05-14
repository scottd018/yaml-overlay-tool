// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func ReadStream(fileName string) (io.Reader, error) {
	if fileName == "-" {
		return bufio.NewReader(os.Stdin), nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to %w", err)
	}

	return file, nil
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Error("error closing file!: %s", err)
	}
}

func GetCommonPrefix(sep byte, paths ...string) string {
	switch len(paths) {
	case 0:
		return ""
	case 1:
		return path.Clean(paths[0])
	}

	c := path.Clean(paths[0]) + string(sep)

	for _, v := range paths[1:] {
		v = path.Clean(v) + string(sep)

		if len(v) < len(c) {
			c = c[:len(v)]
		}

		for i := 0; i < len(c); i++ {
			if v[i] != c[i] {
				c = c[:i]

				break
			}
		}
	}

	for i := len(c) - 1; i >= 0; i-- {
		if c[i] == sep {
			c = c[:i]

			break
		}
	}

	return c
}

func isDirectory(p string) (bool, error) {
	fileInfo, err := os.Stat(p)
	if err != nil {
		return false, fmt.Errorf("could not get file info for %s, %w", p, err)
	}

	return fileInfo.IsDir(), err
}

func getFileNames(p string) ([]string, error) {
	var results []string

	err := filepath.Walk(p,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
				results = append(results, path)
			}

			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("could not walk directory %s, %w", p, err)
	}

	return results, nil
}
