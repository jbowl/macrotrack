# macrotrack
Simple REST API to test database drivers with GO



## To run SqlServer in a podman container

SqlServer 2019 runs rootless using a user with ID 10001

>Use **podman unshare chown** to grant the container user (with id 10001) permissions to write to locally mounted directory <br>
> see https://man7.org/linux/man-pages/man1/unshare.1.html regarding unshare <br>
> podman unshare runs a command in Podman's namespace

I think docker will work the same by simply chown-ing the directory.


$ podman unshare chown 10001:10001 HOST_DIRECTORY <br>
Now you can mount the chown-ed directory


```
$ podman run -d -e "ACCEPT_EULA=Y" \
    -e "MSSQL_SA_PASSWORD=Password \
    --name NAME  --hostname HOSTNAME \
	  -v HOST_DIRECTORY:/var/opt/mssql/data:Z \
	  mcr.microsoft.com/mssql/server:2019-latest
```                 
