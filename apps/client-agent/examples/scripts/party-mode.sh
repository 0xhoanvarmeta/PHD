#!/bin/bash
# Party Mode - Open multiple apps/websites at once

MODE="${1:-default}"

echo "🎉 PARTY MODE ACTIVATED! 🎉"

case "$MODE" in
    birthday)
        echo "Birthday Party Mode..."
        # Open birthday websites
        open "https://www.youtube.com/watch?v=_z-1fTlSDF0" # Happy Birthday song
        open "https://birthdaycake.net"
        # Speak
        say "Happy Birthday! Let's celebrate!"
        # Notification
        osascript -e 'display notification "🎂 Party time! 🎉" with title "Birthday Mode" sound name "Glass"'
        ;;

    work)
        echo "Work Mode..."
        open -a "Slack"
        open -a "Visual Studio Code"
        open -a "Google Chrome" "https://gmail.com"
        say "Time to work!"
        ;;

    chill)
        echo "Chill Mode..."
        open -a "Spotify"
        open "https://www.youtube.com/watch?v=jfKfPfyJRdk" # Lofi music
        say "Relax and enjoy"
        ;;

    prank)
        echo "Prank Mode... 😈"
        # Open random websites
        for i in {1..5}; do
            open "https://www.google.com/search?q=random$i"
            sleep 0.5
        done
        say "Surprise!"
        ;;

    *)
        echo "Default Party Mode..."
        # Open fun stuff
        open "https://www.youtube.com"
        say "Let's have some fun!"
        ;;
esac

echo "✨ Party Mode Complete! ✨"
