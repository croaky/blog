# Circular Favicon from GitHub Avatar

I wrote a script to generate a circular favicon and Apple Touch icon
from my GitHub avatar and generalized it.

Make sure ImageMagick is installed. On macOS:

```
brew install imagemagick
```

Copy the script into a file named `favicon.sh`:

```embed
code/favicon.sh all
```

Make it executable:

```
chmod +x favicon.sh
```

Run it with your GitHub username:

```
./favicon.sh username
```

These files will be created:

```
favicon.ico
apple-touch-icon.png
```
