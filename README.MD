# DisGo
DisGo is CLI music player for streaming music from OneDrive. 

## Usage
1. On first time run, the browser will open and prompt to authenticate with OneDrive. Sign in with OneDrive account where your music is stored and accept the permisions.
2. Use the ``scan <folder>`` command with a valid folder to find all the music albums in that folder. This should be a folder path relative to your OneDrive root. So if in the root of your OneDrive, you have a ``Music/`` folder, and you want to scan the ``Video Game/`` folder inside, you would run the command ``scan Music/Video\ Game/``. This creates a sqlite database of metadata called ``library.db``.
3. Run ``ls`` to list all albums found in ``library.db``.
4. Run ``play <#>`` where # is the number listed against an album as a result of running ``ls``.
5. Run ``help`` to learn about other commands.

## About the project
I have ~100GB audio library stored on my personal OneDrive account. Using [Microsoft Groove](https://en.wikipedia.org/wiki/Groove_Music) on Windows and [Cloudplayer](https://play.google.com/store/apps/details?id=com.doubleTwist.cloudPlayerPro&hl=en_AU&pli=1) on Android, I can stream that audio from OneDrive directly. With no equivalent option for Linux, this is my solution.

### Future Expansion
- Build playlists
- Shuffle playlists
- Save local copies of albums:
  - Choose to download specific albums
  - Option to store any streamed file
- Add local folders to library, independent of OneDrive
- Add other cloud storage providers
- DBUS/MPRIS integration
- UI with [BubbleTea](https://github.com/charmbracelet/bubbletea)
- Support for additional file types
- Other OS support?

## Built using
- [Beep](https://github.com/gopxl/beep)
- SQLite