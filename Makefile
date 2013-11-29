all:
	go install
	rm -rf yo/data.go
	packgo -packname="main" -varname="INIT_DATA" -output="yo/data.go" defaults/
	#gofmt -l -w -s yo/data.go
	cd yo
	go install
