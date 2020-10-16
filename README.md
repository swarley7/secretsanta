# secretsanta
Randomly assigns participants within a Secret Santa system.

## Features

* Randomises matches for secret santa
* supports multiple gifts (e.g. "Every participant has to buy a $100 and a $20 gift") with either unique or repeating participants (e.g. you don't want John to buy more than one gift for Jane).
* SNS for delivery of gift registry to participants

## Usage:

`git clone https://github.com/swarley7/secretsanta.git`

`cd secretsanta && go get -u # install deps`

Use your favourite text editor to copy / create a `config.json` file (an example is provided)

Ensure your AWS env is setup correctly:

`brew install awscli`

`aws configure`

(you'll need some IAM access keys and stuff - follow the AWS guide for it).

Once done, you can run your secret santa by doing

`AWS_REGION=<YOUR PREFERRED AWS REGION HERE> go run secretsanta.go -d=true -c <YOUR CONFIG FILE HERE>`

*Note:* Debug mode disables the actual delivery of SMS' for safety. Use it to validate your secret santa app works before doing it live! (A live run will have different output :) )

When you're ready to do it live, flip the debug flag to false:

`AWS_REGION=<YOUR PREFERRED AWS REGION HERE> go run secretsanta.go -d=false -c <YOUR CONFIG FILE HERE>`
