twiph
=====

``twiph`` draws a fancy graph representation of Twitter connections given a list of names/surnames.

Authentication
==============

You need to create your own app: https://apps.twitter.com/app/new

Then an ``Access Token`` must be created https://apps.twitter.com/app/<app-id>/keys

Examples
========

Once retrieved auth credentils and created a ``cfg`` (a sample is provided in ``conf/sample.cfg``) file::

        export CFGP_FILE_PATH=./conf/mine.cfg;
        twiph -csv list.csv

Run tests
=========

Tests need Twitter api credentials. You can specify them by command line or in configuration file::

        export CFGP_FILE_PATH=./conf/test.cfg; go test

Notes
=====

Graphs with nodes >= ~100 start to be *very* heavy (cpu intensive) to display.

Credits
=======

The amazing ``D3.js`` is used to create graphs http://d3js.org/

Disclaimer
==========

All trademarks, copyrights and other forms of intellectual property belong to their respective owners.

The author is not affiliated with any vendor cited above.
