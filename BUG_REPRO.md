# Bug Reproduction

## What was wrong

Queue publish, idle consumption, retry backoff, worker-pool shutdown, and the no-op consumer did not consistently propagate caller cancellation. Cancelled work could remain queued or keep workers alive.

## How to trigger

Run the queue cancellation regression checks on the bug branch. They cancel publishing, handlers, idle consumers, retry waits, worker pools, and no-op consumers at each boundary.

## Expected behavior

Every queue operation observes the caller context, returns its cancellation cause, and stops without retaining cancelled messages or workers.
