#!/bin/bash
# 6x:
echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org' |\
go run gogogo.go
# 5x:
echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org' |\
go run gogogo.go
# 2x:
echo -e 'https://golang.org\nhttps://golang.org' |\
go run gogogo.go
# 1x:
echo -e 'https://golang.org' |\
go run gogogo.go
