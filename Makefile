.PHONY: run
run:
	go build -o run/Util.bin . && cd run && ./Util.bin
