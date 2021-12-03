### Install Commune
This a detailed guide on how to get Commune up and running on your local machine.

#### Go
Follow the instructions on the official Golang site [here](https://golang.org/doc/install).

#### Synapse
Follow the installation instructions from the Synapse docs [here](https://matrix-org.github.io/synapse/latest/setup/installation.html). Read the section about the directory structure, and config before install Synapse. 
- Update `homeserver.yaml` to match the `shared_secret` code in your `config.toml` file.

#### Redis
Install redis using your system's package manager. Ensure that redis has password authentication. You can update the redis password in the config file `/etc/redis/redis.conf`:
```
requirepass mypassword
```
This password must match the one in `config.toml`.

#### Postgres
Use your system's package manager to install postgres. Create a user called `commune` with superuser privileges.
```
sudo -u postgres createuser -s -i -d -r -l -w commune
```

#### NPM
User your system's package manager, or something like [nvm](https://github.com/nvm-sh/nvm)

### Utilites
Commune requires several utilities for building the project.

Install [modd](https://github.com/cortesi/modd) for hot-reloading the go backend.
```
go get github.com/cortesi/modd/cmd/modd
```

Install [goose](https://github.com/pressly/goose) for database migrations.
```
go get -u github.com/pressly/goose/v3/cmd/goose
```
  
#### Directory structure  
Commune uses relative paths for managing the synapse instance from the make tasks. You'll need to make sure your directory structure is like this:
- Create a directory called `commune`
- Clone the repository inside this directory, `commune/app`
- Install Synapse in `commune/synapse`

#### Config
Copy the `config-sample.toml` file to `config.toml`. Update the config file with relevant flags:
- Redis password
- Shared secret key for Synapse
- Additionally, you might want to include your own Youtube API key for embed previews, or your own Tenor API key for GIF access.

### Running commune
To start Commune, cd into `commune/app` and run:
```
make dev
```
If all goes well, Commune should be running. If not, you should see error messages.
