resource "aws_security_group" "ssh_portfolio" {
  name        = "ssh-portfolio"
  description = "Allow SSH portfolio access"

  # Port 22 — mapped to container 2222, so users just `ssh` normally
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

  tags = {
    Name = "ssh-portfolio"
  }
}
