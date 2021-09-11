#!/bin/sh


docker run -it --rm -v /home/${X_Forwarded_User}:/home/${X_Forwarded_User} -v /etc/passwd:/etc/passwd -v /etc/group:/etc/group --v /etc/netauth:/etc/netauth:ro --user "${X_Forwarded_Uid}:${X_Forwarded_Uid}" alpine:3.14
