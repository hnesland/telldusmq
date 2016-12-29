# MQTT Broker for Telldus Core (Tellstick)
Uses libtelldus-core (Tellstick) to transmit events to MQTT. These events can then be used from different home automation / smart home utilities like [Home Assistant](https://home-assistant.io/) and [OpenHAB](http://www.openhab.org/) or your own custom handlers and loggers.

# Installation
The broker is installed with standard Go tools:
```
go get github.com/hnesland/telldusmq
```

The broker depends on libtelldus-core and libtelldus-core-dev-packges. Please see the [installation instructions](http://developer.telldus.com/wiki/TellStickInstallationUbuntu) for Telldus Core. 

Other dependencies are handled by Go. 

It's currently known to work on Ubuntu 16.04, but should also work on other Linux distributions. The C-based API for Telldus works on Linux, MacOS and Windows. This broker should also work on MacOS, but it is untested. If you're adventurous enough, maybe Windows even works ;-)

# Usage

After installation the broker can be started directly and will look for a configuration file. A sample configuration file is provided in the repository. The broker currently does not accept any parameters.

# Reporting Bugs

Please rebort bugs by raising issues for this project in [Github](https://github.com/hnesland/telldusmq/issues).

