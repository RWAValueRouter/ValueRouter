Getting Started
########################################

.. toctree::
:maxdepth: 3

Json RPC
======================
fast mpc service middleware api details

.. api:


Get group id
-----------------------------------------

To create a mpc address. the first step is to get group id:

- the first param ``"2/2"``: indicate the mpc address we are creating is 2 out of 2 threshold address
- the second array parameter ``["0xac9526c5db81267804a32ec508db780402fd9fec","0xd17831dd9db4ce9a8d331c807329e93015ca2bcb"]`` indicate which addresses are used to create mpc address. we also can add specific node for your designated address using struct like ``["0xac9526c5db81267804a32ec508db780402fd9fec|127.0.0.1:8383","0xd17831dd9db4ce9a8d331c807329e93015ca2bcb"]``

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc":"2.0",
          "method":"smw_getGroupId",
          "params":[
              "2/2",
              [
                  "0xac9526c5db81267804a32ec508db780402fd9fec",
                  "0xd17831dd9db4ce9a8d331c807329e93015ca2bcb"
              ]
          ],
          "id":67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc":"2.0",
          "id":67,
          "result":{
              "Data":{
                  "Gid":"2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                  "Sigs":"2:748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b:0xac9526c5db81267804a32ec508db780402fd9fec:2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96:0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                  "Uuid":"466c90a9-b9ee-42db-97d0-93f85bf40203"
              },
              "Error":"",
              "Status":"success",
              "Tip":""
          }
      }

   - if the response status is ``success`` , indicate the request is success , if the value is ``error`` then indicate something unexpected happened
   - the response value ``Gid, Sigs, Uuid``, you need to store it somewhere , it will be need as request param to create mpc address.


Mpc Address key generating
-----------------------------------------

Key generating used to create mpc address :

- the first param ``"0x3d689c4...f011079c01"`` is the signature of your metamask account signing the second parameter.
- the second param ``"TxType"`` must be exact `REQSMPCADDR`, ``"Account"`` metamask user account, ``"Nonce"`` metamask user account current nonce, ``"Keytype"`` mpc address creating algorithm can be ``EC256K1 or EC256K1`` ``"GroupID"`` gid return by getGroupId, ``"ThresHold"`` threshold can be 2/3, 3/3 , 3/4 etc. ``"Mode"`` must be 2 as it is specified to be used in wallet service, ``"FixedApprover"`` leave it as null ``"AcceptTimeOut"`` request smpc address must be accepted timeout as we are using mode 2. leave it as it is. ``"TimeStamp"`` api invoke timestamp, ``"Sigs"`` return by getGroupId,  ``"Comment"`` comments , ``"Uuid"`` return by getGroupId

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc":"2.0",
          "method":"smw_keyGen",
          "params":[
              "0x3d689c4805f38ca1a12b2e975941c31a3b5aebcff3ac76dd466304da9635fd246bcf387095187e400662bfdfbc303b3439ba0b300d49572f97706e4af011079c01",
              "{\"TxType\":\"REQSMPCADDR\",\"Account\":\"0xac9526C5db81267804a32eC508dB780402fD9fEC\",\"Nonce\":\"3\",\"Keytype\":\"EC256K1\",\"GroupID\":\"2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73\",\"ThresHold\":\"2/2\",\"Mode\":\"2\",\"FixedApprover\":null,\"AcceptTimeOut\":\"604800\",\"TimeStamp\":\"1678095190605845000\",\"Sigs\":\"2:748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b:0xac9526c5db81267804a32ec508db780402fd9fec:2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96:0xd17831dd9db4ce9a8d331c807329e93015ca2bcb\",\"Comment\":\"\",\"Uuid\":\"466c90a9-b9ee-42db-97d0-93f85bf40203\"}"
          ],
          "id":67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": "0x373e1e7d7f84160ee37bf06f07233b7f15090efb81202a9aa467da99f67edd1b",
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response ``Data`` is the ``key id`` of this mpc address. can be used to check mpc address creating status.


Check the creation of mpc address status
-----------------------------------------

To check if a mpc address is already generated:

