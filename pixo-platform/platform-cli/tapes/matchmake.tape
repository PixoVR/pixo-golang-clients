
Output matchmake.gif

Require echo

Set Shell "bash"
Set FontSize 32
Set Width 1200
Set Height 800

Type "pixo mp servers deploy --pre-check" Enter Sleep 3s
Type "/A360" Sleep 1s Enter Sleep 500ms
Type "0.00.04" Sleep 1s Enter Sleep 1s

Type "pixo mp servers deploy" Enter Sleep 3s
Type "/A360" Sleep 1s Enter Sleep 500ms
Type "0.00.04" Sleep 1s Enter Sleep 500ms
Type "gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:0.0.6" Sleep 500ms Enter Sleep 2s

Type "pixo mp servers deploy --pre-check" Enter Sleep 3s
Type "/A360" Sleep 1s Enter Sleep 500ms
Type "0.00.04" Sleep 1s Enter Sleep 1s

Type "pixo mp matchmake --connect" Enter Sleep 3s
Type "/A360" Sleep 1s Enter Sleep 500ms
Type "0.00.04" Sleep 1s Enter

Set PlaybackSpeed 20
Sleep 15s
Set PlaybackSpeed 1

Type "ACK" Sleep 500ms Enter Sleep 500ms
Type "PLAYER_CONNECT walker" Sleep 500ms Enter Sleep 500ms
Type "PLAYER_COUNT" Sleep 500ms Enter Sleep 500ms
Type "PLAYER_DISCONNECT walker" Sleep 500ms Enter Sleep 500ms
Type "PLAYER_COUNT" Sleep 500ms Enter Sleep 500ms
Type "exit" Sleep 500ms Enter

Sleep 5s
