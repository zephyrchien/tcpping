# Tcpping
## Usage
```
tcpping -h <host> -p <port> -c <count> -t <timeout> -q[quiet]
```
# Example
```
tcpping -h github.com -p 8080

TCPPING github.com (13.250.177.223):
seq   1: 13.250.177.223:8080[close] 1000ms
seq   2: 13.250.177.223:8080[close] 1000ms
seq   3: 13.250.177.223:8080[close] 1000ms
seq   4: 13.250.177.223:8080[close] 1000ms
^C----------
total: 4
min/avg/max = 1000/1000/1000ms
```
# PortScan has been moved to [pscan](https://github.com/ZephyrChien/pscan)

