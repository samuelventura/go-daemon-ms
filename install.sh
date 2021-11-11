#!/bin/bash -x

#journalctl -u GoDaemonMs.service
#journalctl -fu GoDaemonMs.service
#sudo systemctl stop GoDaemonMs
#sudo /usr/local/bin/go-daemon-ms -service uninstall
if [[ "$OSTYPE" == "linux"* ]]; then
    SRC=$HOME/go/bin
    DST=/usr/local/bin
    if [[ -f "$DST/go-daemon-ms" ]]; then
        sudo systemctl stop GoDaemonMs
        sudo $DST/go-daemon-ms -service uninstall
        sleep 1
    fi
    go install
    sudo cp $SRC/go-daemon-ms $DST
    sudo $DST/go-daemon-ms -service install
    sudo systemctl restart GoDaemonMs
    sudo systemctl status GoDaemonMs
fi
