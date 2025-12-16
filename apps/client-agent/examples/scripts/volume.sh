#!/bin/bash
# Control Volume and Brightness

ACTION="${1:-status}" # set, mute, unmute, status
LEVEL="${2:-50}"

case "$ACTION" in
    set)
        echo "Setting volume to $LEVEL%..."
        osascript -e "set volume output volume $LEVEL"
        ;;
    mute)
        echo "Muting..."
        osascript -e "set volume output muted true"
        ;;
    unmute)
        echo "Unmuting..."
        osascript -e "set volume output muted false"
        ;;
    status)
        CURRENT=$(osascript -e "output volume of (get volume settings)")
        echo "Current volume: $CURRENT%"
        ;;
    max)
        echo "Setting volume to MAX..."
        osascript -e "set volume output volume 100"
        say "Volume is now at maximum"
        ;;
    *)
        echo "Usage: $0 [set|mute|unmute|status|max] [level]"
        exit 1
        ;;
esac