- the first param ``"0x373e1e7...f67edd1b`` is the returned key id from keygen

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc":"2.0",
          "method":"smw_getReqAddrStatus",
          "params":[
              "0x373e1e7d7f84160ee37bf06f07233b7f15090efb81202a9aa467da99f67edd1b"
          ],
          "id":67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": [
                  {
                      "Status": "1",
                      "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                      "Key_id": "0x373e1e7d7f84160ee37bf06f07233b7f15090efb81202a9aa467da99f67edd1b",
                      "Public_key": "0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                      "Mpc_address": "0x1166239261e699d8cA06dFB0031645c63e5AB528",
                      "Initializer": "1",
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678095190605",
                      "Reply_enode": "748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b",
                      "Gid": "2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                      "Threshold": "2/2"
                  },
                  {
                      "Status": "1",
                      "User_account": "0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                      "Key_id": "0x373e1e7d7f84160ee37bf06f07233b7f15090efb81202a9aa467da99f67edd1b",
                      "Public_key": "0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                      "Mpc_address": "0x1166239261e699d8cA06dFB0031645c63e5AB528",
                      "Initializer": "0",
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678095190605",
                      "Reply_enode": "2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96",
                      "Gid": "2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                      "Threshold": "2/2"
                  }
              ],
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value ``"Status"``: 0 pending , 1 success , 2 fail, 3 timeout, ``"User_account"``: user metamask account, ``"Key_id"``: key id of this mpc address, ``"Public_key"``: mpc address public key, ``"Mpc_address"``: generated mpc address, ``"Initializer"``: 1 initiator 0 not initiator, ``"Reply_status"``: `AGREE` or `DISAGREE` reply status of generating mpc address , ``"Reply_timestamp"``: reply timestamp, ``"Reply_enode"``: reply node, ``"Gid"``: group id of creating this mpc address, ``"Threshold"``: threshold of this mpc address. if ``"status"`` is not 1 then some of these values may be empty.


Get account list
-----------------------------------------

Get mpc account list of metamask user account created:

- the first param ``"0xac9526c...02fd9fec`` is the metamask account

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
        "jsonrpc":"2.0",
        "method":"smw_getAccountList",
        "params":[
            "0xac9526c5db81267804a32ec508db780402fd9fec"
        ],
        "id":67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "jsonrpc": "2.0",
        "id": 67,
        "result": {
            "Data": [
                {
                    "Status": "1",
                    "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                    "Key_id": "0x772087ef9db3f3fb462b047e654f176d37f3f447b55acc33a01ff3c221976ea9",
                    "Public_key": "04264725f78fe99303ffe08e7468be5c176fe3d65520888b36b970ba02d8c3becdcf9c047616f3c2296f9d33a22102b251b0278e6d39e8b70f38ca5c3228b3972e",
                    "Mpc_address": "0x26426df61e1B587a2Dc8e6f3222449072f4Cf8aB",
                    "Initializer": "1",
                    "Reply_status": "AGREE",
                    "Reply_timestamp": "1678074799443",
                    "Reply_enode": "2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96",
                    "Gid": "2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                    "Threshold": "2/2"
                },
                {
                    "Status": "1",
                    "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                    "Key_id": "0x373e1e7d7f84160ee37bf06f07233b7f15090efb81202a9aa467da99f67edd1b",
                    "Public_key": "0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                    "Mpc_address": "0x1166239261e699d8cA06dFB0031645c63e5AB528",
                    "Initializer": "1",
                    "Reply_status": "AGREE",
                    "Reply_timestamp": "1678095190605",
                    "Reply_enode": "748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b",
                    "Gid": "2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                    "Threshold": "2/2"
                }
            ],
            "Error": "",
            "Status": "success",
            "Tip": ""
        }
      }


Get unsigned transaction hash
-----------------------------------------

To get unsigned tx hash:

- the first param ``"from"``: sender address should be mpc address, ``"to"``: receiver address, ``"chainId"``: chain id of evm, ``"value"``: transfer value , ``"nonce"``: mpc address nonce, ``"gas"``: gas limit of this transaction, ``"gasPrice"``: gas price of this transaction, ``"data"``: transaction payload data, ``"originValue"``: transfer value in `eth`, ``"name"``: chain name
- the second param ``0`` means chain type is evm.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
        "jsonrpc":"2.0",
        "method":"smw_getUnsigedTransactionHash",
        "params":[
            "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x6\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}",
            0
        ],
        "id":67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "jsonrpc": "2.0",
        "id": 67,
        "result": {
            "Data": "0x8ed2fca05f6f17a96440dcfe9aa3edebb236737a7817c298c485e12099fac886",
            "Error": "",
            "Status": "success",
            "Tip": ""
        }
      }

   - the response value ``"0x8ed2f...9fac886"``: is the unsigned tx hash. will be used to in ``"sign"``


Sign transaction
-----------------------------------------

Sign transaction:

