VIDEO_URL="${1:-https://kaiko-varmeta-dev.s3.ap-southeast-1.amazonaws.com/7d47ab9b5729a33c1c064858bc18aa85.mp4}" # Can be YouTube URL, direct video URL, or local path
VIDEO_PATH="/tmp/birthday-video.mp4"
if [ -f "$1" ]; then
    VIDEO_PATH="$1"
else
    if [[ "$VIDEO_URL" == *.mp4 ]] || [[ "$VIDEO_URL" == *.mov ]] || [[ "$VIDEO_URL" == *.avi ]]; then
        curl -L "$VIDEO_URL" -o "$VIDEO_PATH"
    else
        VIDEO_PATH="$VIDEO_URL"
    fi
fi
osascript -e 'display notification "🎂 Happy Birthday! 🎉" with title "Special Surprise" sound name "Glass"'
sleep 1
echo "Opening birthday video..."
if [ -f "$VIDEO_PATH" ]; then
    osascript <<EOF
tell application "QuickTime Player"
    activate
    open POSIX file "$VIDEO_PATH"
    tell front document
        present
        play
    end tell
end tell
EOF
else
    open "$VIDEO_PATH"
fi