.SILENT:

build:
	echo "Building: amt - Reconnaissance tool"

	cd src/amt && go build -ldflags="-w -s -buildid=" -trimpath -o ../../bin/amt