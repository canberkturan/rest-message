# Rest-Message
Rest-Message is an Offline RESTful Messaging API that gives you messaging ability on your LAN.

## How to install and run?
- If you didn't install, you need to install go programming language firstly with running commands like below as root(sudo):
> apt install go

> pacman -S go
- If you didn't set GOPATH and GOBIN variables you need to set(<i>if you use different shell instead of BASH you need to change ~/.bashrc to your shell's config file like ~/.zshrc</i>):
> mkdir -p \~/.local/go/bin

> echo "export GOPATH=~/.local/go" >> ~/.bashrc

> echo "export GOBIN=$GOPATH/bin" >> ~/.bashrc

> source .bashrc
- After that you need to clone this repository: 
> git clone https://github.com/canberkturan/rest-message
- When cloning process has been finished, enter project path and get dependencies of the project: 
> cd rest-message && go get
- After that build the project
> go build
- And RUN!
> ./rest-message

## How to use this api
- You can read the manual at index page
> https://localhost:8080/
- Probably you will give an warning because of self-signed tls certificate. You can safely ignore that warning and continue.
- If you don't want to use https connection, just edit main.go file and change <i>IS_SECURE</i> variable's value to false.
