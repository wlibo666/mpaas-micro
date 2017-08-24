#!/bin/bash
pkgname=server2.tar.gz

if [ -f $pkgname ] ; then
    rm -f $pkgname
fi

tar zcvf $pkgname httpserver run.sh