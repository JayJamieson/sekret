#!/bin/bash

(cd ui && pnpm run build)

fly deploy
