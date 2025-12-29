# Copyright (c) 2025 Steve Taranto <staranto@gmail.com>.
# SPDX-License-Identifier: Apache-2.0

# Test file for WalkTokens with various Terraform HCL features

variable "ami" {
  description = "AMI ID"
  type        = string
  default     = "ami-12345678"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

locals {
  common_tags = {
    Environment = "test"
    Project     = "walk_tokens"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "example" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = var.instance_type

  tags = merge(local.common_tags, {
    Name = "example-instance"
  })

  lifecycle {
    create_before_destroy = true
  }

  provisioner "local-exec" {
    command = "echo 'Instance created'"
  }
}

resource "aws_security_group" "example" {
  name_prefix = "example-"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.common_tags
}

output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}

output "instance_public_ip" {
  description = "The public IP of the EC2 instance"
  value       = aws_instance.example.public_ip
}

module "vpc" {
  source = "./modules/vpc"

  cidr_block = "10.0.0.0/16"
  tags       = local.common_tags
}

check "instance_running" {
  assert {
    condition     = aws_instance.example.instance_state == "running"
    error_message = "EC2 instance is not running"
  }
}
