# DDNS-SRV Plugin Builder

This application can be used to easily create [ddns-srv](https://github.com/pbergmam/ddns-srv) plugins from a package that implements the [libdns interface](https://github.com/libdns/libdns). 

It automates the steps needed to create a new provider plugin, including:

- Initializing a new [module](https://go.dev/doc/tutorial/create-module)  
- Adding dependencies  
- Creating a minimal `main` file  
- Building the plugin  

---

### Installing

```bash
go install github.com/pbergman/ddns-srv-plugin-builder@latest
````

---

### Project Template

You can also use this tool to create a project template for developing a plugin. This is useful for testing or setting up a new plugin without building it immediately.

```bash
ddns-srv-plugin-builder -build-dir ./build -no-build -no-cleanup github.com/libdns/example
```

This command will create and configure a new project in the `./build` directory without building the plugin or cleaning temporary files.

---

### Module Replace

The builder supports basic module replacement, for example when testing a fork:

```bash
ddns-srv-plugin-builder example.com/organisation/package=github.com/me/package@package-fork
```

This replaces `example.com/organisation/package` with `github.com/me/package` on the branch `package-fork`.

For more complex replacements, it is recommended to create a project template and manually modify the `go.mod` file before building the plugin.
