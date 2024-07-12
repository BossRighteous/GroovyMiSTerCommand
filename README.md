# GroovyMiSTerHTTP
An HTTP client to allow remote GroovyMAME execution from MiSTer

Actually it's going to be done over UDP but I'm not going to rename/readdress everything yet.

## Alpha Setup Guide
Assumes Go is installed, and I am using WSL2 on windows so my pathing and UDP is wonky but works.

```
git clone git@github.com:BossRighteous/GroovyMiSTerHTTP.git
cd GroovyMiSTerHTTP
cp ./spec/GmcServerConfig.json ./config.json
// Edit JSON for mame bin paths and MiSTer IP
go run cmd/groovymisterhttp/main.go "/path/to/GroovyMiSTerHTTP/config.json"
```

In new terminal screen, send UDP commands to server port
```
echo -n '{"cmd":"mame","vars":{"MACHINE_NAME":"circus"}}' >/dev/udp/localhost/32105

echo -n '{"cmd":"unload"}' >/dev/udp/localhost/32105
```

### Groovy MiSTer Command (GMC) spec

The Groovy MiSTer core does not currently support file browsing as there are no direct roms to load.
The following data type is a proposed file structure the Groovy core can open as a rom-like and perform
UDP RPC commands against.

```
Filename:
{Human ReadableText Title}.gmc
Contents: UTF-8 JSON Text
{
    "cmd": "mame",
    "vars": {
        "MACHINE_NAME": "circus"
    }
}
```

When the core opens the GMC file, it's contents will be streamed over UDP to the server.
It is assumed this will be plain text, but the core will send the binary data as it. It is up
to the server software to read and handle appropriately.

This configuration is managed from the server bin via a settings JSON
```
{
    "mister_host": "192.168.1.255",
    "commands": [
        {
            "cmd": "mame",
            "work_dir": "/path/to/groovymame",
            "exec_bin": "./mame",
            "exec_args": ["${MACHINE_NAME}"]
        }
    ]
}
```

This allows the server to easily control it's command whitelist and intented variable replacements in args.
Since the command will be executed outside of a shell context ENV and command line variable substitution isn't possible.
Arg/Var replacement is done by the server. This also makes it much more difficult to inject escape sequences.

```
/path/to/mybin positional -tf --foo="complex"

... might become
"cmd": "mycmd"
"work_dir": "/path/to/",
"exec_bin": "./mybin",
"exec_args": ["positional", "-tf", "--foo=\"${FOO_VAL}\""]

... and be executed by
{"cmd":"mycmd", "vars":{"FOO_VAL":"custom"}}
```

GroovyMiSTer will know of the server by way of init beacon. Upon server boot, and every n seconds therafter,
an init handshake will be sent to the Groovy core over UDP until the Core acknoledges.

This enables the MiSTer to reply even in cases the server machine amy have a dynamic IP.
Your GroovyMister enabled emulator will need to have it's own configuration including host/IP.

The server will attempt to run asyncronous (long running processes) until exit or error.
Any running processes will be terminated on receipt of the next GMC command over UDP.

Messages including a stderr tail will be displayed on the MiSTer. This is under development for messaging data sources.

#### Groovy MiSTer UX

psakhis has conceptually reviewed the following and may go forward with something like the following

Similar to many other cores a load rom menu setting may be implemented:
```
Load Remote Command *.gmc
```

This will open a file browser interface, pointing to rom-like directory
```
/media/fat/games/Groovy/**.gmc
```

Upon GMC file selection the core may
- load binary contents up to 1024 bytes
- Send bytes as UDP packet to acked Server socket
- Server will manage process state and errors via text blitting back to the MiSTer

#### File Browsing


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

I've noted that .zip files are browseable as filesystem, so a single portable zip may reduce FS block fill and make the build/transfer simpler. I don't know the limits to file density in a zip but could test this.

This would allow the server binary to run a CLI command routine that
- Asks for category/filter preferences
- Checks local MAME rom directory
- Generates temp GroovyMAME/ director
- Adds Utils and All Roms found locally based on filters
- Dupes into category hierarchies selected.
- zips directory

This could then be trivially SCPed and extracted.