# CLI Music app that can stream from OneDrive
## Problem
I have ~100GB audio library stored on my personal OneDrive account. Using [Microsoft Groove](https://en.wikipedia.org/wiki/Groove_Music) on Windows and [Cloudplayer](https://play.google.com/store/apps/details?id=com.doubleTwist.cloudPlayerPro&hl=en_AU&pli=1) on Android, I can stream that audio from OneDrive directly. There does not appear to be an off the shelf solution for Linux.

## Goal
Build a CLI media player to stream audio files from OneDrive.

## Core Features
- Select OneDrive folders to scan for music files
- Display a list of albums
- Stream one/multiple tracks from a chosen album

## Future Expansion
- Save local copies of albums:
  - Choose to download specific albums
  - Option to store any streamed file
- Add local folders to library, independent of OneDrive

## Development Plan
### Playing Audio
- Utiilise [Beep](https://github.com/gopxl/beep)

### Interacting with OneDrive
- Utilise [Microsoft Graph](https://developer.microsoft.com/en-us/graph/graph-explorer)
- Need to authenticate with OneDrive
- Need to display OneDrive folder structure
- Need to select folder(s) as root folders to scan for music
- Need to build a library from scanned OneDrive folders
- Need to list albums in library
- Need to list tracks in album
- Need to stream track 