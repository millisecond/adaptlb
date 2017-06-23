
##### Current version: UNRELEASED / FIRST VERSION IN DEV


AWS-native auto-scaling load balancer with Route53 zone updates.

Fast, modern, and easy to deploy load balancer (HTTP(S) / TCP / UDP) for deploying applications on AWS.

* Gossip cluster allows dynamic scaling, shared state, and fast failover of load-balancer or target nodes.

* Automatically updates Route53 to point A records to LB public IPs

* Valid load balancing targets include List of IPs, Target Groups, and Tagged Instances.

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
