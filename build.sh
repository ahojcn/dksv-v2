BINPATH="./bin/"

echo "build for windows: 386 amd64"
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ${BINPATH}dksv2.0.win-386.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINPATH}dksv2.0.win-amd64.exe

echo "build for linux: 386 amd64 arm"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINPATH}dksv2.0.linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ${BINPATH}dksv2.0.linux-386
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ${BINPATH}dksv2.0.linux-arm

echo "build for darwin: 386 amd64"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o ${BINPATH}dksv2.0.darwin-386
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BINPATH}dksv2.0.darwin-amd64

echo "go build done, push ${BINPATH} to gitee"

cd ${BINPATH} && cp ../deploy.sh ./

git status
git add .
git commit -m "$(date "+%Y-%m-%d %H:%M:%S")"
git push origin master
