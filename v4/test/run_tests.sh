#docker-compose -f test/docker-compose.yaml run --rm -d
echo "Unit tests"
cd ..
go test -v ./weaviate/...
echo "Integration tests"
for pkg in $(go list ./... | grep 'weaviate-go-client/v4/test'); do if ! go test -v -count 1 -race "$pkg"; then echo "Test for $pkg failed" >&2; false; exit; fi; done 
#docker-compose -f test/docker-compose.yaml down
#docker-compose -f test/docker-compose.yaml delete
