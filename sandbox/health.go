package main

import "fmt"

func healthCheck() error {
	output, err := compile([]byte(healthReason), compileToJS)
	if err != nil {
		return err
	}
	if string(output) != wantHealthJs {
		return fmt.Errorf("compiled JS = %q, want %q", string(output), wantHealthJs)
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
