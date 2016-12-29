# MQTT Broker for Telldus Core (Tellstick)
Uses libtelldus-core (Tellstick) to transmit events to MQTT. These events can then be used from different home automation / smart home utilities like [Home Assistant](https://home-assistant.io/) and [OpenHAB](http://www.openhab.org/) or your own custom handlers and loggers.

# Installation
The broker is installed with standard [Go](https://golang.org) tools:
```
go get github.com/hnesland/telldusmq
```

The broker depends on libtelldus-core and libtelldus-core-dev-packages. Please see the [installation instructions](http://developer.telldus.com/wiki/TellStickInstallationUbuntu) for Telldus Core.

Other dependencies are handled by Go.

It's currently known to work on Ubuntu 16.04, but should also work on other Linux distributions. The C-based API for Telldus works on Linux, MacOS and Windows. This broker should also work on MacOS, but it is untested. If you're adventurous enough, maybe Windows even works ;-)

# Usage

After installation the broker can be started directly and will look for a configuration file. A sample configuration file is provided in the repository. The broker currently does not accept any parameters.

# Configuration

The broker is configured using a [Viper](https://github.com/spf13/viper) compatible configuration file (JSON, TOML, YAML, HCL or Java properties). Currently the file must be named `telldusmq.<type>` like `telldusmq.json`. We'll look for the configuration file in `/etc/telldusmq/`, `$HOME/.telldusmq/` or the current working directory.

The MQTT parameters for topic and payload, can contain template variables (from Go text/template package). The templates are executed on an object containing the event from Tellstick with the following properties:

  - Class
  - Protocol
  - Model
  - Code
  - House
  - Unit
  - Group
  - Method
  - Id
  - Temp
  - Humidity
  - Value
  - DataType

We can then use the variables like this: `Temperature for id#{{.Id}} is: {{.Temp}}`

Some sensors supports both temperature and humidity and if you wish to have a separate event for both these types, you can set the configuration item `Tellstick => SplitTemperatureAndHumidity` to `true`. The `Value` property will contain the value (temperature or humidity) and the `DataType` property will contain `temp` or `humidity`. The sample configuration reflects this.

Telldus Core reports device methods as `turnon` or `turnoff` by default. If you wish to map these to other values like `ON` or `OFF`, `1` or `0` etc, this can be configured with setting `Tellstick => MapTurnOnTo` and `Tellstick => MapTurnOffTo` to the appropriate strings. 

# Reporting Bugs

Please rebort bugs by raising issues for this project in [Github](https://github.com/hnesland/telldusmq/issues).
