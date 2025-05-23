# mov2gif

I record screencasts with QuickTime.app, which produces a `.mov` file.
When I want a smaller file to share in product documentation or this blog:

```bash
mov2gif input.mov
```

`mov2gif` takes a `.mov` file as input and outputs a `.output.gif`.

Script:

```bash
#!/bin/bash
#
# Convert a .mov file to a .gif file
#
# Usage: mov2gif [-f fps] [-s scale] input.mov
#
# Options:
#   -f fps     Frames per second (default: 15)
#   -s scale   Scale width in pixels (default: 1400)

set -euo pipefail

# Default parameters
fps=15     # Frames per second
scale=1400 # Scale width in pixels

# Function to display usage information
usage() {
  echo "Usage: $0 [-f fps] [-s scale] input.mov"
  echo "  -f fps     Frames per second (default: 15)"
  echo "  -s scale   Scale width in pixels (default: 1400)"
  exit 1
}

# Parse options
while getopts ":f:s:" opt; do
  case $opt in
  f)
    fps="$OPTARG"
    ;;
  s)
    scale="$OPTARG"
    ;;
  \?)
    echo "Invalid option: -$OPTARG" >&2
    usage
    ;;
  :)
    echo "Option -$OPTARG requires an argument." >&2
    usage
    ;;
  esac
done

# Shift parsed options away
shift $((OPTIND - 1))

# Check for input file
if [ $# -ne 1 ]; then
  echo "Error: Missing input file."
  usage
fi

input_file="$1"

# Verify input file extension
if [[ "${input_file##*.}" != "mov" ]]; then
  echo "Error: Input file must have a .mov extension."
  exit 1
fi

# Check if input file exists
if [ ! -f "$input_file" ]; then
  echo "Error: File '$input_file' not found."
  exit 1
fi

# Check and install ffmpeg if not present
if ! command -v ffmpeg &>/dev/null; then
  echo "ffmpeg not found. Installing via Homebrew..."
  brew install ffmpeg
fi

# Generate a palette for better colors in the GIF
echo "Regenerating palette..."
rm palette.png
ffmpeg -y -i "$input_file" -vf "fps=${fps},scale=${scale}:-1:force_original_aspect_ratio=decrease,palettegen" -frames:v 1 palette.png

# Create the final GIF using the palette
echo "Creating GIF..."
ffmpeg -i "$input_file" -i palette.png -filter_complex "fps=${fps},scale=${scale}:-1:force_original_aspect_ratio=decrease[x];[x][1:v]paletteuse=dither=none" output.gif

echo "GIF created: output.gif"
```
