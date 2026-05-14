# driftcheck

> Compares live Terraform state against deployed cloud resources to surface configuration drift.

---

## Installation

```bash
go install github.com/yourusername/driftcheck@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftcheck.git
cd driftcheck
go build -o driftcheck .
```

---

## Usage

Point `driftcheck` at your Terraform state file and let it compare against your live cloud environment:

```bash
driftcheck --state terraform.tfstate --provider aws --region us-east-1
```

**Example output:**

```
[DRIFT] aws_instance.web_server
  expected: instance_type = "t3.micro"
  actual:   instance_type = "t3.small"

[DRIFT] aws_security_group.app_sg
  expected: ingress port 443 = true
  actual:   ingress port 443 = false

[OK] aws_s3_bucket.assets — no drift detected

2 drifted resources found out of 3 checked.
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--state` | Path to Terraform state file | `terraform.tfstate` |
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) | `aws` |
| `--region` | Cloud region to query | required |
| `--output` | Output format (`text`, `json`) | `text` |

---

## Requirements

- Go 1.21+
- Valid cloud provider credentials (e.g., AWS credentials via environment or `~/.aws/credentials`)

---

## License

MIT © 2024 yourusername