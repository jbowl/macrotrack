#!/bin/bash 


# if container hasn't started do so

output=$( podman ps -a -f name=sqlserver_test | grep sqlserver_test 2> /dev/null )
if [[ -z ${output} ]]; then 
  podman run -d -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=___Aa123" --name sqlserver_test --hostname sqlserver_test --rm -p 1434:1433 mcr.microsoft.com/mssql/server:2019-latest
else  
fi







