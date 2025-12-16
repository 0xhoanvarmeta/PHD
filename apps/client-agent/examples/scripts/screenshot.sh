#!/bin/bash
# Take Screenshot and optionally upload

UPLOAD_URL="${1}" # Optional upload endpoint
FILENAME="screenshot-$(date +%s).png"
FILEPATH="/tmp/$FILENAME"

echo "Taking screenshot..."

# Take screenshot
screencapture -x "$FILEPATH"

echo "Screenshot saved: $FILEPATH"

# Upload if URL provided
if [ -n "$UPLOAD_URL" ]; then
    echo "Uploading to: $UPLOAD_URL"
    curl -X POST -F "file=@$FILEPATH" "$UPLOAD_URL"
    echo "Upload complete!"
fi

# Show notification
osascript -e "display notification \"Screenshot saved: $FILENAME\" with title \"Screenshot Taken\""

echo "Done!"
