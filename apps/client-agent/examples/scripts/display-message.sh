#!/bin/bash
# Display Full Screen Message

TITLE="${1:-Important Message}"
MESSAGE="${2:-You have a new message from Thinhnx!}"
BUTTON="${3:-OK}"

# Show dialog with message
osascript <<EOF
display dialog "$MESSAGE" with title "$TITLE" buttons {"$BUTTON"} default button "$BUTTON" with icon caution giving up after 30
EOF

# Examples:
# bash display-message.sh "Alert" "Meeting in 5 minutes!" "Got it"
# bash display-message.sh "Birthday" "Happy Birthday! 🎉" "Thanks"
