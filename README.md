# RSS feed creator
The RSS feed creator does exactly what it is called, it creates an RSS feed which is accessible on port 8082 and can be filled using the feed item API.

## Setup
For local running and testing purposes:
```
make run
```
This will spin up a `docker-compose` workload containing the RSS feed creator and a MongoDB (v6) instance. You can alter the credentials in the `build/docker-compose.yml` file for both instances.

For running a production workload it is recommended to have a dedicated database for persistence, for which you can find an example in the `build/docker-compose-with-external-db.yml` file. Kubernetes is also supported, for which an example can be found under `build/k8s-deployment.yml`. 
Currently the following databases are supported:
- MongoDB 4
- MongoDB 6 (default)
- MSSQL (SQL Server)

The following environment variables can be set during deployment to alter the application.
When using a MongoDB database:
```bash
DATABASE_TYPE=mongodb6 # required, can also be mongodb4 
MONGODB_URI="mongodb://username:password@localhost:27017" # required when using mongodb
```

When using a MSSQL database:
```bash
DATABASE_TYPE=mssql # required when using mssql
MSSQL_CONNECTION_STRING=server=exampleserver.database.windows.net;user id=testadmin;password=**********;port=1433;database=rssfeed-db; # required when using mssql
```

Other parameters:
```bash
BASICAUTH_USERNAME=test # required, the username you want for the basic authentication
BASICAUTH_PASSWORD=test # required, the password you want for the basic authentication
RSS_TITLE=testfeed # optional, the title that the RSS feed will have
RSS_DESCRIPTION=Test feed description # optional, the RSS feed description
```

## Usage
An example of adding an item to the RSS feed can be found here:
```
curl -X POST -d @examples/test.json http://localhost:8082/api/
```
