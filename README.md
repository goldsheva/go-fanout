# go-fanout

A simple HTTP request relay to multiple endpoints.

## Build & Run App

    git clone git@github.com:goldsheva/go-fanout.git
    cd go-fanout && cp .env.sample .env
    go build -o app ./cmd
    ./app

## Systemd config
 

    [Unit]
    Description=Go-fanout
    After=network.target
    
    [Service]
    Type=simple
    WorkingDirectory=/path-to-binary-app
    ExecStart=/path-to-binary-app/app
    Restart=always
    RestartSec=5
    User=www-data
    
    [Install]
    WantedBy=multi-user.target
