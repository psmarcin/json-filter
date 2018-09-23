#!/bin/sh

now -t ${NOWSH_TOKEN}
now alias podcasts.psmarcin.me -t ${NOWSH_TOKEN}
now rm podcasts --safe --yes -t ${NOWSH_TOKEN}
