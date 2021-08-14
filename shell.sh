#!/bin/sh


PASSWD_ENTRY="${X_Forwarded_User}:x:${X_Forwarded_Uid}:${X_Forwarded_Gid}:${X_Forwarded_Fullname}:/:/bin/sh"

echo "$PASSWD_ENTRY" >> /etc/passwd

docker run -it --rm -v /etc/passwd:/etc/passwd --user "${X_Forwarded_Uid}" alpine:3.14