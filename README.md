# go-curator
==========

A port of the netflix curator zookeeper framework, mostly for service discovery
Provides a zookeeper client with a high tolerance for disconnections. It optionally hooks into the netflix exhibitor
framework, so that it can reconfigure itself as the zookeeper ensemble changes.

## Ensemble providers

An ensemble provider interface represents a dynamic connection string. The library detects when there are changes in the
connection string and will reconnect as needed.

## Event bus

You can register channels as event handlers for connection events. When the message sent to a handler channel isn't
handled withing 50ms (by default) then the handler gets unregistered and a warning is logged. 


### How it works

CuratorConnection 
  * every operation is handled in a retry loop
  * a connection is watched and when a session event occurs an appropriate action is taken, on expired session it reconnects
  * the retry loop does: 
     * check for connection string changes, on change reconnect
     * try the operation when fails check policy (predicate) and retry or bail
  * consumers can register channels as event handlers for connection events so they know how to change internal state

Node Cache
  * caches a single node 

Children Cache
  * caches a map of children

