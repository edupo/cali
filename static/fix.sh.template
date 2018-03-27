#!/bin/sh

# Get standard cali USER_ID variable
USER_ID={{.Uid}}
GROUP_ID={{.Gid}}

distro=`/bin/cat /proc/version`

case $distro in
  *Ubuntu*)
    useradd -ms /bin/sh user

    # Change 'user' uid to host user's uid
    if [ ! -z "$USER_ID" ] && [ "$(id -u user)" != "$USER_ID" ]; then
      # Create the user group if it does not exist
      groupadd --non-unique -g "$GROUP_ID" group

      # Set the user's uid and gid
      usermod --non-unique --uid "$USER_ID" --gid "$GROUP_ID" user
    fi
    chown -R user: /home/user
    ;;
  *)
    echo "Distro $distro is currently not supported."
    exit 1
esac