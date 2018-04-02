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
TBD

License
-------------
The "thumbnail service" is licensed under the Apache License 2.0. Please see the LICENSE file for details.