- the first param ``"0x7a0d1338...3563f400"`` is the signature of your metamask account signing the second parameter.
- the second param ``"TxType"`` must be `SIGN`, ``"Account"`` metamask user account, ``"Nonce"`` metamask user account current nonce, ``"PubKey"`` public key of of your mpc address, ``"InputCode"`` leave it blank, ``"MsgHash"`` msg hash get from `smw_getUnsigedTransactionHash` or also can use `keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))` get hash , ``"MsgContext"`` request param of `smw_getUnsigedTransactionHash` or the `message` , `MsgHash` and `MsgContext` must match. ``"ChainType"`` 0 indicate evm chain. other field is the same as mentioned above.
- you need to ``Get unsigned transaction hash`` first before using this api

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
        "jsonrpc": "2.0",
        "method": "smw_sign",
        "params": [
            "0x7a0d13384b9038a814424263d4fb63f7c0577a2783e0457c9a59d0213e488f03392a311ddae9f4f2f1682c2bb91b0d2d387a50f91c8441267ec6133a053563f400",
            "{\"TxType\":\"SIGN\",\"Account\":\"0xac9526C5db81267804a32eC508dB780402fD9fEC\",\"Nonce\":\"5\",\"PubKey\":\"0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b\",\"InputCode\":\"\",\"MsgHash\":[\"0x8ed2fca05f6f17a96440dcfe9aa3edebb236737a7817c298c485e12099fac886\"],\"MsgContext\":[\"{\\\"from\\\":\\\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\\\",\\\"to\\\":\\\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\\\",\\\"chainId\\\":\\\"0x6\\\",\\\"value\\\":\\\"0xb5e620f48000\\\",\\\"nonce\\\":0,\\\"gas\\\":47984,\\\"gasPrice\\\":5338199,\\\"data\\\":\\\"\\\",\\\"originValue\\\":\\\"0.0002\\\",\\\"name\\\":\\\"ETH Goerli\\\"}\"],\"Keytype\":\"EC256K1\",\"GroupID\":\"2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73\",\"ThresHold\":\"2/2\",\"Mode\":\"2\",\"AcceptTimeOut\":\"604800\",\"TimeStamp\":\"1678414610826\",\"FixedApprover\":null,\"Comment\":\"\",\"ChainType\":0}"
        ],
        "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "jsonrpc": "2.0",
        "id": 67,
        "result": {
            "Data": "0x56370d0d25ef367d992a785a280cf020560ffdc22fb4d6c7e68e3f2dbfb78b41",
            "Error": "",
            "Status": "success",
            "Tip": ""
        }
      }

   - the response value ``"0x56370d...fb78b41"``: is the key id of signed tx, will be used to check the transaction status, signing status etc.

Get approval list
-----------------------------------------

Get approval list of pending transaction that need to be agree or disagree:

- the first param ``"0xD17831dd9db4Ce9a8d331c807329e93015cA2BcB"`` is the metamask user account.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
        "jsonrpc": "2.0",
        "method": "smw_getApprovalList",
        "params": [
            "0xD17831dd9db4Ce9a8d331c807329e93015cA2BcB"
        ],
        "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "jsonrpc": "2.0",
        "id": 67,
        "result": {
            "Data": [
                {
                    "User_account": "0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                    "Group_id": "2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                    "Key_id": "0x56370d0d25ef367d992a785a280cf020560ffdc22fb4d6c7e68e3f2dbfb78b41",
                    "Key_type": "EC256K1",
                    "Mode": "2",
                    "Msg_context": [
                        "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x6\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}"
                    ],
                    "Msg_hash": [
                        "0x8ed2fca05f6f17a96440dcfe9aa3edebb236737a7817c298c485e12099fac886"
                    ],
                    "Nonce": "5",
                    "Public_key": "0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                    "Mpc_address": "0x1166239261e699d8cA06dFB0031645c63e5AB528",
                    "Threshold": "2/2",
                    "TimeStamp": "1678095190605",
                    "Status": 1,
                    "Signed": 0,
                    "Chain_id": 97,
                    "Chain_type": 0
                },
                {
                    "User_account": "0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                    "Group_id": "2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                    "Key_id": "0xcdc1f5ffeadd177bd2db0e14520d17288ce6a72c11990da7629f2e9760e46427",
                    "Key_type": "EC256K1",
                    "Mode": "2",
                    "Msg_context": [
                        "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x5\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}",
                        "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x5\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli111\"}"
                    ],
                    "Msg_hash": [
                        "0x57f42fb3b0f419e6ef0c9805ec7987d955c837814cbd16e11ef1fec12777050d",
                        "0x57f42fb3b0f419e6ef0c9805ec7987d955c837814cbd16e11ef1fec12777050e"
                    ],
                    "Nonce": "4",
                    "Public_key": "0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                    "Mpc_address": "0x1166239261e699d8cA06dFB0031645c63e5AB528",
                    "Threshold": "2/2",
                    "TimeStamp": "1678095190605",
                    "Status": 1,
                    "Signed": 0,
                    "Chain_id": 97,
                    "Chain_type": 0
                }
            ],
            "Error": "",
            "Status": "success",
            "Tip": ""
        }
      }

   - the response value ``"Status"``: `1` means already handled , `0` means need to agree or disagree this transaction.
   - the latest 100 record will returned at most

