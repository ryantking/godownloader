package main

const shellfn = `

cat /dev/null << EOF
------------------------------------------------------------------------
https://github.com/client9/posixshell - portable posix shell functions
Public domain - http://unlicense.org
https://github.com/client9/posixshell/blob/master/LICENSE.md
but credits (and pull requests) appreciated.
------------------------------------------------------------------------
EOF
is_command() {
  type $1 > /dev/null 2> /dev/null
}
uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    x86)    arch="386" ;;
    i686)   arch="386" ;;
    i386)   arch="386" ;;
  esac
  echo ${arch}
}
uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  echo ${os}
}
untar() {
  tarball=$1
  case ${tarball} in
  *.tar.gz|*.tgz) tar -xzf ${tarball} ;;
  *.tar) tar -xf ${tarball} ;;
  *.zip) unzip ${tarball} ;;
  *)
    echo "Unknown archive format for ${tarball}"
    return 1
  esac
}
mktmpdir() {
   test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
   mkdir -p ${TMPDIR}
   echo ${TMPDIR}
}
http_download() {
  local_file=$1
  source_url=$2
  header=$3
  headerflag=''
  destflag=''
  if is_command curl; then
    cmd='curl --fail -sSL'
    destflag='-o'
    headerflag='-H'
  elif is_command wget; then
    cmd='wget -q'
    destflag='-O'
    headerflag='--header'
  else
    echo "http_download: unable to find wget or curl"
    return 1
  fi
  if [ -z "$header" ]; then
    $cmd $destflag "$local_file" "$source_url"
  else
    $cmd $headerflag "$header" $destflag "$local_file" "$source_url"
  fi
}
github_api() {
  local_file=$1
  source_url=$2
  header=""
  case "$source_url" in
  https://api.github.com*)
     test -z "$GITHUB_TOKEN" || header="Authorization: token $GITHUB_TOKEN"
     ;;
  esac
  http_download "$local_file" "$source_url" "$header"
}
github_last_release() {
  OWNER_REPO=$1
  VERSION=$(github_api - https://api.github.com/repos/${OWNER_REPO}/releases/latest | grep -m 1 "\"name\":" | cut -d ":" -f 2 | tr -d ' ",')
  if [ -z "${VERSION}" ]; then
    echo "Unable to determine latest release for ${OWNER_REPO}"
    return 1
  fi
  echo ${VERSION}
}
hash_sha256() {
  TARGET=${1:-/dev/stdin};
  if is_command gsha256sum; then
    hash=$(gsha256sum $TARGET) || return 1
    echo $hash | cut -d ' ' -f 1
  elif is_command sha256sum; then
    hash=$(sha256sum $TARGET) || return 1
    echo $hash | cut -d ' ' -f 1
  elif is_command shasum; then
    hash=$(shasum -a 256 $TARGET 2>/dev/null) || return 1
    echo $hash | cut -d ' ' -f 1
  elif is_command openssl; then
    hash=$(openssl -dst openssl dgst -sha256 $TARGET) || return 1
    echo $hash | cut -d ' ' -f a
  else
    echo "hash_sha256: unable to find command to compute sha-256 hash"
    return 1
  fi
}
hash_sha256_verify() {
  TARGET=$1
  checksums=$2
  if [ -z "$checksums" ]; then
     echo "hash_sha256_verify: checksum file not specified in arg2"
     return 1
  fi
  BASENAME=${TARGET##*/}
  want=$(grep ${BASENAME} "${checksums}" 2> /dev/null | tr '\t' ' ' | cut -d ' ' -f 1)
  if [ -z "$want" ]; then
     echo "hash_sha256_verify: unable to find checksum for '${TARGET}' in '${checksums}'"
     return 1
  fi
  got=$(hash_sha256 $TARGET)
  if [ "$want" != "$got" ]; then
     echo "hash_sha256_verify: checksum for '$TARGET' did not verify ${want} vs $got"
     return 1
  fi
}
cat /dev/null << EOF
------------------------------------------------------------------------
End of functions from https://github.com/client9/posixshell 
------------------------------------------------------------------------
EOF
`