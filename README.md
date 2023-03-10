# syncthing-map

A utility to map syncthing devices and shared folders.

The binaries are automatically released at each update of the stable and development (has `dev` in the version name) version. Use the stable one except if you want to test some development updates.

## Usage

`syncthing-map` relies on the configuration file (`config.xml`) of a device running Syncthing to build a map of relations between folders and devices that they are shared with.

There are currently two ways to provide these configuration files:

- **Manual**: you need to retrieve from the device you want to map its `config.xml` file manually. It leans copying it via `scp`, a Windows share or whatever you use to move files between your devices. Once you have this file locally, you run the appropriate command to add it to the "database of processed devices" (see the "Manual" section for details). You need to repeat this operation for each device (and its `confix.xml` file)

- **Automated**: since retrieving configuration failes each tile something changes is burdensome, you can automate the process by
  - sharing with Syncthing the configuration folder of a device (to one chosen device)
  - updating a file (`syncthing-map-server.yaml`) on that chosen device with that information
  - starting on that one chosen device `syncthing-map` in a server mode: this will start a web server you can connect to to see your map  
  The obvious advantage is that whenever you change something in any of the devices, it is automatically reflected in the map you get via the web server

It is probably a good idea to map one or two devices manually first to get an idea of how the map looks like (and if you like it 🤨) before jumping into the Automated mode.

### Automated (way cooler, requires to share the configurations)

This mode allows to request via HTTP an map of the connections, created on the fly.

#### Preparation of source folders

In order to process several folders, `syncthing-map` reads a configuration file: `syncthing-map-server.yaml`. It is an array of objects, which keys are `folder` (the folder that hold an `config.xml`) and `device` (the name of the device for the folder). An example of such a file on a Windows system:

```yaml
- device: my-laptop
  folder: C:\Users\Y\AppData\Local\Syncthing
- device: srv
  folder: D:\syncthing-configuration\srv
- device: galaxy-s22
  folder: D:\syncthing-configuration\galaxy-22
```

Make sure that `device` has the same as the one defined on the device which configuration is shared (usually its hostname). Otherwise you will see a rather chaotic map.

The natural way to share the configurations is though Syncthing itself! You will typically choose a machine on which you will run `syncthing-map` in server mode and share with it configurations from other devices.

In the example above, the machine `my-laptop` reads its own configuration (this is a Windows 10) and receives from `srv` and `galaxy-s22` configurations it stores in `D:\syncthing-configuration`. Cool, right?

Out of abundance of caution :) it is recommended to share the configuration in a "send only" mode. Just in case.

It is also recommended to use the `.stignore` file to ignore some ever-chnaging files (database, logs, ...). I use

```text
syncthing*
index*
index**
```

##### Where to find the config files

Most systems configurations are mentioned in the [documentation](https://docs.syncthing.net/users/config.html).

For Android you must manually input `/data/data/com.nutomic.syncthingandroid/files` as the folder.

#### Running the server

```text
syncthing-map server
```

If you launch [`http://localhost:3000`](http://localhost:3000) in a browser, you will see the map using the data from `syncthing-map-server.yaml`. Yay!

You can edit `syncthing-map-server.yaml` on the fly - each call to the server reads this config file.

Finally, each call also creates `syncthing-map-server.html`, an offline version of the map (this is what is actually sent back by the server).

### Manual (more traditional, where you retrieve the configurations yourself)

Run the following command repetedely for each `config.xml` you have access to ([how to find it](https://docs.syncthing.net/users/config.html))

```text
syncthing-map add --device <name of the device you took the config.xml from> --file <copied config.xml, possibly renamed>
```

An example of what you should see (with two devices/configs) is

```text
PS D:\syncthing-map> .\syncthing-map.exe add --device srv --file config-srv.xml
2023-01-09T19:46:02+01:00 INF wrote data-cli.json file
PS D:\syncthing-map> .\syncthing-map.exe add --device router --file config-router.xml
2023-01-09T19:46:16+01:00 INF wrote data-cli.json file
```

This added (or updated) two devices to the database file. This file (`data-cli.json` by default) will be initially created if absent, then updated with each `add` command. It gathers relevant information about each of the devices and its folders. The more you run `syncthing-map add`, the better your map will be - otherwise you will see that you are sharing fildes with a crazily named thing (this is the ID of the remote device).

When you are done with adding devices/configurations, run

```text
syncthing-map graph
```

This will create `syncthing-map.html` that you can open with a browser. If everything went well, you should see a comprehensive map of your devices and their folders.

![example of a map](example-1.png)

In the example above, you see a yellow rectange titled `router`. This a device. Its name comes from the `--device` parameter avove and its contents were generated based on the `config.xml` from that device.

It has two folders named `/etc` and `ssh keys`. They correspond to actual folders on your filesystem but the path is not shown here.

The way to read an arrow going from a folder is

> (folder) is connected as (`sendreceive`, `sendonly`, `receiveonly`) to (device)

In the case of the folder called `/etc`, it is shared in a `sendonly` mode with the device (redacted), and in the mode `sendreceive` with the device `srv`.

This is what you will see when sharing folders with a device that has not be "read" (its `config.xml` file was not added)

![example of a missing device](example-2.png)

This means that `Michael Cours` shared by a device (the name is not visible) is connected to a device that has not been processed (only known by its ID). This device is however defined in other `config.xml` files as `Michael Laptop` so the probable name is provided.

It may be that on that this unprocessed device has given itself another name. We will not know until its `config.xml` file is processed.

### Clean up

`syncthing-map clean` will remove the `.json` and `.html` files so that you can start from scratch.

## What next?

- better graph chart, not static, that would allow to maybe filter hosts and move them around.
- unlikely: I will write the above myself based on d3 or something like that

## Conclusion

This is a hobby project done follwong some [discussions on the Synthing forum](https://forum.syncthing.net/t/how-to-graph-my-clients/19554).

Feel free to open Issues if you find bugs, or start Discussions.

I should probbaly add a license but I do not care, so let it be [WTFPL](https://en.wikipedia.org/wiki/WTFPL).

If you have the irrestible need to share your gratitude, call someone you love or send money to a clever charity that helps with education.
