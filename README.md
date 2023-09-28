## Prerequisites

1. Akamai API credentials

## Usage
You can send the `--help` flag to check functionality, following are the most use commands.

#### Regenerating SQS messages

1. Execute the `with_role` script depending on the account you need to regenerate messages
2. Run the script
    ```bash
    go run . sqs create -f {filename} -q {queue_name}
    ```

#### Pulling akamai bw reports

1. Export all the environment variables
    ```bash
    export AKAMAI_HOST=foo
    export AKAMAI_CLIENT_TOKEN=foo
    export AKAMAI_CLIENT_SECRET=foo
    export AKAMAI_ACCESS_TOKEN=foo
    ```
2. Run the script
    ```bash
    go run . reports create -s {start_date} -e {end_date}
    ```