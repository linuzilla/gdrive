# Google Drive Utility

A simple and useless command line utility for Google Drive.

### Descriptions

There are whole brunch of utilities which can be found on the net which can upload, download, or synchronize files
between local to google drive. This is an experimental project did not do much about being a real utility.

This utility use "service account" as a way to access Google Drive. However, a service account could not use account
owner's drive directly, that means, you have to allow your drive share with your service account, or service account's
drive share with yourself. A better way to share between you and your service account is using "shared drive."  As you
may expect, share your "shared drive" with your service account.

If you have no idea about "service account", try to get one from google api console.

### Build

Simply write a simple main program like this.

```go
package main

import "github.com/linuzilla/gdrive"

func main() {
	gdrive.Start()
}
```

### Usage

First, you need a credential of your service account, go to google api console to get one, and the download the
credential in "json" format.

prepare your config json file, something like

```yaml
application:
  name: Google Drive Sync

google-drive:
  credential: /path/to/your/credentials.json

database:
  file: /path/to/database/file/without/extension

plugin:
  commands: /path/to/plugins/directory
```
