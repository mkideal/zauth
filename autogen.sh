#!/bin/bash

midc ./model.mid -Xmeta -Ogo=. -K=beans

midc ./api.mid -Ogo=. -Tgo=./templates/api/ # --log=debug
