#!/usr/bin/env bash
set -euo pipefail

# Generate OpenAPI per api submodule and bundle into single openapi/openapi.yaml
# Requirements: buf, swagger-cli, yq

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
API_DIR="$ROOT_DIR/api"
OUT_DIR="$ROOT_DIR"
TMP_DIR="$ROOT_DIR/tmp_openapi"

echo "API dir: $API_DIR"
echo "Output dir: $OUT_DIR"

command_exists() { command -v "$1" >/dev/null 2>&1; }

for dep in buf swagger-cli; do
  if ! command_exists "$dep"; then
    echo "ERROR: required tool '$dep' not found in PATH. Please install it." >&2
    echo "  - buf: https://docs.buf.build/" >&2
    echo "  - swagger-cli: npm i -g @apidevtools/swagger-cli" >&2
    exit 2
  fi
done

# detect openapi merge tool (openapi-merge-cli)
MERGE_CMD=""
if command_exists openapi-merge-cli; then
  MERGE_CMD=openapi-merge-cli
elif command_exists openapi-merge; then
  MERGE_CMD=openapi-merge
fi
if [[ -z "$MERGE_CMD" ]]; then
  echo "ERROR: required tool 'openapi-merge-cli' not found in PATH. Please install it via npm:" >&2
  echo "  npm i -g openapi-merge-cli" >&2
  exit 2
fi

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR"

# Discover modules: look for directories api/*/* that contain .proto files
modules=()
while IFS= read -r -d '' f; do
  # f is path to proto file e.g. /.../api/cascade/public/auth/v1/auth.proto
  rel="${f#$API_DIR/}"
  # take first two path components as module key
  module_dir=$(echo "$rel" | awk -F/ '{print $1"/"$2}')
  modules+=("$module_dir")
done < <(find "$API_DIR" -name '*.proto' -print0)

# unique
IFS=$'\n' read -r -d '' -a modules_unique < <(printf "%s\n" "${modules[@]}" | awk '!seen[$0]++' && printf '\0')

if [[ ${#modules_unique[@]} -eq 0 ]]; then
  echo "No proto files found under $API_DIR" >&2
  exit 0
fi

echo "Found modules:"; printf "%s\n" "${modules_unique[@]}"

cd "$API_DIR"

for mod in "${modules_unique[@]}"; do
  # module name like cascade/public -> dir_name cascade_public
  mod_dir_name=$(echo "$mod" | tr '/' '_' )
  mod_out="$TMP_DIR/$mod_dir_name"
  mkdir -p "$mod_out"

  echo "Generating OpenAPI for module: $mod -> $mod_out"

  # create temporary buf.gen.yaml based on api/buf.gen.openapi.yaml but override out
  tmp_template="$TMP_DIR/buf.gen.$mod_dir_name.yaml"
  awk -v out="../tmp_openapi/$mod_dir_name" 'BEGIN{p=1} { print }' "$API_DIR/buf.gen.openapi.yaml" > "$tmp_template"
  # replace out: ../tmp_openapi with out: ../tmp_openapi/<mod_dir_name>
  sed -i "s|out: ../tmp_openapi|out: ../tmp_openapi/$mod_dir_name|g" "$tmp_template"

  # run buf generate limited to this module via --path
  echo "Running: buf generate --template $tmp_template --path $mod"
  buf generate --template "$tmp_template" --path "$mod" || { echo "buf generate failed for $mod" >&2; exit 3; }

  # locate generated file (openapi.yaml or openapi.json)
  gen_file=""
  if [[ -f "$mod_out/openapi.yaml" ]]; then
    gen_file="$mod_out/openapi.yaml"
  elif [[ -f "$mod_out/openapi.json" ]]; then
    gen_file="$mod_out/openapi.json"
  else
    # try find
    found=$(find "$mod_out" -maxdepth 2 -type f \( -iname "openapi.*" -o -iname "swagger.*" \) | head -n1 || true)
    if [[ -n "$found" ]]; then
      gen_file="$found"
    fi
  fi

  if [[ -z "$gen_file" ]]; then
    echo "Warning: no generated OpenAPI file found for module $mod in $mod_out" >&2
    continue
  fi

  echo "Generated fragment: $gen_file"

  # normalize/bundle using swagger-cli
  bundled="$mod_out/openapi.bundled.yaml"
  echo "Bundling with swagger-cli: $gen_file -> $bundled"
  swagger-cli bundle "$gen_file" --outfile "$bundled" --type yaml || { echo "swagger-cli bundle failed for $gen_file" >&2; exit 4; }
done

echo "Merging bundled modules into single $TMP_DIR/openapi.yaml"
bundled_files=("$TMP_DIR"/*/openapi.bundled.yaml)
if [[ ${#bundled_files[@]} -eq 0 ]]; then
  echo "No bundled OpenAPI fragments to merge" >&2
  exit 0
fi

echo "Merging ${#bundled_files[@]} bundled files with $MERGE_CMD"
cd "$ROOT_DIR"
"$MERGE_CMD" -c openapi-merge.json || { echo "ERROR: $MERGE_CMD failed to merge files" >&2; exit 5; }
echo "Final OpenAPI: $OUT_DIR/openapi.yaml"

echo "Cleaning up temporary files"
rm -rf "$TMP_DIR"