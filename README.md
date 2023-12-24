# QR Tool

## Use Cases
This tool was created to give you the power to generate and/or read QR Codes which are stored locally in your machine.

## Usage for Supplied Executables

### Option 1 - Double Click
You can double-click the executable file of this application and it will open a 
CMD/Terminal window which can then be used normally in-line with the further
instructions which are displayed during runtime.

*-- OR --*

### Option 2 - Run from Terminal/CMD
The other option to run this application is to open a Terminal/CMD session and
execute it from there. This is different between Windows and Linux/MacOS.

#### Windows
<hr>
Open a CMD Session in any way you want, then

```dos
cd Downloads
qrTool_{arch}.exe
```

<b><i>REMEMBER: </i></b>
Replace `{arch}` with the software architecture on the filename.

#### Linux/MacOS
<hr>
Open a Terminal and type :

```sh
cd Downloads
./qrTool_{arch}_{OS}
```
`{arch}` is the processor architecture.
`{OS}` is the Operating System you are running the app on.

## Compiling from Source
To compile the source code you would need to have installed the go language on your system, which is available from [go.dev](https://go.dev/).

Supposing you have installed correctly the Go Programming Language on your system, You would then type the following to compile the code :  
```sh
go build main.go
```

This is to build the source code to be compatible with your own system. To build versions for other systems you can use the `GOOS` switch and/or the `GOARCH` build switches like  
```sh
GOOS=windows GOARCH=amd64 GOAMD64=v3 go build main.go
```

In <strong>specific</strong> architectures you can use an extra flag to specify what hardware you would like this app to be compatible with. Please consult [this resource]('https://go.dev/wiki/MinimumRequirements#amd64') for available options and [this one]('https://go.dev/doc/install/source#environment') for available combinations.

<strong>For specific operating systems, you need the c/c++ compiler and to make sure Go can find it.</strong>

## Thanks
- [@GeorgeChatzogiannakis](https://github.com/GeorgeChatzogiannakis) - Comments, Recommedations, Testing

## Contributing

Please read the **CONTRIB.md** file to give you a basic TODO list\. \:\)

