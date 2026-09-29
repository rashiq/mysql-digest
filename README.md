# MySQL Digest demo

Build the demo from a separate source checkout:

```sh
make SOURCE=../mysql-digest
```

Use the release tag in that checkout before a deployment.
The build copies the Go runtime from the same toolchain as the WebAssembly compiler.
Update the asset version in `app.js` and `index.html` for each release.
