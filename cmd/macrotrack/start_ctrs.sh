#!/bin/bash

db="mssqldb"
# start sql server
	#running=$( podman ps -a -f name=$db | grep $db 2> /dev/null )
	running=$( podman ps -f name=$db | grep $db 2> /dev/null )
	if [[ -z ${running} ]]; then 
	  podman run -d -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=___Aa123" --name $db --hostname $db  -p 1435:1433  \
	  -v /home/j/jsoft/github.com/jbowl/macrotrack/macros/data/sqlserver:/var/opt/mssql/data:Z \
	  mcr.microsoft.com/mssql/server:2019-latest
    else
      echo "mssqldb running"  
	fi

# postgres
db="postgresdb"

#running=$( podman ps -a -f name=$db | grep $db 2> /dev/null )
running=$( podman ps -f name=$db | grep $db 2> /dev/null )
echo $running
if [[ -z ${running} ]]; then 
    podman run -d --rm --name $db \
      -e "POSTGRES_PASSWORD=postgres" \
	  -e "POSTGRES_USER=postgres" \
	  -e "POSTGRES_DB=postgres" \
	  -p 5434:5432 \
	  -v /home/j/jsoft/github.com/jbowl/macrotrack/macros/data/postgres:/var/lib/postgresql/data:Z postgres \
	  postgres
else
      echo "postgres running"  
fi

# start postgres 



