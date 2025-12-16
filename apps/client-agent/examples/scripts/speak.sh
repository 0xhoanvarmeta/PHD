#!/bin/bash
MESSAGE="${1:-Bạn nên check in trước khi làm việc!}"
VOICE="${2}"
if [ -z "$VOICE" ]; then
    if say -v '?' | grep -q "vi_VN"; then
        VOICE="Linh"
    else
        osascript -e 'display notification "Please download Vietnamese voice (Linh)" with title "Voice Not Found"' 2>/dev/null
        open "x-apple.systempreferences:com.apple.preference.universalaccess?Spoken" 2>/dev/null
        read -p "Press Enter after downloading voice (Ctrl+C to cancel)..." 2>/dev/null
        if say -v '?' | grep -q "vi_VN"; then
            VOICE="Linh"
        else
            VOICE="Samantha"
        fi
    fi
fi
say -v "$VOICE" "$MESSAGE"
