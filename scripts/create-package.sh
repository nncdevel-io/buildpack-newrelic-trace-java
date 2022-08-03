set -euo pipefail

export OUTPUT=buildpack

if [[ -e "${OUTPUT}" ]]; then
  rm -rf "${OUTPUT}"
fi

mkdir "${OUTPUT}"

cp -r bin "${OUTPUT}"
cp -r buildpack.toml "${OUTPUT}"
