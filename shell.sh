#!/bin/sh


PASSWD_ENTRY="${X_Forwarded_User}:x:${X_Forwarded_Uid}:${X_Forwarded_Gid}:${X_Forwarded_Fullname}:/:/bin/sh"

id "${X_Forwarded_user}" || echo "$PASSWD_ENTRY" >> /etc/passwd

docker run -it --rm -v /etc/passwd:/etc/passwd --user "${X_Forwarded_User}" alpine:3.14