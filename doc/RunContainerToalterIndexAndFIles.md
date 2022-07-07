`http_proxy`, `https_proxy` and `no_proxy` can be adjusted.

I'm running a multinode cluster with local storage exposed, hence the `nodeName` parameter is set, such that this "utils" container is spawned on the right node.
This nodeName parameter needs to match the `nodeName` referenced in the _ztp-webserver_ deployment CR.

BE AWARE that the command parameters of _containername_, _image_ and _args_ are being overwritten via the `overrides` provided container spec.
These are mendatory fields for both, commandline as well as overwrites spec. The overwrite spec has precedence.

```
kubectl run -i --rm --tty ztp-webserver-utils --overrides='
{
   "kind":"Pod",
   "apiVersion":"v1",
   "spec":{
      "containers":[
         {
            "name":"ztp-webserver-utils",
            "image":"alpine:latest",
            "args":[
               "sh"
            ],
            "stdin":true,
            "stdinOnce":true,
            "tty":true,
            "volumeMounts":[
               {
                  "mountPath":"/webserver",
                  "name":"ztp-webserver-storage"
               }
            ],
            "env":[
               {
                  "name":"HTTP_PROXY",
                  "value":"http://135.245.192.7:8000"
               },
               {
                  "name":"HTTPS_PROXY",
                  "value":"http://135.245.192.7:8000"
               },
               {
                  "name":"no_proxy",
                  "value":"192.168.200.10,192.168.122.210,192.168.122.253,192.168.122.233"
               }
            ]
         }
      ],
      "nodeName":"node2",
      "volumes":[
         {
            "name":"ztp-webserver-storage",
            "persistentVolumeClaim":{
               "claimName":"ztp-webserver-claim"
            }
         }
      ]
   }
}'  --image=alpine:latest --restart=Never -- sh
```