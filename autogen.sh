#!/bin/bash

midc ./model.mid -Ogo=. -K=storage

midc ./api.mid -Ogo=. -Tgo=./templates/api/ # --log=debug
