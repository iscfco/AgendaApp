# Agenda APP
App web para agendar pedidos

## How to Execute:
```
docker-compose -f .\artifacts\docker-compose.yaml down   
docker-compose -f .\artifacts\docker-compose.yaml build --no-cache 
docker-compose -f .\artifacts\docker-compose.yaml up 
```


## Install Invoke:
```
pip install invoke
```

You will see something like:
```
Downloading invoke-3.0.3-py3-none-any.whl (160 kB)
Installing collected packages: invoke
  WARNING: The scripts inv.exe and invoke.exe are installed in 'C:\Users\progr\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts' which is not on PATH.
  Consider adding this directory to PATH or, if you prefer to suppress this warning, use --no-warn-script-location.
Successfully installed invoke-3.0.3

[notice] A new release of pip is available: 26.1.2 -> 26.2.1
[notice] To update, run: C:\Users\<Widows_UserName>\AppData\Local\Microsoft\WindowsApps\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\python.exe -m pip install --upgrade pip
```

Update PATH:

