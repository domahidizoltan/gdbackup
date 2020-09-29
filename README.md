# gdbackup

A tool to download and backup your Google Drive (GD) contents.  

:warning: This is a pet project, not a professional OSS. Use it at your own risk! It may contain bugs.
> * This tool was tested only on Ubuntu Linux. 
> * GDBackup is using read-only permissions on your Google Drive, but it creates directories and files on your local environment.  

<br/>

## Grant permissions

* Create a project on [Google Console](https://console.cloud.google.com/)
* Go to [APIs & Services/Credentials](https://console.cloud.google.com/apis/credentials)
* Create an OAuth client ID and download the json file next to `gdbackup` executable as `credentials.json`

<br/>

## Usage

### The gdignore.yaml file

gdignore.yaml gives you some level of control to avoid downloading your whole Google Drive content.

Let's take the `gdignore.yaml` example:
```
Documents: 
    - public
    - private:
        - ignore-file1
        - ignore-file2
        - office*.zip
        - '*.mp4'
Images:
    - rawfiles
    - public:
        - office*
root:
    - shared-for-gmail.txt
```

* Only folders on your GD root level will be included in the downloads having at the root level of the yaml. So if you have *Document*, *Images* and *Videos* on the root level of your GD, only *Documents* and *Images* will be downloaded
* Files on the root level will be ignored by default except if you define `root` on the root level of the yaml
* Values of the yaml file can be folders or files.
* Path to a yaml value is the same as the path to your GD file. *Documents/private/ignore-file1* is the corresponding GD file what should not be downloaded
* Only yaml values will be ignored, so the rest of the *Documents/private* folder will be downloaded
* Ignore values can have a single * wildcard, so these files will be also ignored:
  > *Documents/private/office1.zip*  
  > *Documents/private/office2.zip*  
  > *Documents/private/longvideo.mp4*  
  > *Images/public/office-pics.zip*  

  however these will be downloaded:
  > *Documents/private/office1.doc*  
  > *Documents/private/longvideo.mpg*  
  > *Images/public/outofoffice.zip*  
  

### Program arguments

  * *-backup-path*: path where to download files
  * *-delay*: an interval in seconds to randomly wait before a file download
  * *-gdignore-path*: path to custom *gdignore.yaml* file
  * *-loglevel*: for debugging purpose
  * *-max-parallel-downloads*: max number of parallel downloads can be used

*delay* and *max-parallel-downloads* let you control the download rate if you would be banned for excessive downloads  

See `./gdbackup -h` for help

<br/>

## Build and Run

build:  
```
go build -ldflags="-w -s" -o gdbackup main/*
```
Use [goupx](https://github.com/pwaller/goupx) to compress the executable from ~10MB to ~4MB

run:
```
./gdbackup
```

help:
```
./gdbackup -h
```


*gdbackup* will create a folder containing the current date i.e. `gdbackup20200928` where the files will be downloaded.