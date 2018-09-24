#!/bin/sh

now -t ${NOWSHTOKEN} -e PS_GOOGLE_API=@ps_google_api
now alias podcasts.psmarcin.me -t ${NOWSHTOKEN}
now rm podcasts.psmarcin.me --safe --yes -t ${NOWSHTOKEN}
