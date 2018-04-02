# Cloudinary Assignment
Introduction
------------
coded by Moshe Tubul 2018

Compatibility
-------------
The Application was tested on:
1. Ubuntu 17.10

Go Version: 1.9.2

Installation and usage
----------------------

To download project run:

    go get github.com/go-yaml/yaml
    go get github.com/moshetbl/go
    go get -u github.com/disintegration/imaging
    
Compile the project:

    go build -o thumbnail -a main.go
    
Run:

    ./thumbnail ../Config/config.yaml

Tests
-------------
* For now tests are not fully implemented, just a simple example of tests

To run tests:

    cd HttpServices
    go test
 
    
Deployment:
-------------
Heroku Deployment Procecdure:

Install heroku and glide packages:

    sudo add-apt-repository "deb https://cli-assets.heroku.com/branches/stable/apt ./"
    curl -L https://cli-assets.heroku.com/apt/release.key | sudo apt-key add -
    sudo apt-get update
    sudo apt-get install heroku
    sudo add-apt-repository ppa:masterminds/glide
    sudo apt-get update
    sudo apt-get install glide
    
Deploy:

    cd $GOPATH/github.com/moshetbl/thumbnail
    heroku create -b https://github.com/heroku/heroku-buildpack-go.git
    glide create
    glide install
    echo "web: thumbnail Config/config.yaml" > Procfile
    git add glide.yaml glide.lock Procfile
    git commit -m "glide Procfile"
    git push heroku master
    heroku open

License
-------------
The "thumbnail service" is licensed under the Apache License 2.0. Please see the LICENSE file for details.
