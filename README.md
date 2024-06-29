# GroovyMiSTerHTTP
An HTTP client to allow remote GroovyMAME execution from MiSTer

## Draft Data Contracts

### Groovy MiSTer

The Groovy MiSTer core does not currently support file browsing as there are no direct roms to load.
The following data type is a proposed file structure the Groovy core can open as a rom-like and perform
an HTTP RPC command against.

```
Filename:
{Text Title}.{GroovyMisterCommand}
Contents: UTF-8 Text
{UtilKeyword | MameSlug}

Examples:
// Game Slug
Hanabi de Doon! - Don-chan Puzzle.gmc
doncdoon

// Util Keywords
Test Connection.gmc
connect

Unload Rom.gmc
unload

Kill Server Process.gmc
kill

Reboot Host Machine.gmc
reboot

Shutdown Host Machine.gmc
shutdown
```

This is a minimal amount of information and should provide just enough data for the server to execute tasks without security concerns.

I initially considered a {command, args} schema but limiting args to single 'word' will make verification much easier.

#### Groovy MiSTer UX

Similar to many other cores a load rom menu setting may be implemented:
```
Load Remote GroovyMAME *.gmc
```

This will open a file browser interface, pointing to rom-like directory
```
/media/fat/games/GroovyMAME/**.gmc
```

Upon GMC file selection the core may
- load the contents into string/chars CMD_STR
- Check new ini/env key GROOVY_MAME_HOST (Example: "192.168.1.128")
- Execute syncronous HTTP POST
  - Address: "http://{GROOVY_MAME_HOST}:32105/cmd"
  - Body: "{CMD_STR}"
- For 200/ok response, close OSD, wait for game broadcast from PC server
- ? For non-200 response, display error message from body or http status code

This could trivially be a GET request with query parameters as well. Depends on what the Core's build would best support.

#### File Browsing
I am making some assumptions here that there is a file browsing interface module the core
can point to. There are probably opportunities to expand this or make it server-listing dependent.

For now I am assuming a 1:1 file per ROM load command. In its most basic form a directory tree may be used for
command presentation and display.

Because the MAME rom list is massive I am proposing using either symlink or pure file cloning to create additional hierarchies.
```
/media/fat/games/GroovyMAME/
    Utils/
        Unload Rom.gmc
        Kill Server Process.gmc
        Reboot Host Machine.gmc
        Shutdown Host Machine.gmc
    All Roms/
        [Full Titles Set].gmc
```

To keep the Core's command logic simplified, companion server-side scripts can be used to populate this original listing set as well as duped categorizational directories. It is unclear to me if a file browsing interface would support symlinks but <1kb file dupes even in the thousands aren't the end of the world.

```
/media/fat/games/GroovyMAME/
    Utils/...
    All Roms/...
    Genre/
        Action/
        ...
    Year
        1987/
        ...
    System
        CPS2/
        ...
```

I've noted that .zip files are browseable as filesystem, so a single portable zip may reduce FS block fill and make the build/transfer simpler.

This would allow the server binary to run a CLI command routine that
- Asks for category/filter preferences
- Checks local MAME rom directory
- Generates temp GroovyMAME/ director
- Adds Utils and All Roms found locally based on filters
- Dupes into category hierarchies selected.
- zips directory

This could then be trivially SCPed and extracted.

This could also be exposed as a special HTTP endpoint for a MiSTer Script to process against the remote

### Groovy MiSTer HTTP Server

This is a binary CLI app that will live on the GroovyMAME host machine.

The IP address or hostname of the server machine should be defined in the MiSTer ini for GroovyMiSTer settings.

The server binary itself should be placed in the GroovyMAME application directory for easier cross-platform build pathing. May include .ini configuration for the server's paths and ports to decouple.

On execution the server will listen generically to any host, but specifically on port 32105. This port allows for additional GroovyMiSTer UDP expansion but occupies similar range.

The server listens for commands on the ` POST /cmd` endpoint. It expects basic single word text data as defined and generated for the .gmc spec.
Based on the command, a single subprocess may be terminated or replaced.

For example: a command will
- Compare POST body command to known UtilKeywords
- Execute Utility routine on match
- Else pattern match for MAME game slug
    - Check GroovyMAME executable for game slug verification
        - Fail if not verifiable
    - Check for existing GroovyMAME subprocess, exit if operational
    - Attempt to create new GroovyMAME subprocess based on game slug
    - Runs as goroutine/async
- Else fail
- Return HTTP response

The CLI app must be executed/ran on the Server before the GroovyMiSTer core can access it.

For a delegated mini-pc this could be run on startup. The server process is extremely light weight and can safely run in the background without significant resource drain.

PC, OSX, and Linux builds are possible from the single go repository.