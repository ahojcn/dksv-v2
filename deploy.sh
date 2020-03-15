LINUX_X86_64="dksv2.0.linux-amd64"
LINUX_X86_32="dksv2.0.linux-386"
LINUX_ARM="dksv2.0.linux-arm"
DARWIN_X86_64="dksv2.0.darwin-amd64"
DARWIN_X86_32="dksv2.0.darwin-386"
WIN_X86_64="dksv2.0.win-amd64.exe"
WIN_X86_32="dksv2.0.win-386.exe"

SERVER_ROOT_URL="https://ahojcn.gitee.io/"
ROOT_PATH="/root/bin/"

SYSNAME=`uname -s`
SYSLONG=`uname -m`

mkdir ${ROOT_PATH} -p
cd ${ROOT_PATH}

echo "${1}" > "ip.txt"  # 保存 ip
echo "${2}" > "apiversion.txt"  # 保存 api 版本
echo "${3}" > "port.txt"  # 保存 api 端口

function downloadDksv() {
  curl "${SERVER_ROOT_URL}${1}" --output "${1}"
  chmod +x "${1}"
  curl "${SERVER_ROOT_URL}${2}" --output "${2}"
  chmod +x "${2}"
  "${2}"
  # nohup "${1}" &
  # "${1}"
}

if [ ${SYSNAME} == "Linux" ]; then
    if [ ${SYSLONG} == "x86_64" ]; then
      downloadDksv ${LINUX_X86_64} "deploy-linux-amd64"
    fi
fi

if [ ${SYSNAME} == "Darwin" ]; then
    if [ ${SYSLONG} == "x86_64" ]; then
      downloadDksv ${DARWIN_X86_64}
    fi
fi
