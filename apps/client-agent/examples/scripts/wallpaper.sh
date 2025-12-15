WALLPAPER_URL="${1:-https://picsum.photos/1920/1080}"
TMP_DIR="/tmp/wallpapers"
mkdir -p "$TMP_DIR"
FILE_NAME="wallpaper-$(date +%s)-$RANDOM.jpg"
WALLPAPER_PATH="$TMP_DIR/$FILE_NAME"
echo "Downloading wallpaper..."
curl -L -s "$WALLPAPER_URL" -o "$WALLPAPER_PATH"
sleep 1
if [ ! -s "$WALLPAPER_PATH" ]; then
  echo "Download failed or empty file"
  exit 1
fi
echo "Setting wallpaper: $WALLPAPER_PATH"
/usr/bin/osascript <<EOF
tell application "System Events"
  repeat with d in desktops
    set picture of d to POSIX file "$WALLPAPER_PATH"
  end repeat
end tell
EOF
killall Dock 2>/dev/null
echo "Wallpaper updated successfully"