#!/bin/bash

gcc_images=(
  gcc:4
  gcc:5
  gcc:6
  gcc:7
  gcc:8
  gcc:9
  gcc:10
  gcc:11
  gcc:12
  gcc:13
  gcc:14
  gcc:15
)

# 确定工作目录
workdir="$(dirname "$(dirname "$(dirname "$(realpath ${0})")")")"
output_root="${workdir}/samples/gcov-data"
code_root="${workdir}/samples/hello"

# 确定 docker 命令
_docker="docker"
if which podman > /dev/null 2>&1; then
  _docker="podman"
fi
if ! which "${_docker}" > /dev/null; then
  echo "${_docker} not found"
  exit 1;
fi

echo "Workdir    : ${workdir}"
echo "Output root: ${output_root}"
echo "Code root  : ${code_root}"
echo "Docker     : ${_docker}"

set -e

# 逐个 gcc 生成
for img in ${gcc_images[*]}; do
  output_dir="$(echo "${img}" | sed -nre 's:.*(^|/)([^/]+)$:\2:gp' | sed -re 's/[:@]/_/g')"

  echo "----------------"
  echo "GCC: ${img}"
  echo "Output dir: ${output_dir}"

  rm -rf "${output_root}/${output_dir}"
  mkdir -p "${output_root}/${output_dir}"
  $_docker run \
    -it \
    --platform linux/amd64 \
    --rm \
    -v "${code_root}/src:/workdir/src:ro" \
    -v "${code_root}/CMakeLists.txt:/workdir/CMakeLists.txt:ro" \
    -v "${code_root}/cmake-3.26.6-linux-x86_64.sh:/workdir/cmake-3.26.6-linux-x86_64.sh:ro" \
    -v "${output_root}/${output_dir}:/output/workdir" \
    --workdir "/workdir" \
    "${img}" \
    bash -c 'sh /workdir/cmake-3.26.6-linux-x86_64.sh --skip-license --prefix=/usr/local && export PATH=/usr/local/cmake-3.26.6-linux-x86_64/bin:${PATH} && cmake . && make hello && GCOV_PREFIX=/output bin/hello && find . -name "*.gcno" -exec cp --parents -t /output/workdir {} +'
done
