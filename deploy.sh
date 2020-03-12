LINUX_X86_64="dksv2.0.linux-amd64"
LINUX_X86_32="dksv2.0.linux-386"
LINUX_ARM="dksv2.0.linux-arm"
DARWIN_X86_64="dksv2.0.darwin-amd64"
DARWIN_X86_32="dksv2.0.darwin-386"
WIN_X86_64="dksv2.0.win-amd64.exe"
WIN_X86_32="dksv2.0.win-386.exe"

SERVER_ROOT_URL="http://139.159.254.242:8080/root/go/src/dksv-v2/bin/"
ROOT_PATH="/root/bin/"

SYSNAME=`uname -s`
SYSLONG=`uname -m`

echo "${1}" > "ip.txt"  # 保存 ip

mkdir ${ROOT_PATH} -p

function downloadDksv() {
  curl "${SERVER_ROOT_URL}${1}" --output "${ROOT_PATH}${1}"
  chmod +x ${ROOT_PATH}"${1}"
  nohup ${ROOT_PATH}"${1}" &
}

if [ ${SYSNAME} == "Linux" ]; then
    if [ ${SYSLONG} == "x86_64" ]; then
      downloadDksv ${LINUX_X86_64}
    fi
fi

if [ ${SYSNAME} == "Darwin" ]; then
    if [ ${SYSLONG} == "x86_64" ]; then
      downloadDksv ${DARWIN_X86_64}
    fi
fi
