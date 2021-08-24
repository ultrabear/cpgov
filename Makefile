GOCMD=go
EXE=cpgov

${EXE}: *.go
	${GOCMD} build -o ${EXE}

install: ${EXE}
	cp ./${EXE} /usr/sbin/${EXE}

install-setuid: install
	chown root.root /usr/sbin/${EXE}
	chmod 4711 /usr/sbin/${EXE}
