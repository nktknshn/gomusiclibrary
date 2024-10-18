#!/bin/bash

cd /home/user/Workspace/MusicLibrary/golib;
rm data.db
./scripts/migrate.sh; ./scripts/run.sh scan ~/Library
