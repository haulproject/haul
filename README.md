# haul

Inventory management system for patchwork components and assets.

## procedure

Deploy mongodb and haul server

`$ docker-compose up -d --build`

Make sure the server is up

`$ docker-compose logs haul`

### optional - local cli

Haul is accessed mainly through it's API.

Install the cli locally.

`$ go install`

Then, add a config locally to `$HOME/.haul.yaml` from the `example/haul.yaml`, or specify a different config location with the `--config` flag. 

*Note that if you installed the server from docker-compose without any modifications, this should be unnecessary to access the server on the same host.*

Once that is done, make sure you can access the server with a healthcheck

`$ haul ping`

You can then list all api routes

`$ haul api`

See `$ haul help` for more options
