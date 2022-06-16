default:
	@echo "Try 'make demo'"

demo: hidapi
	go run demo/main.go
.PHONY: demo

clean:
	rm -rf hidapi hidapi.zip temp

hidapi:
	curl -L -o hidapi.zip https://github.com/libusb/hidapi/archive/refs/heads/master.zip
	unzip -d temp hidapi.zip
	mv temp/*/ hidapi
	rm -rf temp hidapi.zip
