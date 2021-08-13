Credrpc
=======

Using UNIX socket and `SO_PASSCRED` for privilegied actions.

Somewhere in `man unix` you can find the documentation :

```
SO_PASSCRED
    Enabling  this  socket option causes receipt of the credentials of the sending process in an
    SCM_CREDENTIALS ancillary message in each subsequently received message.  The returned
    credentials are those specified by the sender using SCM_CREDENTIALS, or a default that
    includes the sender's PID, real user ID, and real group ID, if the sender did not specify
    SCM_CREDENTIALS ancillary data.

    When this option is set and the socket is not yet connected, a unique name in the abstract
    namespace will be generated automatically.

    The value given as an argument to setsockopt(2) and returned as the result of getsockopt(2)
    is an integer boolean flag.
```

`SO_PASSCRED` is Linux only, Darwin should use `LOCAL_PEERCRED`, the patch is merged, but not in current Golang version.

Implementation
--------------

The code is über simple, not optimised, with few abstraction, it should be completly read before usage.
If any error happened, connection is closed. Nothing is reused.
The server is designed to not trust the client, but the kernel.

Protocol
--------

Message are Pascal String, 4 bytes for the lentgh, an `uint32`, and n bytes for the message.

RPC
---

The protocol use one shot UNIX socket, one socket per call, without streaming.

There is no routing, just one handler per socket.

The RPC send one message for argument, and get two messages for response : error and payload.
The error is a plain string, if its length is 0, there is no error, so the payload is read.

The serialisation is not part of this project, bring your own `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler`.

The handler get the call argument (a plain old `[]byte`), and unix credential (process ID, User ID, Group ID).

The handler returns a response (an other plain old `[]byte`) and an error.
