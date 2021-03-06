# `mirror-mirror`

Check yourself out. Lookin' good.

Go web server that lets you know where you're coming from. Incase you forgot.

## Install

```
go get github.com/ajpatri/mirror-mirror
```

## Usage

```bash
$ mirror-mirror --help
Usage of mirror-mirror:
  -host string
        Host address to listen on (default "0.0.0.0")
  -https
        Serve over HTTPS
  -port int
        Port to listen on (default 8080)
  -private string
        Private key (.pem) - Requires https flag
  -public string
        Public key (.pem) - Requires https flag

# Over HTTPS
$ mirror-mirror -https -public cert.pem -private key.pem

# Only locally
$ mirror-mirror -host 127.0.0.1
```

## Docker Usage

```
git clone git@github.com:ajpatri/mirror-mirror.git
cd mirror-mirror
sudo docker build -t mirror-mirror .
sudo docker run -d -p 127.0.0.1:8080:8080 mirror-mirror
```

