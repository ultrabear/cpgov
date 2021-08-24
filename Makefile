GOCMD=go
EXE=cpgov

${EXE}: *.go
	${GOCMD} build -o ${EXE}

install: ${EXE}
	cp ./${EXE} /usr/bin/${EXE}

install-setuid: install
	chown root.root /usr/bin/${EXE}
	chmod 4711 /usr/bin/${EXE}
