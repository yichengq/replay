package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type compileType int

const (
	compileToJS compileType = iota + 1
	compileToML
	compileToRun
)

func parseCompileType(str string) (compileType, error) {
	switch str {
	case "to_js":
		return compileToJS, nil
	case "to_ml":
		return compileToML, nil
	case "to_run":
		return compileToRun, nil
	default:
		return 0, errors.New("invalid input")
	}
}

func compile(code []byte, typ compileType) (*result, error) {
	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, err
	}
	srcFile := path.Join(tmpDir, "prog.re")
	if err := ioutil.WriteFile(srcFile, code, 0755); err != nil {
		return nil, fmt.Errorf("failed to write source file: %+v", err)
	}

	switch typ {
	case compileToJS:
		destFile := path.Join(tmpDir, "prog.js")
		cmd := exec.Command(os.Getenv("BSC_BIN"), "-pp", os.Getenv("REFMT_BIN"), "-impl", srcFile)
		if output, err := cmd.CombinedOutput(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				errStr := strings.Replace(string(output), srcFile, "prog.re", -1)
				return &result{errStr: errStr}, nil
			}
			return nil, fmt.Errorf("failed to run bsc: %+v, output=%q", err, string(output))
		}
		bs, err := ioutil.ReadFile(destFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to read file: %+v", err)
		}
		return &result{output: bs}, nil
	case compileToRun:
		destFile := path.Join(tmpDir, "prog.native")
		cmd := exec.Command(os.Getenv("REBUILD_BIN"), "prog.native")
		cmd.Dir = tmpDir
		if output, err := cmd.CombinedOutput(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				return &result{errStr: polishReasonBinsOutput(output)}, nil
			}
			return nil, fmt.Errorf("failed to compile: %+v, output=%q", err, string(output))
		}

		cmd = exec.Command(destFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				return &result{errStr: polishReasonBinsOutput(output)}, nil
			}
			return nil, fmt.Errorf("failed to run: %+v, output=%q", err, string(output))
		}
		return &result{output: output}, nil
	default:
		// unsupported yet
		return nil, errors.New("unsupported compile type")
	}
}
