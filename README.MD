# DisGo
DisGo is CLI music player for streaming music from OneDrive. 

![demo](demo.gif)
(Gif created using [VHS](https://github.com/charmbracelet/vhs))


## Installation
### Binaries
- Binaries are provided for Linux, Windows, and MacOS (amd64 and arm64), and should execute without dependencies
   - The MacOS binaries are presently untested
### Source
- Requires [Go to be installed](https://go.dev/doc/install)
- Make sure your [Go install directory is in your PATH]()
- To install, either:
   - Download the source and run ``go install .`` from inside the source, or
   - Run ``go install github.com/jthughes/disgo``

## Usage
1. On first time run, ``.disgo/`` will be created in your home directory to store your library metadata and other configuration data. The browser will also open and prompt to authenticate with OneDrive. Sign in with OneDrive account where your music is stored and accept the permisions.
2. Use the ``scan <folder>`` command with a valid folder to find all the music albums in that folder. This should be a folder path relative to your OneDrive root. So if in the root of your OneDrive, you have a ``Music/`` folder, and you want to scan the ``Video Games/`` folder inside, you would run the command ``scan /Music/Video Games/``. This creates a sqlite database of metadata called ``library.db``.
3. Run ``ls`` to list all albums found in ``library.db``.
4. Run ``play <#>`` where # is the number listed against an album as a result of running ``ls``.
5. Run ``help`` to learn about other commands.

## About the project
I have ~100GB audio library stored on my personal OneDrive account. Using [Microsoft Groove](https://en.wikipedia.org/wiki/Groove_Music) on Windows and [Cloudplayer](https://play.google.com/store/apps/details?id=com.doubleTwist.cloudPlayerPro&hl=en_AU&pli=1) on Android, I can stream that audio from OneDrive directly. With no equivalent option for Linux, this is my solution.

### Future Expansion
- Build playlists
- Shuffle playlists
- Add file sources:
  - Dropbox
  - Google Drive
  - Local files
  - Local cache of cloud files
- DBUS/MPRIS integration (interact with media keys/bluetooth headphone controls)
- UI with [BubbleTea](https://github.com/charmbracelet/bubbletea)
  - Message history
  - Visual feedback for current playing track (sliders, etc)
- Support for additional file types

## Built using
- [Beep](https://github.com/gopxl/beep)
- SQLite