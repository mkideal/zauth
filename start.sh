#!/bin/bash

app=./authd
pidfile=./authd.pid
log_level=trace
log_providers=file
log_opts="dir=./log"
etcd="127.0.0.1:2379"
service_name="authd"
third_party="wechat,qq"
addr=0.0.0.0:5200
mode=release
cookie_key=authd_cookie
session_expire_duration=3600
html=html
html_router=/

$app daemon \
	--pid $pidfile \
	--log-level $log_level \
	--log-providers $log_providers \
	--log-opts $log_opts \
	--etcd $etcd \
	--service-name $service_name \
	--third-party $third_party \
	--addr $addr \
	--mode $mode \
	--cookie $cookie_key \
	--session-expire-duration $session_expire_duration \
	--html $html \
	--html-router $html_router 
