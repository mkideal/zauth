#!/bin/bash

app=./accountd
daemon=true
pidfile=./accountd.pid
log_level=trace
log_providers=file
log_opts="dir=./log"
etcd="127.0.0.1:2379"
service_name="accountd"
third_party="wechat/qq"
addr=127.0.0.1:5200
mode=release
cookie_key=accountd
session_expire_duration=3600
html_dir=html
html_router=/

$app -d $daemon \
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
	--html-dir $html_dir \
	--html-router $html_router \
