#!/bin/bash
# Lock or Sleep Computer

ACTION="${1:-lock}" # lock, sleep, shutdown, restart

echo "Action: $ACTION"

case "$ACTION" in
    lock)
        echo "Locking computer..."
        osascript -e 'tell application "System Events" to keystroke "q" using {control down, command down}'
        ;;
    sleep)
        echo "Putting computer to sleep..."
        pmset sleepnow
        ;;
    shutdown)
        echo "Shutting down in 60 seconds..."
        osascript -e 'tell application "System Events" to shut down'
        ;;
    restart)
        echo "Restarting in 60 seconds..."
        osascript -e 'tell application "System Events" to restart'
        ;;
    *)
        echo "Unknown action: $ACTION"
        echo "Usage: $0 [lock|sleep|shutdown|restart]"
        exit 1
        ;;
esac
