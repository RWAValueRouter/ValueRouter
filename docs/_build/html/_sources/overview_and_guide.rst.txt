System Requirements and Installation Guide
########################################

.. toctree::
  :maxdepth: 2

Installation
============
Following is how we build fastmpc service in ubuntu.

Ubuntu
******

Install service in linux ::

    $ install golang 1.17 or above
    $ install mysql 8.0
    $ import docs/sql/* into database
    $ make
    $ nohup ./smw --rpcport 8888 &

add config.json following content::

    {
      "DbConfig": {
        "DbDriverName": "mysql",
        "DbDriverSource": "root:12345678@tcp(127.0.0.1:3306)/smw"
      }
    }
