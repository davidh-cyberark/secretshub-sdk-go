<img alt="CyberArk Banner" src="images/cyberark-banner.jpg">

# CyberArk Secrets Hub - Go SDK

<!--
Author:   David Hisel <david.hisel@cyberark.com>
Updated:  <2025/06/18 16:53:00>
-->

Note that this SDK is a **WORK IN PROGRESS**, and does not implement all the endpoints yet.

## Description

This Go SDK interfaces with CyberArk Identity Adminstration REST endpoints.  It uses the Secrets Hub openapi specification to generate the code.

Here is [a link to the Secrets Hub documentation](https://docs.cyberark.com/secrets-hub-privilege-cloud/latest/en/content/secretshubcontent/sh-developer-lp.htm) for reference.

The openapi spec in this project is modified from the original to allow the `oapi-codegen` tool to generate code.  The original spec can be found in the [Secrets Hub API docs site](https://api-docs.cyberark.com/docs/secretshub-api/).

## Requirements

- go v1.24.3 (https://go.dev/doc/install)

See the [contributing guide](CONTRIBUTING.md) for additional requirements.

## Project Status

**WORK IN PROGRESS** -- breaking changes are expected

Check the examples directory for examples of endpoints that are implemented.

## Maintainers

This project is maintained and updated when there are issues submitted.

[CODEOWNERS](.github/CODEOWNERS)

## Example Code

Look in the [examples/](./examples) folder for example implementations.

## Contributing

We welcome contributions of all kinds to this repository. For
instructions on how to get started and descriptions of our development
workflows, please see our [contributing guide](CONTRIBUTING.md).

## License

Copyright (c) 2025 CyberArk Software Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

For the full license text see [LICENSE](LICENSE).

## Code of Conduct

Summary of Key Principles

- Be respectful to others in the community at all times.
- Report harassing or abusive behavior that you experience or witness at <ReportAbuse@cyberark.com>
- The CyberArk community will not tolerate abusive or disrespectful behavior towards its members; anyone engaging in such behavior will be suspended from the CyberArk community.

For the full document see the [Code of Conduct](CODE_OF_CONDUCT.md).

## Reporting a Vulnerability

If you believe you have found a vulnerability in this repository,
we ask that you follow responsible disclosure guidelines and
contact <product_security@cyberark.com>.

For the full document see the [Security Policies and Procedures](SECURITY.md).
