
.PHONY: zia clean test_target_server saypcap

clean:
	rm -rf ./build/ 2>/dev/null 2>&1
	mkdir build

saypcap:
	echo "setcap CAP_NET_BIND_SERVICE=+eip /tmp/zia"

zia: saypcap
	GOOS=linux GOARCH=amd64 go build -a -o build/zia -ldflags "-s -w" zia.go; chmod +x build/zia

test_target_server:
	GOOS=linux GOARCH=amd64 go build -a -o build/test_target_server -ldflags "-s -w" testTargetServer/test_target_server.go; chmod +x build/test_target_server

zip: clean zia test_target_server
	cp zia.service build/
	cd build && zip -r ../zia_reverse_proxy.zip *

cert:
	cd config/cert && openssl genrsa -out ziaca.key 4096
	cd config/cert && openssl req -new -x509 -days 3650 -key ziaca.key -out ziacert.pem -subj "/C=RO/ST=AR/L=Arad/O=Zia/OU=IT/CN=ZiaRootCA"

install:
	sudo mkdir -p /opt/zia/
	sudo chmod +x build/zia
	sudo setcap CAP_NET_BIND_SERVICE=+eip build/zia
	cp -fr static/ /opt/zia && cp -fr templates/ /opt/zia && cp -f zia /opt/zia/
	endif
