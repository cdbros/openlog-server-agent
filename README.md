# Openlog server agent

Golang program used to loop through a log folder and send csv files to openlog server

## How to use it
- Build the executable with ```GOOS=target-OS GOARCH=target-architecture go build ```
- Move the file to target server
- ```./executable "/path/to/log/folder" "OPENLOG_HOST" ```

## Author
- Danilo Cadeddu (danilo_cadeddu@outlook.it)