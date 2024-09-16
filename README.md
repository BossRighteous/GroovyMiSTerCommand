# GroovyMiSTerCommand
A UDP Server to allow local Emulator execution from the remote (Groovy_MiSTer)[https://github.com/psakhis/Groovy_MiSTer] Core.

Very sincere thank you to (@psakhis)[https://github.com/psakhis].

I check the following discords somewhat frequently but not daily. \
GroovyArcade > nogpu
MiSTerFPGA > mister-cores > Groovy MiSTer

## Installation

This may be run via normal go runtime commands.

The releases section has builds for major platforms, let me know if you need a one-off.
Go is very friendly to building cross-platform from source.

- Download the latest release binary for your platform
- Save to new folder (anywhere is fine)
- Download/copy config.dist.json to the same directory as the binary.
- Update the IP with your MiSTer's static IP or current IP
- Update all applicable Paths for your system type. Example is Windows.
- If using MAME generator, download mamelist.arcade-no-clone.266.json and rename as mamelist.json

By default, executing with no command line args will run the server instance.

There are also cli args that can be used to run GMC generator scripts
- `GroovyMiSterCommand.exe server` - Runs server (same as no command)
- `GroovyMiSterCommand.exe generate:retroarch` - Scans config playlists_dir for your RA playlist files, generating GMCs from each as a folder
- `GroovyMiSterCommand.exe generate:mame` - Scans mamelist.json and cross references against your defined roms_dir for existence. This is opinionated and tries a bit of folder organization too.
- `GroovyMiSterCommand.exe generate:directory {name}` - Uses a generic directory walker to generate GMCs by filetype. Limited to depth of 4. Can be configured for any emulator with standard filename based args. Used with mednafen as example. `"extensions": [],` will match any extension.

### Groovy MiSTer Command (GMC) spec

The GroovyMiSTer core has been extended to allow loading of GMC files as if they were ROMs.

When this server handshakes with the Groovy core, the connection will allow the loading of the GMC files to broadcast to the server.

The GMC file defines a whitelist of command execution strategies, as defined in the config.json

{Human ReadableText Title}.gmc
```
{
    "cmd": "mame",
    "vars": {
        "MACHINE_NAME": "circus"
    }
}
```

config.json
```
{
    ...
    "commands": [
        {
            "cmd": "mame",
            "work_dir": "/path/to/groovymame",
            "exec_bin": "mame",
            "exec_args": ["${MACHINE_NAME}"]
        }
        ...
    ]
    ...
}
```

There is a fixed relationship between the config's command execution instruction and what is required by the GMC file to fulfil the command request. The `cmd` property is how the GMC is routed to a command handler. The `vars` object values are effectively replaced into the all the `exec_args` as needed.

GroovyMiSTer will know of the server by way of init beacon. Upon server boot, and every n seconds therafter,
an init handshake will be sent to the Groovy core over UDP until the Core acknoledges.

This enables the MiSTer to reply even in cases the server machine amy have a dynamic IP.
Your GroovyMister enabled emulator will need to have it's own configuration including host/IP.

The server will attempt to run asyncronous (long running processes) until exit or error.
Any running processes will be terminated on receipt of the next GMC command over UDP.

Messages including a stderr tail will be displayed on the MiSTer. This is under development and honestly pretty poor at the moment.

### Generators
Because the GMC spec is easily templated it's somewhat trivial to create new parsers. The `directories` generator config is a good example.
Generators will auto-create new drectories as needed in the application dir/Groovy.

These files can be SCPed or otherwise transfered to your `/path/to/games/Groovy` directory like any other Core's ROMs.

#### Directory Generator
```
{
    "name": "mednafen",
    "dir": "C:\\Users\\bossr\\Downloads\\mednafen-1.32.1-win64\\roms",
    "extensions": ["zip", "chd"],
    "template": {
        "cmd": "mednafen",
        "vars": {
            "ROM_PATH": "${ROM_FULL_PATH}"
        }
    }
}
```
- `name` acts as the CLI parameter reference as well as the folder name in resultant GMC files
- `dir` indicates what directory to scan for files. Each being looped and created as a GMC file based on the template
- `extensions` filters files by matching extensions. `[]` as a value allows all when you don't know specific file extensions, but this can also produce unplayable GMCs pointing to non-ROM files. Limited to 4 levels of recursion.
- `template` This template is the schema of the resulting GMC.

Special variables will be replaced into the `vars` object values based on the file being looped where defined.
- `${ROM_FULL_PATH}` - Full path including absolute directory structure and filename
- `${ROM_FULL_DIR}` - Full path to containing directory of file
- `${ROM_RELATIVE_PATH}` - directory path and filename relative to the `dir` defined in the generator
- `${ROM_RELATIVE_PATH}` - directory relative path to the `dir` defined in the generator

#### RetroArch Generator
```
"retroarch": {
    "playlists_dir": "C:\\RetroArch-Win64\\playlists"
},
```
RetroArch is actually pretty convenient to generate GMCs for, because it does it's own indexing in the form of Playlist LPL files.

Simply point the config to your playlists_dir and a folder will be generated for each playlist.
Each playlist's entry will result in a single GMC file.

*Important Note* RetroArch via command line cannot infer the correct Core to use. Each Playlist will need a DEFAULT CORE defined in the GUI or the .lpl file itself.
```
"default_core_path": "C:\\RetroArch-Win64\\cores\\mupen64plus_next_libretro.dll",
```
If this is not present the playlist will be skipped.

#### MAME Generator
```
"mame": {
    "roms_dir": "C:\\Users\\bossr\\groovymame\\groovymame_0266.221d_win-7-8-10\\roms",
    "mamelist_path": "mamelist.json"
},
```

MAME romsets are both massive, and completely full of related/dependent rom files. Without digging in too deep I sought to provide the following:
- Ability to categorize based on Genre, Manufacturer, Year
- Ability to generate GMCs for playable roms where present
- Ability to further allow filtering of wanted/unwanted roms.

I took an opinionated approach and first tried to source a list of playable roms parent/clone. This is viewable in `mamelist.all-playable.json`. I then used JQ to filter the set to my preferences for arcade/no-clone sets.

Some helpful JQ
```
// Unique genres list
map(.genre) | unique

// Full -> Arcade
map(select(.genre | startswith("Arcade") or startswith("Ball") or startswith("Climbing") or startswith("Driving") or startswith("Fighter") or startswith("Maze") or startswith("MultiGame") or startswith("Platform") or startswith("Puzzle") or startswith("Shooter") or startswith("Sports")))

// Arcade -> no-clone
map(select(.cloneof == "" and (.manufacturer | startswith("bootleg") | not)))
```

This gives us a bunch of these
```
{
    "name": "nbajam",
    "description": "NBA Jam (rev 3.01 4-07-93)",
    "cloneof": "",
    "manufacturer": "Midway",
    "year": "1993",
    "genre": "Sports - Basketball"
}
```

Which allow us to create GMCs in categorical folders, referencing the machine name for the GMC Var, and the description as the filename.

Your personal filtered list as `mamelist.json` will be cross-referenced against your MAME rom dir for existence as well. This allows people to generate the whole set, or use a broader mamelist.json even with only a few local ROMs.