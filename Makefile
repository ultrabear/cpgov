GOCMD=go
EXE=cpgov

${EXE}: *.go
	${GOCMD} build -o ${EXE}

install: ${EXE}
	cp ./${EXE} /bin/${EXE}

install-setuid: install
	chown root.root /bin/${EXE}
	chmod 4711 /bin/${EXE}
