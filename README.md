
##### Current version: UNRELEASED / FIRST VERSION IN DEV

Cloud-native auto-scaling load balancer with built-in DNS zone updates.

Gossip cluster powers:
* Dynamic scaling 
* Shared state
* Fast fail-over of load-balancer or target nodes.
* Automatically updates Route53 to point A records to LB public IPs

Load Balancer Targets:
* List of IPs/Hostnames
* Auto-scaling groups
* Target Groups
* Tagged Instances.

AdaptLB is developed and maintained by [Casey Haakenson](https://twitter.com/millisecond).

## Getting started

To install it locally, [go here]()

1. Install from source, [binary](https://github.com/millisecond/adaptlb/releases) or
   [Docker](https://hub.docker.com/r/millisecond/adaptlb/).
    ```
	# go 1.8 or higher is required
    go get github.com/millisecond/adaptlb                     (>= go1.8)

    docker pull millisecond/adaptlb                           (Docker)

    https://github.com/millisecond/adaptlb/releases           (pre-built binaries)
    ```

2. Create a Target Group

2a. Route 53

3. Create an Auto Scaling Group

* User properties

* IAM Role

* Public/private routes

4. 

5. 

6. 

7. Done

## Maintainers

* Casey Haakenson [@millisecond](https://twitter.com/millisecond)

## License

See [LICENSE](https://github.com/millisecond/adaptlb/blob/master/LICENSE) for details.
