#!/bin/bash

midc ./model.mid -Xmeta -Ogo=. -K=storage

midc ./api.mid -Ogo=. -Tgo=./templates/api/ # --log=debug
