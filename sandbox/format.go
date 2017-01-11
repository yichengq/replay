package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func formatReason(code []byte) (*result, error) {
	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, err
	}
	srcFile := path.Join(tmpDir, "hello.re")
	if err := ioutil.WriteFile(srcFile, code, 0755); err != nil {
		return nil, fmt.Errorf("failed to write source file: %+v", err)
	}
	cmd := exec.Command(os.Getenv("REFMT_BIN"), "--parse", "re", "--print", "re", srcFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			errStr := strings.Replace(string(output), srcFile, "prog.re", -1)
			return &result{errStr: errStr}, nil
		}
		return nil, fmt.Errorf("failed to run bsc: %+v, output=%q", err, string(output))
	}
	return &result{output: output}, nil
}
