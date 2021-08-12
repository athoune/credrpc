Chown me
========

Using UNIX socket and `SO_PASSCRED` for privilegied actions.

Somewhere in `man unix` you can find the documentation :

```
SO_PASSCRED
    Enabling  this  socket option causes receipt of the credentials of the sending process in an SCM_CRE‐
    DENTIALS ancillary message in each subsequently received message.  The returned credentials are those
    specified by the sender using SCM_CREDENTIALS, or a default that includes the sender's PID, real user
    ID, and real group ID, if the sender did not specify SCM_CREDENTIALS ancillary data.

    When this option is set and the socket is not yet connected, a unique name in the abstract  namespace
    will be generated automatically.

    The value given as an argument to setsockopt(2) and returned as the result of getsockopt(2) is an in‐
    teger boolean flag.
```

`SO_PASSCRED` is Linux only, Darwin should use `LOCAL_PEERCRED`.
