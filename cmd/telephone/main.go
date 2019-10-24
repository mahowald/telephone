package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "telephone"
	app.Usage = "Wrap the specified application with a simple webserver"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "port for the webserver",
		},
	}

	app.Action = func(c *cli.Context) error {
		cmd := exec.Command(c.Args().First(), c.Args().Tail()...)

		// to support streaming writes to stdin and reads from stdout
		// use StdinPipe/StdoutPipe instead of just cmd.Stdin/cmd.Stdout
		cmdWriter, _ := cmd.StdinPipe()
		cmdReader, _ := cmd.StdoutPipe()

		// four channels:
		// inchan/outchan for actually supplying data to/from the command
		// and echo variants for logging purposes
		inchan := make(chan string)
		inchanecho := make(chan string)
		outchan := make(chan string)
		outchanecho := make(chan string)

		// goroutine for grabbing cmd's stdout and putting it
		// onto outchan and outchanecho
		go func() {
			scanner := bufio.NewScanner(cmdReader)
			for scanner.Scan() {
				msg := scanner.Text()
				outchan <- msg
				outchanecho <- msg
			}
		}()

		// goroutine for writing anything put on inchan
		// to stdin for cmd (and copying that input to inchanecho)
		go func() {
			for msg := range inchan {
				io.WriteString(cmdWriter, msg+"\n")
				inchanecho <- msg
			}
		}()

		// kick off the command
		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		// goroutine to do basic logging of inputs & outputs
		go func() {
			for msg := range inchanecho {
				res := <-outchanecho
				fmt.Println(msg, "->", res)
			}
		}()

		// set up a simple webserver
		http.HandleFunc("/", makeRequestHandler(inchan, outchan))
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.String("port")), nil))
		}()

		err = cmd.Wait()
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// makeRequestHandler returns a handler function that will copy inputs recieved from POST requests
// onto inchan, and then write corresponding outputs on outchan as a response to the request.
func makeRequestHandler(inputs chan string, outputs chan string) func(w http.ResponseWriter, r *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			inputs <- string(body)
			resp := <-outputs
			fmt.Fprint(w, resp+"\n")
		}
	}

	return f
}
