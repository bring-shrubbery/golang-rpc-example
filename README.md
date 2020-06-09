# Golang RPC Example

This is an example of how to setup/use golang's built-in RPC. This example contains a client, scheduler (for primitive load balancing) and two servers, which are contacted by only the scheduler via RPC.

## How to run

To run the full example you will need to run 4 processes: client, scheduler and 2 servers. You can start all four processes by running following commands separately (in separate terminal windows) in order as they are presented below:

```bash
go run main.go -s1 -p 4001
```

```bash
go run main.go -s2 -p 4002
```

```bash
go run main.go -scheduler
```

```bash
go run main.go
```
