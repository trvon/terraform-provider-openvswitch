#!/usr/bin/env bash

set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

ver="${1:-0.0.1}"
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch_raw="$(uname -m)"

case "${arch_raw}" in
  x86_64) arch="amd64" ;;
  aarch64) arch="arm64" ;;
  arm64) arch="arm64" ;;
  *) echo "unsupported arch: ${arch_raw}" >&2; exit 1 ;;
esac

zip_name="terraform-provider-openvswitch_${ver}_${os}_${arch}.zip"

mirror_base="${root_dir}/.provider-mirror"
tf_host_dir="${mirror_base}/registry.terraform.io/trvon/openvswitch"
tofu_host_dir="${mirror_base}/registry.opentofu.org/trvon/openvswitch"

mkdir -p "${tf_host_dir}" "${tofu_host_dir}"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "${tmp_dir}"' EXIT

cp "${root_dir}/bin/terraform-provider-openvswitch" "${tmp_dir}/terraform-provider-openvswitch_v${ver}"
chmod +x "${tmp_dir}/terraform-provider-openvswitch_v${ver}"

(cd "${tmp_dir}" && zip -q "${zip_name}" "terraform-provider-openvswitch_v${ver}")

cp -f "${tmp_dir}/${zip_name}" "${tf_host_dir}/${zip_name}"
cp -f "${tmp_dir}/${zip_name}" "${tofu_host_dir}/${zip_name}"

echo "wrote ${tf_host_dir}/${zip_name}"
echo "wrote ${tofu_host_dir}/${zip_name}"
