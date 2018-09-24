#!/bin/sh

now -t ${NOWSHTOKEN}
now alias podcasts.psmarcin.me -t ${NOWSHTOKEN}
now rm podcasts.psmarcin.me --safe --yes -t ${NOWSHTOKEN}
