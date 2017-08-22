#!/bin/bash

# README: https://golang.org/doc/install/source#environment
#
# This build script basically does:
# 
# echo 'building windows, amd64'
# rm -rf ./windows/amd64
# env GOOS=windows GOARCH=amd64 go build -o ./windows/amd64/geeny
#
# You can add extra steps in `main`

### params ####

#system
readonly FALSE=1
readonly TRUE=0
readonly SCRIPT_NAME=`basename "$0"`
FORCE=$FALSE

#program locations
readonly env="/usr/bin/env"
readonly go=$(which go)
readonly rm="/bin/rm"
readonly date=$(which date)
readonly git=$(which git)

#os
readonly WINDOWS="windows"
readonly DARWIN="darwin"
readonly LINUX="linux"

#arch
readonly AMD64="amd64"

### functions ###

function log {
	echo "  > $SCRIPT_NAME: $1" > /dev/stderr
}

function logf {
	printf "  > $SCRIPT_NAME: $1" > /dev/stderr
}

function checkArgs {
	local readonly PARAM_COUNT=2
	if [ "$#" -ne $PARAM_COUNT ]; then
    	log "illegal number of parameters. found:[$#], required:[$PARAM_COUNT]"
    	return $FALSE
	fi
	if [ "$1" -ne "$2" ]; then
    	log "illegal number of parameters. found:[$1], required:[$2]"
    	return $FALSE
	fi
	return $TRUE
}

function getPermissionAndOverride {
	if [ $FORCE = $TRUE ]; then
		return $TRUE
	fi

	local readonly PARAM_COUNT=1
	if ! $(checkArgs $# $PARAM_COUNT); then
		return $FALSE
	fi

	local readonly LOC=$1
	if [ -f "$LOC" ]; then
		logf "[$LOC] exists. delete [$LOC] and resave? [y/n] "
		read decision
		if [ "$decision" = "y" ]; then
			$rm -rf $LOC
			return $TRUE
		else
			log "did not delete:[$LOC]"
			return $FALSE
		fi
	fi
	return $TRUE
}

function build {
	local readonly PARAM_COUNT=3
	if ! $(checkArgs $# $PARAM_COUNT); then
		return $FALSE
	fi
	local readonly OS="$1"
	local readonly ARCH="$2"
	local readonly LOC="$3"

	if $(getPermissionAndOverride $LOC); then
		local readonly BUILD_MSG="OS:[$OS], ARCH:[$ARCH], LOC:[$LOC]"
		log "building $BUILD_MSG..."
		cd src/geeny
		$env GOOS=$OS GOARCH=$ARCH $go build -ldflags "-X version.timestamp=$($date -u '+%d.%m.%Y_%H:%M:%S') -X version.version=$($git describe --tags)" -o $LOC
		cd ..
		log "built $BUILD_MSG"
	fi
	return $TRUE
}

function main {
	local readonly PARAM_COUNT=1
	if ! $(checkArgs $# $PARAM_COUNT); then
		return $FALSE
	fi
	if ! [ -d $1 ]; then
		log "this is not a valid directory: [$1]"
		return $FALSE
	fi

	local PATH="$1/$WINDOWS/$AMD64"
	if ! $(build $WINDOWS $AMD64 $PATH); then
		return $FALSE
	fi

	local PATH="$1/$DARWIN/$AMD64"
	if ! $(build $DARWIN $AMD64 $PATH); then
		return $FALSE
	fi

	local PATH="$1/$LINUX/$AMD64"
	if ! $(build $LINUX $AMD64 $PATH); then
		return $FALSE
	fi

	return $TRUE
}

### script ###

if [ "$1" = "-h" ]; then
	echo "\n\
  Help:\n\n\
  save binaries to given location:
    $SCRIPT_NAME [-f: force all] <PATH_TO_SAVE_BINARIES>\n
  display this help:
    $SCRIPT_NAME -h\n"
	exit 0	
fi

ARGS=$@
if [ "$1" = "-f" ]; then
	FORCE=$TRUE
	ARGS=$2
fi

if $(main $ARGS); then
	log "finished"
	exit 0
else
	log "failed"
	exit -1
fi 