Paginated approval list
-----------------------------------------

Get approval list of specific page and page size:

- the first param ``"0xac9526c5db81267804a32ec508db780402fd9fec"`` is the metamask user account.
- the second param is the ``"status"`` of returned data it can only be `0` (to be approval) or `1` (already approval or timeout etc).
- the third param is the page number
- the fourth param is the page size

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "method": "smw_getApprovalListByPagination",
          "params": [
              "0xac9526c5db81267804a32ec508db780402fd9fec",
              0,
              1,
              1
          ]
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "result": {
              "Data": [
                  {
                      "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                      "Group_id": "1f29316aef29dad0f2b1a4a5a53943318a65df56116a836200f8ed45f70c2f1b2965e1ac53d2c938be530af74c6f19e60bd612524cdc14264af48f7776c87aeb",
                      "Key_id": "0xf7f6ce06d0faaedbe5f7b6491041719dfd023077e9bba974f049eb5993796ef5",
                      "Key_type": "EC256K1",
                      "Mode": "2",
                      "Msg_context": [
                          "{\"from\":\"0xDb5c6ec1a096EB4cDB7ab52718abE332a95A6f8d\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x5\",\"value\":\"0x110d9316ec000\",\"nonce\":0,\"gas\":48016,\"gasPrice\":67727197933,\"data\":\"\",\"originValue\":\"0.0003\",\"name\":\"ETH Goerli\"}"
                      ],
                      "Msg_hash": [
                          "0x00f4dba42d61491ed98be5224dfa3bfed7c713aae6f2eb909f7c66ad7ee8bf90"
                      ],
                      "Nonce": "5",
                      "Public_key": "04033dcd0505bc2b78479a102b0940c0ee2cf2a3080c6a19e88e685396af7fe0e977f506c89aa475c21c758ea5b611d91f06ee27e20a12e77e17c6ff55b230cfaa",
                      "Mpc_address": "0xDb5c6ec1a096EB4cDB7ab52718abE332a95A6f8d",
                      "Threshold": "2/2",
                      "Timestamp": "1678763889487",
                      "Status": 0,
                      "Signed": 0,
                      "Chain_id": 97,
                      "Chain_type": 0

                  }
              ],
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the data is using descending order of create time, so the latest are the first to grab


Get approval list through key id
-----------------------------------------

get approval list through signed key id:

- the first param ``"0xce8e22a9....6f4498678"`` is the signed key id

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_getApprovalListByKeyId",
          "params": [
              "0xce8e22a90bc7c79b230e29cb850e70e5e3c42d426aa628d5ce45c1c6f4498678"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": [
                  {
                      "User_account": "0x8d5992c1151439c0a6f421564588f556722eece6",
                      "Group_id": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Key_id": "0xce8e22a90bc7c79b230e29cb850e70e5e3c42d426aa628d5ce45c1c6f4498678",
                      "Key_type": "EC256K1",
                      "Mode": "2",
                      "Msg_context": [
                          "{\"from\":\"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x61\",\"value\":\"0x1C6BF52634000\",\"nonce\":0,\"gas\":100000,\"gasPrice\":10000000000,\"data\":\"\",\"originValue\":\"0.0005\",\"name\":\"BSC testnet\"}"
                      ],
                      "Msg_hash": [
                          "0x95c58a2bf0cff8f69c7bae7c1b1f834de1da132022a31befbbbeba8a45733314"
                      ],
                      "Nonce": "5",
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Threshold": "2/3",
                      "Timestamp": "1679628957289",
                      "Status": 1,
                      "Signed": 3,
                      "Chain_id": 97,
                      "Chain_type": 0,
                      "Reply_status": "TIMEOUT",
                      "Reply_initializer": 0
                  },
                  {
                      "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                      "Group_id": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Key_id": "0xce8e22a90bc7c79b230e29cb850e70e5e3c42d426aa628d5ce45c1c6f4498678",
                      "Key_type": "EC256K1",
                      "Mode": "2",
                      "Msg_context": [
                          "{\"from\":\"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x61\",\"value\":\"0x1C6BF52634000\",\"nonce\":0,\"gas\":100000,\"gasPrice\":10000000000,\"data\":\"\",\"originValue\":\"0.0005\",\"name\":\"BSC testnet\"}"
                      ],
                      "Msg_hash": [
                          "0x95c58a2bf0cff8f69c7bae7c1b1f834de1da132022a31befbbbeba8a45733314"
                      ],
                      "Nonce": "5",
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Threshold": "2/3",
                      "Timestamp": "1679628957289",
                      "Status": 1,
                      "Signed": 3,
                      "Chain_id": 97,
                      "Chain_type": 0,
                      "Reply_status": "AGREE",
                      "Reply_initializer": 1
                  },
                  {
                      "User_account": "0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                      "Group_id": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Key_id": "0xce8e22a90bc7c79b230e29cb850e70e5e3c42d426aa628d5ce45c1c6f4498678",
                      "Key_type": "EC256K1",
                      "Mode": "2",
                      "Msg_context": [
                          "{\"from\":\"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x61\",\"value\":\"0x1C6BF52634000\",\"nonce\":0,\"gas\":100000,\"gasPrice\":10000000000,\"data\":\"\",\"originValue\":\"0.0005\",\"name\":\"BSC testnet\"}"
                      ],
                      "Msg_hash": [
                          "0x95c58a2bf0cff8f69c7bae7c1b1f834de1da132022a31befbbbeba8a45733314"
                      ],
                      "Nonce": "5",
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Threshold": "2/3",
                      "Timestamp": "1679628957289",
                      "Status": 1,
                      "Signed": 3,
                      "Chain_id": 97,
                      "Chain_type": 0,
                      "Reply_status": "AGREE",
                      "Reply_initializer": 0
                  }
              ],
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value is a list of approval account approving details.
   - ``Status`` 0: not handled, 1: handled. ``Signed`` the number of account has handled this transaction. ``Reply_status`` account reply. ``Reply_initializer`` 0: not initiator, 1:initiator


Accept transaction
-----------------------------------------

Accept a signed transaction:

- the first param ``"0xf3a41504...2f1abcc00"`` is the signature of your metamask account signing the second parameter.
- the second param ``"TxType"`` must be `ACCEPTSIGN`, ``"Account"`` metamask user account, ``"Nonce"`` metamask user account current nonce, ``"Key"`` key id of signed tx, ``"Accept"`` must be `AGREE` or `DISAGREE`, ``"MsgHash"`` the same with `smw_sign`, ``"MsgContext"`` the same with `smw_sign`, `MsgHash` and `MsgContext` must be the same with the same key id of signed tx. other field is the same as mentioned above.
- you need to ``get approval list`` first before accept transaction

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
        "jsonrpc":"2.0",
        "method":"smw_acceptSign",
        "params":[
            "0xf3a41504407af0fff3c54ca1ae383a15faa16e2dd1add76a6e10e052b8160a771995555962dbab8df09bd146c2ff25fcba73b97f0dc09033b9d920f852f1abcc00",
            "{\"TxType\":\"ACCEPTSIGN\",\"Account\":\"0xD17831dd9db4Ce9a8d331c807329e93015cA2BcB\",\"Nonce\":\"5\",\"Key\":\"0x56370d0d25ef367d992a785a280cf020560ffdc22fb4d6c7e68e3f2dbfb78b41\",\"Accept\":\"AGREE\",\"MsgHash\":[\"0x8ed2fca05f6f17a96440dcfe9aa3edebb236737a7817c298c485e12099fac886\"],\"MsgContext\":[\"{\\\"from\\\":\\\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\\\",\\\"to\\\":\\\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\\\",\\\"chainId\\\":\\\"0x6\\\",\\\"value\\\":\\\"0xb5e620f48000\\\",\\\"nonce\\\":0,\\\"gas\\\":47984,\\\"gasPrice\\\":5338199,\\\"data\\\":\\\"\\\",\\\"originValue\\\":\\\"0.0002\\\",\\\"name\\\":\\\"ETH Goerli\\\"}\"],\"TimeStamp\":\"1678416408863\",\"ChainType\":0}"
        ],
        "id":67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "jsonrpc": "2.0",
        "id": 67,
        "result": {
            "Data": "Success",
            "Error": "",
            "Status": "success",
            "Tip": ""
        }
      }

   - the response value ``"Success"``: indicate request is ok.

