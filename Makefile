.SILENT:

build:
	echo "Building: amt"

	cd src/amt && go build -ldflags="-w -s -buildid=" -trimpath -o ../../bin/amt