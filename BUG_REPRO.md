# Bug Reproduction

## What was wrong

Event filtering, searching, aggregation, merging, active-event selection, and statistics shared caller-owned slice storage or label maps. Later operations could reorder earlier results or mutate historical labels.

## How to trigger

Run the event ownership regression checks on the bug branch. They exercise filtering, searching, aggregation, merge, active selection, and statistics in sequence and assert that prior slices and label maps remain unchanged.

## Expected behavior

Every returned event collection owns its slice and label data, merge preserves historical labels, and statistics does not reorder caller input.
