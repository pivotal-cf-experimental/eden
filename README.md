# eden for [topham](https://github.com/pivotal-cf-experimental/topham-controller)

A fork of eden, modified to store no state locally, and instead rely on topham-controller as a stateful services controller.

Provides a CF-like workflow (provision/bind/unbind/deprovision with no state stored on your machine), minus the CF.


## Installation

```
go get -u gopkg.in/pivotal-cf-experimental/eden.v1
```


## Usage

Use environment variables to target your topham-controller server:

```bash
export SB_BROKER_URL=https://topham-controller.com
export SB_BROKER_USERNAME=topham-username
export SB_BROKER_PASSWORD=topham-password
```

To see the available services and plans:

```bash
eden catalog
```

To create (`provision`) a new service instance, and to generate a set of access credentials (`bind`):

```bash
eden provision -s servicename -p planname -i my-db-name # all flags are required
eden bind
```

To view the credentials for your binding:

```bash
eden credentials
```
