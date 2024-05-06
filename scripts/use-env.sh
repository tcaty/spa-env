#!/bin/bash

dotenv=$1

export $(grep -v '^#' $dotenv | xargs -d '\n')
