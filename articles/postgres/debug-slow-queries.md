# postgres / debug slow queries

If the query is slow, I add this to the top of the file:

```
EXPLAIN (ANALYZE, COSTS, VERBOSE, BUFFERS, FORMAT JSON)
```

Then, run:

```bash
:!psql -qAt -d $(cat .db) -f % | pbcopy
```

Paste into <http://tatiyants.com/pev/#/plans/new>
and delete the trailing line to make it valid JSON:

```
Time: 1111.111 ms (00:01.111)
```

The output is an interactive visualization that makes it
easy to identify which parts of the query are
slowest, largest, and costliest.

![EXPLAIN visualizer](/images/postgres-explain-visualizer.png)
