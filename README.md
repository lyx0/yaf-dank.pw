# yaf - Yet Another (Temporary) Fileshare
yaf is a simple Go program to handle file uploads.
If you also want to serve the uploaded files, consider a web server like [nginx](https://nginx.org/en/).

## Installation with docker-compose and local build
**Clone** the directory:
```bash
git clone https://github.com/lyx0/yaf.git
```
Run **tests** (optional):
```bash
go test
```

## Manual Installation
**Clone** the directory:
```bash
git clone https://github.com/lyx0/yaf.git
```
**Build** the executable:
```bash
go build
```
Run **tests** (optional):
```bash
go test
```

If you plan on using a systemd service or another init system, you might want to move the `yaf` executable to a different directory (e.g. `/opt`) at this point; you know your setup best.

## Configuration
### yaf
There are just a few parameters that need to be configured for yaf.
Refer to the `example.conf` file:
```
Port:       4711
# a comment
LinkPrefix: https://yaf.example.com/
FileDir:    /var/www/yaf/
LinkLength: 5
ScrubExif: true
# Both IDs also refer to the "Orientation" tag, included for illustrative purposes only
ExifAllowedIds: 0x0112 274
ExifAllowedPaths: IFD/Orientation
ExifAbortOnError: true
FileExpiration: false
```

Option             | Use
------------------ | -------------------------------------------------------------------
`Port`             | the port number yaf will listen on
`LinkPrefix`       | a string that will be prepended to the file name generated by yaf
`FileDir`          | path to the directory yaf will save uploaded files in. if using docker-compose needs to be the same as the target mount point (the right side)
`LinkLength`       | the number of characters the generated file name is allowed to have
`ScrubExif`        | whether to remove EXIF tags from uploaded JPEG and PNG images (`true` or `false`)
`ExifAllowedIds`   | a space-separated list of EXIF tag IDs that should be preserved through EXIF scrubbing (only relevant if `ScrubExif` is `true`)
`ExifAllowedPaths` | a space-separated list of EXIF tag paths that should be preserved through EXIF scrubbing (only relevant if `ScrubExif` is `true`)
`ExifAbortOnError` | whether to abort JPEG and PNG uploads if an error occurs during EXIF scrubbing (only relevant if `ScrubExif` is `true`)
`FileExpiration`   | whether to automatically remove files after a given time or not (`true` or `false`)


Make sure the user running yaf has suitable permissions to read, and write to, `FileDir`.
Also note that `LinkLength` directly relates to the number of files that can be saved.
Since yaf only uses alphanumeric characters for file name generation, a maximum of `(26 + 26 + 10)^LinkLength` names can be generated.

#### A Note on EXIF Scrubbing
EXIF scrubbing can be enabled via the `ScrubExif` config key.
When enabled, all standard EXIF tags are removed on uploaded JPEG and PNG images per default.
It is meant as a last-line "defense mechanism" against leaking PII, such as GPS information on pictures.
**If possible, you should always prefer disabling capturing potentially sensitive EXIF tags when creating the images!**

Obviously, EXIF tags serve a purpose and you may want to keep _some_ of the information, e.g., image orientation.
The `ExifAllowedIds` and `ExifAllowedPaths` config keys can be used to selectively allow specific tags to survive the scrubbing.
The IDs for standard tags can be found in [1].
You may specify tag IDs in decimal and hexadecimal notation.
(In the latter case, the ID _must_ start with `0x`.)

The path specification for `ExifAllowedPaths` relies on the format implemented in [`go-exif`](https://github.com/dsoprea/go-exif) which is "documented" in machine-readable format in [2].
Multiple paths can be specified, separated by a space.
The path format is as follows:

1. For tags in the main section: `IFD/<GroupName>/<FieldName>`.
   Examples: `IFD/Orientation`, `IFD/Exif/Flash`, `IFD/GPSInfo/GPSTimeStamp`.
   You will probably want to use both [1] and [2] in combination if you plan to specify allowed tags by path.

2. Tags in the thumbnail section follow the same format but paths start with `IFD1/` instead of `IFD`.

### nginx
If you use a reverse-proxy to forward requests to yaf, make sure to correctly forward the original request headers.
For nginx, this is achieved via the `proxy_pass_request_headers on;` option.

If you want to limit access to yaf (e.g. require basic authentication), you will also need to do this via your reverse-proxy.

### caddy
I provided a `Caddyfile.example` for you that should be pretty self explanatory. Copy the contents to your own `Caddyfile` and be sure to move the contents of the `dist` folder to your file directory so you can enjoy the really pretty high quality frontend page.
```Caddyfile
yaf.example.com {
    root * /path/to/filedir/
    file_server

    reverse_proxy /upload localhost:4711
    reverse_proxy /uploadweb localhost:4711
}
```

## Running

### Manually
After adjusting the configuration file to your needs, run:
```bash
yaf -configFile yaf.conf
```
Of course, you can also write a init system script to handle this for you.

### Running from docker-compose
**Copy** configuration file and fill it in:
```bash
cp example.conf yaf.conf
```
**Configure** the `docker-compose.yml` volume paths:
```bash
vim docker-compose.yml
```
**Build** the local docker file:
```bash
make build
```
**Run** the local docker file with docker-compose:
```bash
make run
```

### Running from Docker
Building the Docker image and running it locally
```bash
docker build -t yaf .
docker run \
    -p 4711:4711 \
    -v /path/to/your/yaf.conf:/app/yaf.conf \
    -v /path/to/local/filedir:/var/www/yaf \
    yaf
```

Port 4711 is the default port for the server in `example.conf`, if you've changed this in your config you'll need to change this in the `docker run` invocations above too.  
The above runs forwards the yaf port from 4711 in the container to 4711 on your local system.

## Usage
You can use yaf with any application that can send POST requests (e.g. ShareX/ShareNix or just `curl`).
Make sure the file you want to upload is attached as a `multipart/form-data` field named `file`.
In `curl`, a request to upload the file `/home/alice/foo.txt` could look like this:
```bash
curl -L -F "file=@/home/alice/foo.txt" yaf.example.com/upload
```
The response will include a link to the newly uploaded content.
Note that you may have to add additional header fields to the request, e.g. if you have basic authentication enabled.

## Inspiration
- [i](https://github.com/fourtf/i) by [fourtf](https://github.com/fourtf) – a project very similar in scope and size
- [filehost](https://github.com/nuuls/filehost) by [nuuls](https://github.com/nuuls) – a more integrated, fully-fledged solution that offers a web interface and also serves the files


[1]: https://exiv2.org/tags.html
[2]: https://github.com/dsoprea/go-exif/blob/a6301f85c82b0de82ceb8501f3c4a73ea7df01c2/assets/tags.yaml