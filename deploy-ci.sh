#!/bin/sh

now -t ${NOWSHTOKEN} -e PS_GOOGLE_API=@PS_GOOGLE_API
now alias podcasts.psmarcin.me -t ${NOWSHTOKEN}
now rm podcasts.psmarcin.me --safe --yes -t ${NOWSHTOKEN}
