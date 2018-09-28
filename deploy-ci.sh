#!/bin/sh

now -t ${NOWSHTOKEN}
now alias podcast.psmarcin.me -t ${NOWSHTOKEN}
now rm podcast.psmarcin.me --safe --yes -t ${NOWSHTOKEN}
