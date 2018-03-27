go-shell - execute shell command in parallel
=======

How to run

```bash
$ go-shell -h
Usage of go-shell:
  -l	loop the commands
  -n int
    	number of processes (default 5)
  -s int
    	sleep beween tasks per process, in seconds (default 10)

```

```bash
$ go-shell -n 2 << EOF
echo hello
echo world
EOF
```
or 
```bash
printf "echo hello\necho world" | go-shell -n 2
```

or loop the commands and run infinitely
```bash
echo date | go-shell -n 1 -s 1 -l
```