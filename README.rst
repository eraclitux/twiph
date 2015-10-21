twiph
=====

``twiph`` draws fancy interactive graphs representation of Twitter connections (aka friends) given a list of names/surnames. It could be uselful to spot and analyze social connections whithin an arbitrary list of people.

Authentication
==============

You need to create your own app: https://apps.twitter.com/app/new

Then an ``Access Token`` must be created https://apps.twitter.com/app/<app-id>/keys

Examples
========

These are the output for the 100 younger italian parliamentarians grouped by political parties:

|image0|_ |image1|_

.. |image0| image:: http://www.eraclitux.com/img/twiph_100_younger.png
.. _image0: http://www.eraclitux.com/img/twiph_100_younger.png

and with avatars:

.. |image1| image:: http://www.eraclitux.com/img/twiph_100_younger_avatar.png
.. _image1: http://www.eraclitux.com/img/twiph_100_younger_avatar.png

Live demo
=========

http://www.eraclitux.com/twiph_demo/100_younger_camera_groups/index_groups.html

http://www.eraclitux.com/twiph_demo/100_younger_camera_groups/index_avatar.html

Usage
=====

Once retrieved auth credentils and created a ``cfg`` file (a sample is provided in ``conf/sample.cfg``)::

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

Political data has been retrieved by ``Openpolis`` REST APIs http://api3.openpolis.it/

Disclaimer
==========

All trademarks, copyrights and other forms of intellectual property belong to their respective owners.

The author is not affiliated with any vendor cited above.
