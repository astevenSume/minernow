#!/bin/bash

# game
bin/dbgenerate --models-dir src/utils/game/models/ \
    --sql-dir src/utils/game/sql/ \
    --config src/tools/dbgenerate/game.json \
    && echo "execute game dbgenerate succeed!"