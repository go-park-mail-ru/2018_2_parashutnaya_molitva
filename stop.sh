#!/bin/bash

tag="kekmateg[o]"

kill $(ps aux | grep "$tag" | awk '{print $2}')
