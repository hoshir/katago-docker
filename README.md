# katago-docker
A Docker image to run KataGo in GCP


# Important notice
Please don't forget to stop your GCP instance to avoid incurring charges to your account when you don't use it. 



# Quick Start
## Launch a new instance 
See instructions in Google Cloud Platform to launch an instance with GPU.

This Docker image is tested with the "GPU Optimized Debian m32 (with CUDA 10.0)" image with NVIDIA Tesla K80 GPU.

## SSH to the instance


## Set up NVIDIA Container Toolkit
You need to install NVIDIA Container Toolkit to build run GPU accelerated Docker containers.

See [Quick Start](https://github.com/NVIDIA/nvidia-docker#quickstart) section in the [nvidia-docker](https://github.com/NVIDIA/nvidia-docker) page.

You will need to run following commands for Debian Stretch in the instance.
```sh
# Add the package repositories
$ distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
$ curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
$ curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list

$ sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
$ sudo systemctl restart docker
```

## Prepare a model
Download one of the model files from [the release page](https://github.com/lightvector/KataGo/releases) and extract it.

For example, 
```sh
$ wget https://github.com/lightvector/KataGo/releases/download/v1.1/g104-b15c192-s297383936-d140330251.zip
$ unzip g104-b15c192-s297383936-d140330251.zip
```

## Prepare a config file
If you would like to run KataGo as an engine for Go client such as Lizzie, you will want to use gtp_example.cfg in the KataGo [repository](https://github.com/lightvector/KataGo) as a base line.

For example, 
```sh
$ wget https://raw.githubusercontent.com/lightvector/KataGo/master/cpp/configs/gtp_example.cfg
```

## Running KataGo from the Docker image
Now, you are ready to run KataGo Docker image. Run the follwoing command.

```sh
$ docker run --rm -it \
    -v $(pwd):/data \
    --runtime nvidia \
    hoshir/katago-v1.2-cuda10.0-linux-x64 gtp \
    -model /data/g104-b15c192-s297383936-d140330251/model.txt.gz \
    -config /data/gtp_example.cfg \
    -override-version 0.17

```

If KataGo works successfully, it will show its version and some outputs like this.
```
KataGo v1.2
Loaded model /data/g104-b15c192-s297383936-d140330251/model.txt.gz
GTP ready, beginning main protocol loop
```

Try some `showboard` GTP command, and quit.
```sh
KataGo v1.2
Loaded model /data/g104-b15c192-s297383936-d140330251/model.txt.gz
GTP ready, beginning main protocol loop
showboard
= MoveNum: 0 HASH: 7BF12F3F24903F0C225CD6C55BA8BA1F
   A B C D E F G H J K L M N O P Q R S T
19 . . . . . . . . . . . . . . . . . . .
18 . . . . . . . . . . . . . . . . . . .
17 . . . . . . . . . . . . . . . . . . .
16 . . . . . . . . . . . . . . . . . . .
15 . . . . . . . . . . . . . . . . . . .
14 . . . . . . . . . . . . . . . . . . .
13 . . . . . . . . . . . . . . . . . . .
12 . . . . . . . . . . . . . . . . . . .
11 . . . . . . . . . . . . . . . . . . .
10 . . . . . . . . . . . . . . . . . . .
 9 . . . . . . . . . . . . . . . . . . .
 8 . . . . . . . . . . . . . . . . . . .
 7 . . . . . . . . . . . . . . . . . . .
 6 . . . . . . . . . . . . . . . . . . .
 5 . . . . . . . . . . . . . . . . . . .
 4 . . . . . . . . . . . . . . . . . . .
 3 . . . . . . . . . . . . . . . . . . .
 2 . . . . . . . . . . . . . . . . . . .
 1 . . . . . . . . . . . . . . . . . . .

quit
= 

$
```



# Communicating with the engine from your local machine
You can communicate with the engine from your local machine with SSH.

## Register you SSH key to the instance


## Create a launcher script in your local machine
Create a small shell script named `katago.sh` by adding GCP_USER and GCP_IP in it.


```sh
#!/bin/bash

set -eu

# User name of the GCP insance
GCP_USER=

# External IP address of the GCP instance
GCP_IP=

ssh -o 'StrictHostKeyChecking no' $GCP_USER@$IP \
  docker run \
    -i \
    --rm \
    --runtime nvidia \
    -v /home/$GCP_USER:/data \
    hoshir/katago-v1.2-cuda10.0-linux-x64 gtp \
    -model /data/g104-b15c192-s297383936-d140330251/model.txt.gz \
    -config /data/gtp_example.cfg \
    -override-version 0.17
```
The script assumes
  - The model file is stored as `/home/$GCP_USER/g104-b15c192-s297383936-d140330251/model.txt.gz` in the GCP instance.
  - The config file is stored as `/home/$GCP_USER/gtp_example.cfg` in the GCP instance.

The command mounts `/home/$GCP_USER` directory in the GCP instance as `/data` when running the image.

## Register the script to Go client
Register the script to your Go client. If you use Lizzie, the entry of the engine will be `./katago.sh --weights %network-file`. You will need to add `--weights %network-file` argument since Lizzie checks it.