Sign history
-----------------------------------------

Get Sign history of metamask user account:

- the first param ``"0xD17831dd9db4Ce9a8d331c807329e93015cA2BcB"`` is metamask user account.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_getSignHistory",
          "params": [
              "0xD17831dd9db4Ce9a8d331c807329e93015cA2BcB"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc":"2.0",
          "id":67,
          "result":{
              "Data":[
                  {
                      "User_account":"0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                      "Group_id":"2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                      "Key_id":"0x56370d0d25ef367d992a785a280cf020560ffdc22fb4d6c7e68e3f2dbfb78b41",
                      "Key_type":"EC256K1",
                      "Mode":"2",
                      "Msg_context":[
                          "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x6\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}"
                      ],
                      "Msg_hash":[
                          "0x8ed2fca05f6f17a96440dcfe9aa3edebb236737a7817c298c485e12099fac886"
                      ],
                      "Public_key":"0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                      "Mpc_address":"0x1166239261e699d8cA06dFB0031645c63e5AB528",
                      "Threshold":"2/2",
                      "Txid":"",
                      "Status":1,
                      "Reply_status":"AGREE",
                      "Reply_timestamp":"1678416408863",
                      "Signed": 3,
                      "Local_timestamp": "1680839066000",
                      "Chain_id": 97,
                      "Chain_type": 0
                  },
                  {
                      "User_account":"0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                      "Group_id":"2bc9ac6c25f2e47fa1f0d2f6968d19b13c261179f4b783414ac86e9a2db6501f71eb20114c759dfb58a2708dc2b180afe4d44d4a2d88b53e312d97cbc3134b73",
                      "Key_id":"0xcdc1f5ffeadd177bd2db0e14520d17288ce6a72c11990da7629f2e9760e46427",
                      "Key_type":"EC256K1",
                      "Mode":"2",
                      "Msg_context":[
                          "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x5\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}",
                          "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x5\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli111\"}"
                      ],
                      "Msg_hash":[
                          "0x57f42fb3b0f419e6ef0c9805ec7987d955c837814cbd16e11ef1fec12777050d",
                          "0x57f42fb3b0f419e6ef0c9805ec7987d955c837814cbd16e11ef1fec12777050e"
                      ],
                      "Public_key":"0457404abfd62efff9bfcf1ee0aaa74ffc1355b494e8294f2f3309f52f8376a92d778b5f7b7c14cf696284f1d565a6b7bae031fc25c00b3dfca1d65e4c8da4af0b",
                      "Mpc_address":"0x1166239261e699d8cA06dFB0031645c63e5AB528",
                      "Threshold":"2/2",
                      "Txid":"",
                      "Status":1,
                      "Reply_status":"AGREE",
                      "Reply_timestamp":"1678331951040"
                      "Signed": 3,
                      "Local_timestamp": "1680839066000",
                      "Chain_id": 97,
                      "Chain_type": 0
                  }
              ],
              "Error":"",
              "Status":"success",
              "Tip":""
          }
      }

   - the response value is the transaction history list.
   - the latest 100 record can be returned at most


