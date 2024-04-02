#!/bin/bash
set -euo pipefail

mkdir -p tmp
pg_dump -Fc "$(cb uri app-production)" > tmp/latest.backup
