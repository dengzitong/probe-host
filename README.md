# probe-host

* requires
```
install goland 1.12+
```

* build
```
make && make install
```
* run
```
./bin/probe-host --addr 127.0.0.1 --port 9999 --probe-timeout 1 
```
* test
```
http://localhost:9999/?data=[{ "ip" : "127.0.0.1" , "port" :["3380..3389","6443"] },{ "ip" : "127.0.0.1" , "port" :["8080..8081","9999"] }]
```