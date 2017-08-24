#!/bin/bash
pkgname=server1.tar.gz

if [ -f $pkgname ] ; then
    rm -f $pkgname
fi

tar zcvf $pkgname microserver config.ini run.sh