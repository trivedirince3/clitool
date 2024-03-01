build:
	go build -o dictionary_lookup .

test:
	go test -v ./...

run:
	MW_API_KEY=${MW_API_KEY} ./dictionary_lookup -word beautiful

clean:
	rm -f dictionary_lookup

.PHONY: clean build test run