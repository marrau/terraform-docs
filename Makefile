
V=`git describe --tags --always`
B="-X main.version=$(V)"

CGO_ENABLED=0 

dist:
	@gox \
		--os "linux" \
		--output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		--ldflags=$(B)

clean:
	rm -rf dist/

install:
	go install
