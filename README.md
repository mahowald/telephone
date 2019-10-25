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

## Containerized version

1. Build the image:
    ```
    docker build -t mahowald/telephone .
    ```
2. Run the image:
    ```
    docker run -it --net=host mahowald/telephone --port 8123 python test.py
    ```
3. Pass some input:
    ```
    curl -X POST -d "5" http://localhost:8123
    ```



(If it's not obvious, this is a toy proof of concept, and probably not the best way to expose your Python programs as RESTful APIs.)