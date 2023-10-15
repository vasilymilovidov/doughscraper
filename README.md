## Doughscraper

Doughscraper is a CLI tool for preparing samples for [`tidalcycles/strudel`](https://github.com/tidalcycles/strudel). 

![doughscraper](https://github.com/vasilymilovidov/doughscraper/blob/main/doughscraper1.png?raw=true)

### Usage
Download the latest release for your OS from the release page and put `doughscraper` and `pitchdetector` binaries in your PATH.

#### Modes:
- **Rename pitched files** — 
 identifies files with note/octave pairs in their names and renames them by removing everything else. It requires the path to the local folder containing the pitched samples, for example, `Users/username/samples/piano`.
- **Generate JSON** — takes a folder containing samples, identifies subfolders with pitched and oneshot samples based on their filenames. It then generates a JSON file that can be used with strudel. You need to provide the path to the root of the local samples folder (e.g., `'Users/username/samples'`) and the path to a remote folder where you will upload your samples (e.g., `'https://raw.githubusercontent.com/vasilymilovidov/samples/main'`).
- **Detect pitch and rename** — identifies the pitch of the samples and renames them accordingly. It requires the path to the local folder containing the pitched samples, for example, ``Users/username/samples/piano``. Note that this script keeps the original filenames and adds the note/octave pairs to them. So, if you want to properly add these samples to the strudel.json file, run the **Rename pitched files** script first.

Do not include the closing backslash in the paths.

### Build
To build, you need to have Go and Rust installed.

Clone the repo and run:
```go
make all
```
Then copy both generated binaries (`doughscraper` and `pitchdetector`) to your PATH.

