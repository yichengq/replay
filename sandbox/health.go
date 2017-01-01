package main

import (
	"errors"
	"fmt"
)

func healthCheck() error {
	res, err := compile([]byte(healthReason), compileToJS)
	if err != nil {
		return err
	}
	if res.errStr != "" {
		return errors.New(res.errStr)
	}
	if string(res.output) != wantHealthJs {
		return fmt.Errorf("compiled JS = %q, want %q", string(res.output), wantHealthJs)
	}
	return nil
}

const healthReason = `
Js.log "HELLO, REASON & BUCKLESCRIPT"
`

const wantHealthJs = `// GENERATED CODE BY BUCKLESCRIPT VERSION 0.7.1 , PLEASE EDIT WITH CARE
'use strict';


console.log("HELLO, REASON & BUCKLESCRIPT");

/*  Not a pure module */
`
