# telephone

telephone is a minimal command line tool for wrapping commands in simple webservers written in Go.

## Example usage

1. Kick off the server and a test command (`python test.py`)
    ```
    telephone --port 8123 python test.py
    ```
2. Pass some input:
    ```
    curl -X POST -d "4" http://127.0.0.1:8123
    ```

The response will be the output of the script `test.py` (i.e., "5").
Note that `test.py` expects to get its inputs from stdin and writes its outputs to stdout.


(If it's not obvious, this is really a toy proof of concept, and probably not the best way to expose your Python programs as RESTful APIs.)