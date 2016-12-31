package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type compileType int

const (
	compileToJS compileType = iota + 1
	compileToML
)

func compile(code []byte, typ compileType) ([]byte, error) {
	if typ == compileToML {
		// unsupported yet
		return nil, errors.New("unsupported compile type")
	}

	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, err
	}
	srcFile := path.Join(tmpDir, "hello.re")
	destFile := path.Join(tmpDir, "hello.js")
	if err := ioutil.WriteFile(srcFile, code, 0755); err != nil {
		return nil, fmt.Errorf("failed to write source file: %+v", err)
	}
	cmd := exec.Command(os.Getenv("BSC_BIN"), "-pp", os.Getenv("REFMT_BIN"), "-impl", srcFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to run bsc: %+v, output=%q", err, string(output))
	}
	bs, err := ioutil.ReadFile(destFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %+v", err)
	}
	return bs, nil
}
