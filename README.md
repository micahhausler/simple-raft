# Simple Raft implementation

I wanted to try to get my head around Raft leader election so this is a naive
implementation written in go.

## Instructions

You must have docker/docker-compose

```bash
docker-compose up
```

And you'll see a log like this

```log
node2_1  | 2018/01/25 02:15:48 Starting event loop http://node2:3000
node3_1  | 2018/01/25 02:15:48 Starting event loop http://node3:3000
node3_1  | 2018/01/25 02:15:48 starting server
node2_1  | 2018/01/25 02:15:48 starting server
node2_1  | 2018/01/25 02:15:49 Heartbeat
node3_1  | 2018/01/25 02:15:49 Heartbeat
node1_1  | 2018/01/25 02:15:49 Starting event loop http://node1:3000
node1_1  | 2018/01/25 02:15:49 starting server
node2_1  | 2018/01/25 02:15:49 Heartbeat
node2_1  | 2018/01/25 02:15:49 No leader, proposing self
node1_1  | 2018/01/25 02:15:49 Got vote for http://node2:3000
node1_1  | 2018/01/25 02:15:49 No leader, accepting request from http://node2:3000
node2_1  | 2018/01/25 02:15:49 http://node1:3000 voted for me to be leader
node3_1  | 2018/01/25 02:15:49 Got vote for http://node2:3000
node3_1  | 2018/01/25 02:15:49 No leader, accepting request from http://node2:3000
node2_1  | 2018/01/25 02:15:49 http://node3:3000 voted for me to be leader
node2_1  | 2018/01/25 02:15:49 I am the candidate!
node1_1  | 2018/01/25 02:15:49 Got confirm for http://node2:3000
node1_1  | 2018/01/25 02:15:49 Confirming http://node2:3000 as leader
node2_1  | 2018/01/25 02:15:49 http://node1:3000 confirmed me as leader
node3_1  | 2018/01/25 02:15:49 Got confirm for http://node2:3000
node3_1  | 2018/01/25 02:15:49 Confirming http://node2:3000 as leader
node2_1  | 2018/01/25 02:15:49 http://node3:3000 confirmed me as leader
node2_1  | 2018/01/25 02:15:49 I am the Leader!
node1_1  | 2018/01/25 02:15:49 Heartbeat
node3_1  | 2018/01/25 02:15:49 Heartbeat
node3_1  | 2018/01/25 02:15:49 Heartbeat
node1_1  | 2018/01/25 02:15:49 Heartbeat
```