Paginated sign history
-----------------------------------------

Get Sign history of specific page and pageSize :

- the first param ``"0xD17831dd9db4Ce9a8d331c807329e93015cA2BcB"`` is metamask user account.
- the second param is page , the third param is the page size

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "method": "smw_getSignHistoryByPagination",
          "params": [
              "0xac9526c5db81267804a32ec508db780402fd9fec",
              1,
              2
          ]
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "result": {
              "Data": [
                  {
                      "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                      "Group_id": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Key_id": "0x19b613c480776317fd321761ec06ded42797b6edfae94c0e1ac2de6e07164e6f",
                      "Key_type": "EC256K1",
                      "Mode": "2",
                      "Msg_context": [
                          "{\"from\":\"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x61\",\"value\":\"0x1C6BF52634000\",\"nonce\":0,\"gas\":100000,\"gasPrice\":10000000000,\"data\":\"\",\"originValue\":\"0.0005\",\"name\":\"bsc testnet\"}"
                      ],
                      "Msg_hash": [
                          "0x95c58a2bf0cff8f69c7bae7c1b1f834de1da132022a31befbbbeba8a45733314"
                      ],
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Threshold": "2/3",
                      "Txid": "",
                      "Status": 7,
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678948494520",
                      "Signed": 3,
                      "Local_timestamp": "1680839066000",
                      "Chain_id": 97,
                      "Chain_type": 0
                  },
                  {
                      "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                      "Group_id": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Key_id": "0xce8e22a90bc7c79b230e29cb850e70e5e3c42d426aa628d5ce45c1c6f4498678",
                      "Key_type": "EC256K1",
                      "Mode": "2",
                      "Msg_context": [
                          "{\"from\":\"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F\",\"to\":\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\",\"chainId\":\"0x61\",\"value\":\"0x1C6BF52634000\",\"nonce\":0,\"gas\":100000,\"gasPrice\":10000000000,\"data\":\"\",\"originValue\":\"0.0005\",\"name\":\"BSC testnet\"}"
                      ],
                      "Msg_hash": [
                          "0x95c58a2bf0cff8f69c7bae7c1b1f834de1da132022a31befbbbeba8a45733314"
                      ],
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Threshold": "2/3",
                      "Txid": "0x1269c0e1ac1344db7493e74de99c4335c882e60b0bf5adc3cd34b3fa33b25e87",
                      "Status": 5,
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678948494520",
                      "Signed": 3,
                      "Local_timestamp": "1679628993000",
                      "Chain_id": 97,
                      "Chain_type": 0
                  }
              ],
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the data is using descending order of create time, so the latest are the first to grab


Get transaction hash
-----------------------------------------

Get transaction hash through sign key id:

- the first param ``"0xaf9a7bcfd...530d6da877e8"`` is `sign` returned key id.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_getTxHashByKeyId",
          "params": [
              "0xaf9a7bcfd420e0c618f25eafb27ff80c94d0b8439115f42698e5530d6da877e8"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": "0x33a8737c32a4cf65f4a2469027b92b2cb1682588484fd164de0d010f59bf72c7",
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value `Data` is the transaction hash which has been sent on chain.
   - if the transaction has not been sent on chain then the returned value should be empty.


Get signed transaction status
-----------------------------------------

Get signed transaction status through signed key id:

- the first param ``"0xaf9a7bcfd...530d6da877e8"`` is `sign` returned key id.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_getTxStatusByKeyId",
          "params": [
              "0xaf9a7bcfd420e0c618f25eafb27ff80c94d0b8439115f42698e5530d6da877e8"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": 4,
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value `Data` ``4`` means the transaction has been sent on chain.there are some other value. which are 0: ``"Mpc-Pending"`` , 1: ``"Mpc-Success"`` , 2: ``"MPC-Fail"``, 3: ``"MPC-Timeout"``, 4: ``"Tx-Pending"``, 5: ``"Tx-Confirmed"``, 6: ``"Tx-Failed"``, 7: ``"Tx-NotValid"``, 8: ``"Tx-Replaced"``


Add asset
-----------------------------------------

Add asset of specific metamask user account:

- the first param ``"0x21fb6...0e033b501"`` is the signature of your metamask account signing the second parameter.
- the second param ``"TxType"`` must be `ADDASSET`, ``"Account"`` metamask user account, ``"Nonce"`` metamask user account current nonce, ``"Symbol"`` asset symbol, ``"Name"`` asset name, ``"Decimal"`` asset decimal, ``"Contract"`` asset contract, ``"ChainId"`` chain id of evm or 0 for none evm, ``"ChainType"`` 0 for evm chain,

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_addAsset",
          "params": [
              "0x21fb6c871f4f93d52eb1adc90f485b573cfe5f2075c964dc4eccd749a20955ab2abd7a289d5dd1f0bc87e11f0f18c922bd0cf9c84160cf27a242564620e033b501",
              "{\"TxType\":\"ADDASSET\",\"Account\":\"0xac9526C5db81267804a32eC508dB780402fD9fEC\",\"Nonce\":\"5\",\"Symbol\":\"BBD\",\"Name\":\"BBD Coin\",\"Decimal\":18,\"Contract\":\"0x168f6dec26cbbb3749654e0e3cc4fc29314fdf6d\",\"TimeStamp\":\"1679906687499\",\"ChainId\":100,\"ChainType\":0}"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": "success",
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }


Add asset for mpc address
-----------------------------------------

Add asset for specific mpc address, only owner can add asset:

- compare with ``"AddAsset"`` there is one more field need to be add, which is the mpc address
- the first param ``"0x21fb6...0e033b501"`` is the signature of your metamask account signing the second parameter.
- the second param ``"TxType"`` must be `ADDASSETFORMPCADDRESS`, ``"Account"`` metamask user account, ``"Nonce"`` metamask user account current nonce, ``"Symbol"`` asset symbol, ``"Name"`` asset name, ``"Decimal"`` asset decimal, ``"Contract"`` asset contract, ``"ChainId"`` chain id of evm or 0 for none evm, ``"ChainType"`` 0 for evm chain, ``"MpcAddress"`` the mpc address which add this address.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_addAssetForMpcAddress",
          "params": [
              "0x21fb6c871f4f93d52eb1adc90f485b573cfe5f2075c964dc4eccd749a20955ab2abd7a289d5dd1f0bc87e11f0f18c922bd0cf9c84160cf27a242564620e033b501",
              "{\"TxType\":\"ADDASSETFORMPCADDRESS\",\"Account\":\"0xac9526C5db81267804a32eC508dB780402fD9fEC\",\"Nonce\":\"5\",\"Symbol\":\"BBD\",\"Name\":\"BBD Coin\",\"Decimal\":18,\"Contract\":\"0x168f6dec26cbbb3749654e0e3cc4fc29314fdf6d\",\"TimeStamp\":\"1679906687499\",\"ChainId\":100,\"ChainType\":0,\"MpcAddress\":\"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F\"}"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": "success",
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

Get asset
-----------------------------------------

Get asset of specific metamask user account or mpc address:

- the first param ``"0xac9526C5db81267804a32eC508dB780402fD9fEC"`` is the metamask account or mpc address
- the second param is the chain id ,if non-evm then the value should be 0
- the third param is the chain type , 0 indicate EVM chains

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_getAsset",
          "params": [
              "0xac9526C5db81267804a32eC508dB780402fD9fEC",
              100,
              0
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": [
                  {
                      "Symbol": "bbc",
                      "Contract": "0x168f6dec26cbbb3749654e0e3cc4fc29314fdf6d",
                      "Name": "bbc coin",
                      "Decimal": 18
                  },
                  {
                      "Symbol": "bbd",
                      "Contract": "0x168f6dec26cbbb3749654e0e3cc4fc29314fdf6d",
                      "Name": "bbd coin",
                      "Decimal": 18
                  }
              ],
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value is a list of asset that added by this account

Get mpc address details
-----------------------------------------

Get mpc address contained account details:

- the first param ``"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F"`` is the mpc address

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "method": "smw_getMpcAddressDetail",
          "params": [
              "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F"
          ],
          "id": 67
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 67,
          "result": {
              "Data": [
                  {
                      "Status": "1",
                      "User_account": "0xac9526c5db81267804a32ec508db780402fd9fec",
                      "Key_id": "0x5f3cd990d9dbc1a2bc5f217829246cb24aca3eb571fea0dd171bb2fe32f62fde",
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Initializer": "1",
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678948494520",
                      "Reply_enode": "748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b",
                      "Gid": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Threshold": "2/3",
                      "Mode": "2",
                      "Key_type": "EC256K1"
                  },
                  {
                      "Status": "1",
                      "User_account": "0xd17831dd9db4ce9a8d331c807329e93015ca2bcb",
                      "Key_id": "0x5f3cd990d9dbc1a2bc5f217829246cb24aca3eb571fea0dd171bb2fe32f62fde",
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Initializer": "0",
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678948494520",
                      "Reply_enode": "2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96",
                      "Gid": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Threshold": "2/3",
                      "Mode": "2",
                      "Key_type": "EC256K1"
                  },
                  {
                      "Status": "1",
                      "User_account": "0x8d5992c1151439c0a6f421564588f556722eece6",
                      "Key_id": "0x5f3cd990d9dbc1a2bc5f217829246cb24aca3eb571fea0dd171bb2fe32f62fde",
                      "Public_key": "04f0f6b326aea6a3eaa6e2603fcadb82131be6afa176145a62b7fbe649a9142eed6b34679b1e6bb1f29128649caadf37b7d7ff6132e2d407dff980ffa93492529e",
                      "Mpc_address": "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
                      "Initializer": "0",
                      "Reply_status": "AGREE",
                      "Reply_timestamp": "1678948494520",
                      "Reply_enode": "08ba43b0715bb27e03911592d3fed49a22f49ecaf1628d44c4c7d4a8914b86d423716d7316d98a80f4d280b852a33ba7972f0a744e386375e6d58a12fab96752",
                      "Gid": "874a5f6f4aca0ea16487be4269476c5231df8ab59f91340c723b59e487877a383fff449f351e6fe28543a4c109b91a46f02e782e01640c9c54a74f1ec1a50824",
                      "Threshold": "2/3",
                      "Mode": "2",
                      "Key_type": "EC256K1"
                  }
              ],
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value is a list of account that this mpc address contains.


Get latest mpc address working status
-----------------------------------------

get the latest mpc address status:

- the first param ``"0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F"`` is the mpc address
- the second param the chain_id
- the second param the chain_type, 0 for evm.

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "method": "smw_getLatestMpcAddressStatus",
          "params": [
              "0xFDE75431218d8a16E20B43e333Cb45D57Ce70D7F",
              97,
              0
          ]
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "result": {
              "Data": 5,
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value data ``"5"`` means transaction is confirmed on chain.  means the transaction has been sent on chain.there are some other value. which are -1: ``"no record find"``, 0: ``"Mpc-Pending"`` , 1: ``"Mpc-Success"`` , 2: ``"MPC-Fail"``, 3: ``"MPC-Timeout"``, 4: ``"Tx-Pending"``, 5: ``"Tx-Confirmed"``, 6: ``"Tx-Failed"``, 7: ``"Tx-NotValid"``, 8: ``"Tx-Replaced"`` . ``"0,1,4"`` means mpc address is currently working.


Get node number
-----------------------------------------

get the number of node on server:

   **Example request**:

   .. sourcecode:: http

    POST http://localhost:8888 HTTP/1.1
    context-type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "method": "smw_getNodesNumber",
          "params": []
      }

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "jsonrpc": "2.0",
          "id": 97,
          "result": {
              "Data": 3,
              "Error": "",
              "Status": "success",
              "Tip": ""
          }
      }

   - the response value is the number of nodes that server used in service.
