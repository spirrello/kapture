# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [0.2.11] - 2019-10-23

### Changed
- corrected nodeport in configmap
- imported ENV vars in nodeapi
- Testing linux commands in nodeapi
- added network policies
- looping through all kcapture pods to post pod data


## [0.2.10] - 2019-10-19

### Changed
- nodeapi port
- namespace to kcapture
- hostNetwork for kcapture-node


## [0.2.9] - 2019-10-16

### Changed
- Shared package for common functions
- Compound Docker files to separate build vs application image
- Testing exec commands

## [0.2.8] - 2019-10-14

### Changed
- Using slice in nodeapi instead of map


## [0.2.7] - 2019-10-14

### Fixed
- Fix nodeapi call

## [0.2.6] - 2019-10-14

### Changed
- start posting to node-api again


## [0.2.5] - 2019-10-14

### Changed
- use kcapture-node ip instead of name


## [0.2.4] - 2019-10-14

### Changed
- Updated kcapture-node to no copy bashrc or Dockerfile

### Changed
- busybox for kcapture-node


## [0.2.2] - 2019-09-21

### Changed
- renamed podinfo.go to common.go
- moved LogFormat struct to models package

## [0.2.1] - 2019-09-21

### Changed

- updated logging
- readme
- healthcheck response

## [0.2.0] - 2019-09-14

### Added
- struct for all pod info

### Changed
- moved pod struct to package
- using ENV for go-client connection

## [0.1.4] - 2019-09-08

### Added
- print entire pod object

## [0.1.4] - 2019-09-08

### Added
- Stop printing pod status

## [0.1.3] - 2019-09-08

### Added
- Test pod filter

## [0.1.2] - 2019-09-08

### Changed
- update changelog

## [0.1.1] - 2019-09-08

### Changed
- Remove dynamic links

## [0.1.0] - 2019-08-26

### Added
- Jenkinsfile
- K8s in cluster client