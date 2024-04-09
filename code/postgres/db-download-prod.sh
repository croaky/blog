#!/bin/bash
set -euo pipefail

# Delete/create target directory
backup_dir="tmp/latest_backup_dir"
rm -rf "$backup_dir"
mkdir -p "$backup_dir"

# Detect the number of CPU cores
case "$(uname -s)" in
    Linux*)     cores=$(nproc);;
    Darwin*)    cores=$(sysctl -n hw.ncpu);;
    *)          cores=1;;
esac

# Use one less than the total number of cores, but ensure at least 1 is used
(( jobs = cores - 1 ))
if (( jobs < 1 )); then
    jobs=1
fi

echo "Downloading with $jobs parallel job(s)"

# Use the directory format and specify the number of jobs for parallel dumping
pg_dump -Fd "$(cb uri app-prod --role application)" -j "$jobs" -f "$backup_dir"
