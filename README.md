# civkeyer

[![Tests](https://github.com/bbathe/civkeyer/workflows/Tests/badge.svg)](https://github.com/bbathe/civkeyer/actions) [![Release](https://github.com/bbathe/civkeyer/workflows/Release/badge.svg)](https://github.com/bbathe/civkeyer/actions)

civkeyer is a Windows application that interfaces with [Icom](https://www.icomamerica.com/en/amateur/) radios to send predefined CI-V commands.  The original purpose was to send CI-V commands to initiate the playback of recorded Tx's but there is no constraint for having it send whatever CI-V commands you want.

## Description

I'm using a [BlueMax49ers FTDI Icom CI-V Cat Control Programming Cable](https://smile.amazon.com/gp/product/B074JRWYRP) to connect to an [Icom IC-7300](https://icomamerica.com/en/products/amateur/hf/7300/default.aspx).  I'm sure other cables will work but this one is really high quality and the instructions and support are from [KJ6ZWL](https://www.qrz.com/db/KJ6ZWL).

## Installation

To install this application:

1. Create the folder `C:\Program Files\civkeyer`
2. Download the `civkeyer.exe.zip` file from the [latest release](https://github.com/bbathe/civkeyer/releases) and unzip it into that folder
3. Create the `civkeyer.yaml` file (a [YAML](https://en.wikipedia.org/wiki/YAML) file is a plain text file with the `.yaml` extension) in that folder, with these attributes:
    ```yaml
    connection:
      port: COM9
      baud: 9600
    functions:
      -
        label: F1 Callsign
        message: FEFE94E0280001FD
      -
        label: F2 Exchange
        message: FEFE94E0280002FD
      -
        label: F3 CQ
        message: FEFE94E0280003FD
      -
        label: F4 QRZ
        message: FEFE94E0280004FD
    ```
    - `connection` defines how the application connects to the radio.
      - `port` is the name of the Windows COM port to use to send the CI-V commands, this is setup when you install the device driver for the Cat Control Programming Cable.  You should be able to find this in Device Manager.
      - `baud` is the rate at which information is transferred to the COM port, this is a setting on the port that is setup when you install the device driver for the Cat Control Programming Cable.  You should be able to find this in Device Manager, check the "Port Settings" tab for the device.
    - `functions` defines the commands that can be sent to the radio.  You can include up to 12 functions.
      - `label` is the button label you want associated with the function.  Hotkeys will be associated with each function command, starting with F1 for the first function, F2 for the second, etc.  So, you will probably want to include that in the label and then whatever note to remind you what it will do.
      - `message` is the byte code sequence to send to the radio, encoded as hex.  
4. You can now double-click on the `civkeyer.exe` file to start the application.  Creating a shortcut somewhere will make it easier to find in the future.  You can press the button or Hotkey to execute the function.  The application has to have focus in order for the Hotkeys to work.

There will be a log file created in the same directory as the executable and any errors are logged there.

## References

[DF4OR: ICOM CI-V Information](http://www.plicht.de/ekki/civ)