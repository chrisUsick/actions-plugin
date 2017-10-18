SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"
OS="darwin linux"

gox -arch="amd64" -os="$OS" -output "./bin/{{.OS}}/actions-plugin" ./

NAME=vault-integration
docker rm -f $NAME || true
shasum -a 256 -b bin/actions-plugin | awk '{print $1 }' > bin/checksum
for GOOS in $OS; do 
    shasum -a 256 -b bin/$GOOS/actions-plugin | awk '{print $1 }' > bin/$GOOS/checksum
done
docker build -t chrisusick/vault-integration .
ID=$(docker run --cap-add=IPC_LOCK -d --name $NAME -p 8201:8200 -v $PWD/.logs:/vault/logs chrisusick/vault-integration)
echo Container ID=$ID
sleep 2
curl -H 'X-VAULT-TOKEN: 1234' -X PUT http://127.0.0.1:8201/v1/sys/plugins/catalog/actions-plugin -d "{\"sha_256\":\"$(cat bin/linux/checksum)\", \"command\":\"actions-plugin\"}"
docker exec $ID vault auth 1234
docker exec $ID vault mount -path=actions -plugin-name=actions-plugin plugin